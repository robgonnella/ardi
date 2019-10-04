package ardi

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
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
var ArdiDir = path.Join(homeDir, ".ardi")

// DataDir location of data directory inside ~/.ardi
// To avoid polluting an existing arduino-cli installation, ardi
// uses its own data directory to keep cores, libraries and the likes.
var DataDir = path.Join(ArdiDir, "arduino-rpc-client")

// LibConfig used to tell arduino-cli where to find libraries
var LibConfig = "ardi.yaml"

// DepConfig used to tell ardi which libraries to use for a specific project
var DepConfig = "ardi.json"

// GlobalLibConfig returns path to global library directory config file
var GlobalLibConfig = path.Join(DataDir, LibConfig)

// LibraryDirConfig represents yaml config for telling arduino-cli where to find libraries
type LibraryDirConfig struct {
	ProxyType      string                 `yaml:"proxy_type"`
	SketchbookPath string                 `yaml:"sketchbook_path"`
	ArduinoData    string                 `yaml:"arduino_data"`
	BoardManager   map[string]interface{} `yaml:"board_manager,flow"`
}

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

type boardInfo struct {
	FQBN string
	Name string
}

// SetLogLevel sets log level for ardi
func SetLogLevel(level log.Level) {
	logger.SetLevel(level)
}

// GetDesiredBoard prints list of supported boards and asks user to choose
func GetDesiredBoard(client rpc.ArduinoCoreClient, instance *rpc.Instance) string {
	logger.Debug("Getting list of supported boards...")
	listResp, err := client.PlatformList(
		context.Background(),
		&rpc.PlatformListReq{Instance: instance},
	)

	if err != nil {
		logger.Fatalf("List error: %s", err)
	}

	var boards []boardInfo

	for _, p := range listResp.GetInstalledPlatform() {
		for _, b := range p.GetBoards() {
			b := boardInfo{
				FQBN: b.GetFqbn(),
				Name: b.GetName(),
			}
			boards = append(boards, b)
		}
	}

	var boardIdx int
	printSupportedBoardsWithIndices(boards)

	fmt.Print("\nEnter number of board for which to compile: ")
	if _, err := fmt.Scanf("%d", &boardIdx); err != nil {
		logger.WithError(err).Fatal("Failed to parse target board")
	}

	if boardIdx < 0 || boardIdx > len(boards)-1 {
		logger.WithError(errors.New("Invalid board selection")).Fatal("Failed to get desired board")
	}

	return boards[boardIdx].FQBN
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
		printConnectedBoardsWithIndices(list)
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

// LibSearch searches all available libraries with optional search filter
func LibSearch(client rpc.ArduinoCoreClient, instance *rpc.Instance, searchArg string) {
	searchResp, err := client.LibrarySearch(
		context.Background(),
		&rpc.LibrarySearchReq{
			Instance: instance,
			Query:    searchArg,
		},
	)
	if err != nil {
		logger.WithError(err).Fatal("Error searching library")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "Library\tLatest\tReleases")
	for _, lib := range searchResp.GetLibraries() {
		releases := []string{}
		for _, rel := range lib.GetReleases() {
			releases = append(releases, rel.GetVersion())
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", lib.GetName(), lib.GetLatest().GetVersion(), strings.Join(releases, ", "))
	}
}

// LibInstall installs library either globally or for project
func LibInstall(client rpc.ArduinoCoreClient, instance *rpc.Instance, name, version string) string {
	logger.Infof("Installing library: %s %s", name, version)
	installRespStream, err := client.LibraryInstall(context.Background(),
		&rpc.LibraryInstallReq{
			Instance: instance,
			Name:     name,
			Version:  version,
		})

	if err != nil {
		logger.WithError(err).Fatal("Error installing library")
	}

	foundVersion := ""

	for {
		installResp, err := installRespStream.Recv()
		if err == io.EOF {
			logger.Info("Lib install done")
			break
		}

		if err != nil {
			logger.WithError(err).Fatal("Library install error")
		}

		if installResp.GetProgress() != nil {
			logger.Infof("DOWNLOAD: %s\n", installResp.GetProgress())
		}
		if installResp.GetTaskProgress() != nil {
			msg := installResp.GetTaskProgress()
			lib := msg.GetName()
			logger.Infof("TASK: %s\n", msg)
			if foundVersion == "" {
				foundVersion = strings.Split(lib, "@")[1]
			}
		}
	}
	return foundVersion
}

// LibUnInstall installs library either globally or for project
func LibUnInstall(client rpc.ArduinoCoreClient, instance *rpc.Instance, name string) {
	logger.Infof("Uninstalling library: %s", name)
	uninstallRespStream, err := client.LibraryUninstall(
		context.Background(),
		&rpc.LibraryUninstallReq{
			Instance: instance,
			// Assume spaces in name were intended to be underscore. This indicates
			// a potential bug in the arduino-cli package manager as names
			// potentially do not have a one-to-one mapping with regards to install
			// and remove commands. It seems as though arduino should be forcing
			// devs to name their library according to the github url.
			// @todo there has to be a better way - find it!
			Name: strings.ReplaceAll(name, " ", "_"),
		})

	if err != nil {
		logger.WithError(err).Fatal("Error uninstalling library")
	}

	for {
		uninstallRespStream, err := uninstallRespStream.Recv()
		if err == io.EOF {
			logger.Info("Lib uninstall done")
			break
		}

		if err != nil {
			logger.WithError(err).Fatal("Library install error")
		}

		if uninstallRespStream.GetTaskProgress() != nil {
			logger.Infof("TASK: %s\n", uninstallRespStream.GetTaskProgress())
		}
	}
}

// IsInitialized returns whether or not ardi had been initialized (has a data directory)
func IsInitialized() bool {
	_, err := os.Stat(DataDir)
	return !os.IsNotExist(err)
}

// IsProjectDirectory returns whether or not current directory is configured with
func IsProjectDirectory() bool {
	_, err := os.Stat(LibConfig)
	return !os.IsNotExist(err)
}

// StartDaemonAndGetConnection starts daemon as goroutine and return connection, client, and rpc-instance
func StartDaemonAndGetConnection(pathToConfig string) (*grpc.ClientConn, rpc.ArduinoCoreClient, *rpc.Instance) {
	go startDaemon(pathToConfig)
	conn := getServerConnection()
	client := rpc.NewArduinoCoreClient(conn)
	instance := getRPCInstance(client, pathToConfig)
	return conn, client, instance
}

// ListPlatforms list all available platforms or filter with a search arg
func ListPlatforms(platform string) {
	if !IsInitialized() {
		createDataDir()
	}
	conn, client, instance := StartDaemonAndGetConnection(GlobalLibConfig)
	defer conn.Close()

	updateIndexes(client, instance)
	searchResp, err := client.PlatformSearch(
		context.Background(),
		&rpc.PlatformSearchReq{
			Instance:   instance,
			SearchArgs: platform,
		},
	)

	if err != nil {
		logger.Fatalf("Search error: %s", err)
	}

	platforms := searchResp.GetSearchOutput()
	logger.Info("------AVAILABLE PLATFORMS------")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', 0)
	defer w.Flush()
	fmt.Fprintln(w, "Platform\tID")
	for _, plat := range platforms {
		fmt.Fprintf(w, "%s\t%s\n", plat.GetName(), plat.GetID())
	}
}

// ListBoards lists all available boards with optional search filter
func ListBoards(board string) {
	if !IsInitialized() {
		createDataDir()
	}
	conn, client, instance := StartDaemonAndGetConnection(GlobalLibConfig)
	defer conn.Close()

	updateIndexes(client, instance)
	searchResp, err := client.PlatformSearch(
		context.Background(),
		&rpc.PlatformSearchReq{
			Instance:   instance,
			SearchArgs: board,
		},
	)

	if err != nil {
		logger.Fatalf("Search error: %s", err)
	}

	platforms := searchResp.GetSearchOutput()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', 0)
	defer w.Flush()
	fmt.Fprintln(w, "Board\tPlatform\tFQBN")
	for _, plat := range platforms {
		for _, board := range plat.GetBoards() {
			fmt.Fprintf(w, "%s\t%s\t%s\n", board.GetName(), plat.GetID(), board.GetFqbn())
		}
	}
}

// Initialize downloads and installs all available platforms for maximum board support
func Initialize(platform, version string) {
	if !IsInitialized() {
		createDataDir()
	}

	conn, client, rpcInstance := StartDaemonAndGetConnection(GlobalLibConfig)
	defer conn.Close()

	quit := make(chan bool, 1)
	// Show simple "processing" indicator if not logging verbosely
	if !isVerbose() {
		logger.Info("Installing platforms...")
		ticker := time.NewTicker(2 * time.Second)
		go func() {
			for {
				select {
				case <-ticker.C:
					fmt.Print(".")
				case <-quit:
					ticker.Stop()
				}
			}
		}()
	}
	updateIndexes(client, rpcInstance)
	if platform == "" {
		loadAllPlatforms(client, rpcInstance)
	} else {
		platParts := strings.Split(platform, ":")
		platPackage := platParts[0]
		arch := platParts[len(platParts)-1]
		version := ""
		done := make(chan platformInstallMessage, 1)
		platformInstall(client, rpcInstance, platPackage, arch, version, done)
	}
	if !isVerbose() {
		quit <- true
		fmt.Println("")
	}
	platformList(client, rpcInstance)
}
