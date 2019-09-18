package ardi

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	arduino "github.com/arduino/arduino-cli/cli"
	rpc "github.com/arduino/arduino-cli/rpc/commands"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
	"google.golang.org/grpc"
)

var cli = arduino.ArduinoCli
var logger = log.New()
var homeDir, _ = os.UserHomeDir()

// ArdiDir location of .ardi directory in users home directory
var ArdiDir = fmt.Sprintf("%s/.ardi", homeDir)

// DataDir location of data directory inside ~/.ardi
// To avoid polluting an existing arduino-cli installation, ardi
// uses its own data directory to keep cores, libraries and the likes.
var DataDir = fmt.Sprintf("%s/arduino-rpc-client", ArdiDir)

// TargetInfo represents all necessary info for compiling, and uploading
type TargetInfo struct {
	FQBN         string
	Device       string
	SketchDir    string
	SketchFile   string
	Baud         int
	Stream       *serial.Port
	Compiling    bool
	CompileError bool
	Uploading    bool
	Logging      bool
}

type platformUpgradeMessage struct {
	platformPackage string
	architecture    string
	success         bool
}

type platformInstallMessage struct {
	platformPackage string
	architecture    string
	version         string
	success         bool
}

// SetLogLevel sets log level for ardi
func SetLogLevel(level log.Level) {
	logger.SetLevel(level)
}

// WatchLogs connects to a serial port at a specified baud rate and prints
// any logs received.
func WatchLogs(target *TargetInfo) {
	logFields := log.Fields{"baud": target.Baud, "device": target.Device}

	stopLogs(target)
	waitForPreviousCompile(target)
	waitForPreviousUpload(target)

	logger.Info("Watching logs...")
	target.Logging = true

	config := &serial.Config{Name: target.Device, Baud: target.Baud}
	stream, err := serial.OpenPort(config)
	if err != nil {
		logger.WithError(err).WithFields(logFields).Fatal("Failed to read from device")
		return
	}

	target.Stream = stream

	for {
		if target.Stream == nil {
			target.Logging = false
			break
		}
		var buf = make([]byte, 128)
		n, err := stream.Read(buf)
		if err != nil {
			logger.WithError(err).WithFields(logFields).Fatal("Failed to read from serial port")
		}
		fmt.Printf("%s", buf[:n])
	}

}

// Upload compiled sketches to the specified board
func Upload(client rpc.ArduinoCoreClient, instance *rpc.Instance, target *TargetInfo) {

	stopLogs(target)
	waitForPreviousCompile(target)
	waitForPreviousUpload(target)

	logger.Info("Uploading...")

	target.Uploading = true

	uplRespStream, err := client.Upload(context.Background(),
		&rpc.UploadReq{
			Instance:   instance,
			Fqbn:       target.FQBN,
			SketchPath: target.SketchDir,
			Port:       target.Device,
			Verbose:    isVerbose(),
		})

	if err != nil {
		logger.WithError(err).Fatal("Failed to upload")
	}

	for {
		uplResp, err := uplRespStream.Recv()
		if err == io.EOF {
			target.Uploading = false
			logger.Info("Upload done")
			break
		}

		if err != nil {
			logger.WithError(err).Fatal("Failed to upload")
			break
		}

		// When an operation is ongoing you can get its output
		if resp := uplResp.GetOutStream(); resp != nil {
			logger.Debugf("STDOUT: %s", resp)
		}
		if resperr := uplResp.GetErrStream(); resperr != nil {
			logger.Debugf("STDERR: %s", resperr)
		}
	}
}

// Compile the specified sketch
func Compile(client rpc.ArduinoCoreClient, instance *rpc.Instance, target *TargetInfo) {

	stopLogs(target)
	waitForPreviousCompile(target)
	waitForPreviousUpload(target)

	logger.Info("Compiling...")

	target.Compiling = true
	target.CompileError = false

	compRespStream, err := client.Compile(context.Background(),
		&rpc.CompileReq{
			Instance:   instance,
			Fqbn:       target.FQBN,
			SketchPath: target.SketchDir,
			Verbose:    true,
		})

	if err != nil {
		logger.WithError(err).Fatal("Failed to compile")
	}

	// Loop and consume the server stream until all the operations are done.
	for {
		compResp, err := compRespStream.Recv()

		// The server is done.
		if err == io.EOF {
			target.Compiling = false
			logger.Info("Compilation done")
			break
		}

		// There was an error.
		if err != nil {
			target.CompileError = true
			target.Compiling = false
			logger.WithError(err).Error("Failed to compile")
			break
		}

		// When an operation is ongoing you can get its output
		if resp := compResp.GetOutStream(); resp != nil {
			logger.Debugf("STDOUT: %s", resp)
		}
		if resperr := compResp.GetErrStream(); resperr != nil {
			logger.Errorf("STDERR: %s", resperr)
		}
	}
}

// GetTargetInfo returns a connected board if found. If more than
// one board is connected it will ask the user to choose.
func GetTargetInfo(list []TargetInfo) TargetInfo {
	var boardIndex int
	target := TargetInfo{}
	listLength := len(list)

	if listLength == 0 {
		logger.WithError(errors.New("No boards detected")).Fatal("Failed to get target board")
	} else if listLength == 1 {
		boardIndex = 0
	} else {
		printBoardListWithIndices(list)
		fmt.Print("\nEnter number of board to upload to: ")
		if _, err := fmt.Scanf("%d", &boardIndex); err != nil {
			logger.WithError(err).Fatal("Failed to parse target board")
		}
	}

	if boardIndex < 0 || boardIndex > listLength-1 {
		logger.WithError(errors.New("Invalid board selection")).Fatal("Failed to parse target board")
	}

	target = list[boardIndex]
	return target
}

// GetTargetList returns a list of connected boards and their corresponding info
func GetTargetList(client rpc.ArduinoCoreClient, instance *rpc.Instance, sketchDir, sketchFile string, baud int) []TargetInfo {
	logger.Debug("Getting target list...")

	var boardList []TargetInfo

	boardListResp, err := client.BoardList(
		context.Background(),
		&rpc.BoardListReq{Instance: instance},
	)

	if err != nil {
		logger.Fatalf("Board list error: %s\n", err)
	}

	for _, port := range boardListResp.GetPorts() {
		for _, board := range port.GetBoards() {
			target := TargetInfo{
				FQBN:       board.GetFQBN(),
				Device:     port.GetAddress(),
				Baud:       baud,
				SketchDir:  sketchDir,
				SketchFile: sketchFile,
			}
			boardList = append(boardList, target)
		}
	}

	return boardList
}

// WatchSketch responds to changes in a given sketch file by automatically
// recompiling and re-uploading.
func WatchSketch(client rpc.ArduinoCoreClient, instance *rpc.Instance, target *TargetInfo) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.WithError(err).Fatal("Failed to watch directory for changes")
	}
	defer watcher.Close()

	err = watcher.Add(target.SketchFile)
	if err != nil {
		logger.WithError(err).Fatal("Failed to watch directory for changes")
	}

	go WatchLogs(target)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				break
			}
			logger.Debugf("event: %+v", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				logger.Debugf("modified file: %s", event.Name)
				Compile(client, instance, target)
				if !target.CompileError {
					Upload(client, instance, target)
					go WatchLogs(target)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logger.WithError(err).Warn("Watch error")
		}
	}

}

// Initialize downloads and installs all available platforms for maximum board support
func Initialize() (*grpc.ClientConn, rpc.ArduinoCoreClient, *rpc.Instance) {
	createDataDirIfNeeded()
	go startDaemon()

	conn := getServerConnection()
	client := rpc.NewArduinoCoreClient(conn)
	rpcInstance := getRPCInstance(client, DataDir)
	quit := make(chan bool)
	// Show simple "processing" indicator if not logging verbosely
	if !isVerbose() {
		ticker := time.NewTicker(2 * time.Second)
		go func() {
			for {
				select {
				case <-ticker.C:
					fmt.Print(".")
				case <-quit:
					fmt.Print(".\n")
					ticker.Stop()
				}
			}
		}()
	}
	updateIndex(client, rpcInstance)
	loadPlatforms(client, rpcInstance)
	platformList(client, rpcInstance)
	quit <- true
	return conn, client, rpcInstance
}

// private
func printBoardListWithIndices(list []TargetInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 5, 0, '\t', 0)
	defer w.Flush()
	fmt.Fprintln(w, "No.\tBoard\tDevice")
	for i, board := range list {
		fmt.Fprintf(w, "%d\t%s\t%s\n", i, board.FQBN, board.Device)
	}
}

func platformList(client rpc.ArduinoCoreClient, instance *rpc.Instance) {
	listResp, err := client.PlatformList(context.Background(),
		&rpc.PlatformListReq{Instance: instance})

	if err != nil {
		logger.Fatalf("List error: %s", err)
	}

	logger.Debug("------INSTALLED PLATFORMS------")
	for _, plat := range listResp.GetInstalledPlatform() {
		logger.Debugf("Installed platform: %s - %s", plat.GetID(), plat.GetInstalled())
	}
	logger.Debug("-------------------------------")
}

func platformUpgrade(client rpc.ArduinoCoreClient, instance *rpc.Instance, platPackage, arch string, done chan platformUpgradeMessage) {
	logger.Debugf("Upgrading platform: %s:%s\n", platPackage, arch)

	upgradeRespStream, err := client.PlatformUpgrade(context.Background(),
		&rpc.PlatformUpgradeReq{
			Instance:        instance,
			PlatformPackage: platPackage,
			Architecture:    arch,
		})

	if err != nil {
		logger.WithError(err).Warn("Error upgrading platform")
	}

	message := platformUpgradeMessage{
		platformPackage: platPackage,
		architecture:    arch,
		success:         false,
	}

	// Loop and consume the server stream until all the operations are done.
	for {
		upgradeResp, err := upgradeRespStream.Recv()

		// The server is done.
		if err == io.EOF {
			logger.Debug("Upgrade done")
			message.success = true
			done <- message
			break
		}

		// There was an error.
		if err != nil {
			if !strings.Contains(err.Error(), "platform already at latest version") {
				logger.WithError(err).Warn("Cannot upgrade platform")
			}
			done <- message
			break
		}

		// When a download is ongoing, log the progress
		if upgradeResp.GetProgress() != nil {
			logger.Debugf("DOWNLOAD: %s", upgradeResp.GetProgress())
		}

		// When an overall task is ongoing, log the progress
		if upgradeResp.GetTaskProgress() != nil {
			logger.Debugf("TASK: %s", upgradeResp.GetTaskProgress())
		}
	}
}

func upgradePlatforms(client rpc.ArduinoCoreClient, instance *rpc.Instance, platforms []*rpc.Platform) {
	count := len(platforms)
	completed := 0
	done := make(chan platformUpgradeMessage, count)
	for _, plat := range platforms {
		id := plat.GetID()
		idParts := strings.Split(id, ":")
		platPackage := idParts[0]
		arch := idParts[len(idParts)-1]
		go platformUpgrade(client, instance, platPackage, arch, done)
	}
	for message := range done {
		if message.success {
			logger.Debugf("Successfully upgraded %s:%s", message.platformPackage, message.architecture)
		}
		completed++
		if completed == count {
			close(done)
		}
	}
}

func platformInstall(client rpc.ArduinoCoreClient, instance *rpc.Instance, platPackage, arch, version string, done chan platformInstallMessage) {
	logger.Debugf("Installing platform: %s:%s\n", arch, version)

	installRespStream, err := client.PlatformInstall(context.Background(),
		&rpc.PlatformInstallReq{
			Instance:        instance,
			PlatformPackage: platPackage,
			Architecture:    arch,
			Version:         version,
		})

	if err != nil {
		logger.WithError(err).Warn("Failed to install platform")
	}

	message := platformInstallMessage{
		platformPackage: platPackage,
		architecture:    arch,
		version:         version,
		success:         false,
	}

	// Loop and consume the server stream until all the operations are done.
	for {
		installResp, err := installRespStream.Recv()

		// The server is done.
		if err == io.EOF {
			logger.Debug("Install done")
			message.success = true
			done <- message
			break
		}

		// There was an error.
		if err != nil {
			logger.WithError(err).Warn("Failed to install platform")
			done <- message
			break
		}

		// When a download is ongoing, log the progress
		if installResp.GetProgress() != nil {
			logger.Debugf("DOWNLOAD: %s", installResp.GetProgress())
		}

		// When an overall task is ongoing, log the progress
		if installResp.GetTaskProgress() != nil {
			logger.Debugf("TASK: %s", installResp.GetTaskProgress())
		}
	}
}

func loadPlatforms(client rpc.ArduinoCoreClient, instance *rpc.Instance) {
	searchResp, err := client.PlatformSearch(context.Background(), &rpc.PlatformSearchReq{
		Instance: instance,
	})

	if err != nil {
		logger.Fatalf("Search error: %s", err)
	}

	platforms := searchResp.GetSearchOutput()
	count := len(platforms)
	completed := 0
	done := make(chan platformInstallMessage, count)
	for _, plat := range platforms {
		id := plat.GetID()
		idParts := strings.Split(id, ":")
		platPackage := idParts[0]
		arch := idParts[len(idParts)-1]
		latest := plat.GetLatest()
		logger.Debugf("Search result: %s: %s - %s", platPackage, id, latest)
		go platformInstall(client, instance, platPackage, arch, latest, done)
	}
	for message := range done {
		if message.success {
			logger.Debugf("Successfully installed %s:%s - %s", message.platformPackage, message.architecture, message.version)

		}
		completed++
		if completed == count {
			close(done)
		}
	}
	upgradePlatforms(client, instance, platforms)
}

func updateIndex(client rpc.ArduinoCoreClient, instance *rpc.Instance) {
	logger.Debug("Updating index...")
	uiRespStream, err := client.UpdateIndex(context.Background(), &rpc.UpdateIndexReq{
		Instance: instance,
	})
	if err != nil {
		logger.Fatalf("Error updating index: %s", err)
	}

	// Loop and consume the server stream until all the operations are done.
	for {
		uiResp, err := uiRespStream.Recv()

		// the server is done
		if err == io.EOF {
			logger.Debug("Update index done")
			break
		}

		// there was an error
		if err != nil {
			logger.Fatalf("Update error: %s", err)
		}

		// operations in progress
		if uiResp.GetDownloadProgress() != nil {
			logger.Debugf("DOWNLOAD: %s", uiResp.GetDownloadProgress())
		}
	}
}

func getRPCInstance(client rpc.ArduinoCoreClient, dataDir string) *rpc.Instance {
	// The configuration for this example client only contains the path to
	// the data folder.
	initRespStream, err := client.Init(context.Background(), &rpc.InitReq{
		Configuration: &rpc.Configuration{
			DataDir: dataDir,
		},
	})
	if err != nil {
		logger.Fatalf("Error creating server instance: %s", err)
	}

	var instance *rpc.Instance
	// Loop and consume the server stream until all the setup procedures are done.
	for {
		initResp, err := initRespStream.Recv()
		// The server is done.
		if err == io.EOF {
			break
		}

		// There was an error.
		if err != nil {
			logger.Fatalf("Init error: %s", err)
		}

		// The server sent us a valid instance, let's print its ID.
		if initResp.GetInstance() != nil {
			instance = initResp.GetInstance()
			logger.Debugf("Got a new instance with ID: %v", instance.GetId())
		}

		// When a download is ongoing, log the progress
		if initResp.GetDownloadProgress() != nil {
			logger.Debugf("DOWNLOAD: %s", initResp.GetDownloadProgress())
		}

		// When an overall task is ongoing, log the progress
		if initResp.GetTaskProgress() != nil {
			logger.Debugf("TASK: %s", initResp.GetTaskProgress())
		}
	}

	return instance
}

func getServerConnection() *grpc.ClientConn {
	backgroundCtx := context.Background()
	ctx, _ := context.WithTimeout(backgroundCtx, time.Second)
	// Establish a connection with the gRPC server, started with the command: arduino-cli daemon
	conn, err := grpc.DialContext(ctx, "localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Fatal("error connecting to arduino-cli rpc server, you can start it by running `arduino-cli daemon`")
	}
	return conn
}

func startDaemon() {
	logger.Debug("Starting daemon")
	cli.SetArgs([]string{"daemon"})
	if err := cli.Execute(); err != nil {
		logger.WithError(err).Fatal("Failed to start rpc server")
	}
	logger.Debug("Daemon started")
}

func createDataDirIfNeeded() {
	logger.Debug("Creating data directory if needed")
	_ = os.MkdirAll(DataDir, 0777)
}

func stopLogs(target *TargetInfo) {
	if target.Stream != nil {
		logWithField := logger.WithField("device", target.Device)
		logWithField.Info("Closing serial port connection")
		if err := target.Stream.Close(); err != nil {
			logWithField.WithError(err).Fatal("Failed to close serial port connection")
		}
		if err := target.Stream.Flush(); err != nil {
			logWithField.WithError(err).Fatal("Failed to flush serial port connection")
		}
		target.Stream = nil
		// block until all logs have stopped
		for {
			if !target.Logging {
				break
			}
		}
	}
}

func waitForPreviousUpload(target *TargetInfo) {
	// block until target is no longer uploading
	for {
		if !target.Uploading {
			break
		}
		logger.Info("Waiting for previous upload to finish...")
		time.Sleep(time.Second)
	}
}

func waitForPreviousCompile(target *TargetInfo) {
	// block until target is no longer compiling
	for {
		if !target.Compiling {
			break
		}
		logger.Info("Waiting for previous compile to finish...")
		time.Sleep(time.Second)
	}
}

func isVerbose() bool {
	return logger.Level == log.DebugLevel
}
