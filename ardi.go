/*
ardi is a command-line tool for compiling, uploading code, and
watching logs for your usb connected arduino board. This allows you to
develop in an environment you feel comfortable in, without needing to
use arduino's web or desktop IDEs.

Usage:
  ardi [sketch] [flags]
  ardi [command]

Available Commands:
  clean       Delete all ardi data
  help        Help about any command
  init        Download and install platforms

Flags:
  -b, --baud int   specify sketch baud rate (default 9600)
  -h, --help       help for ardi
*/
package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	arduino "github.com/arduino/arduino-cli/cli"
	rpc "github.com/arduino/arduino-cli/rpc/commands"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tarm/serial"
	"google.golang.org/grpc"
)

var cli = arduino.ArduinoCli
var logger = log.New()

// To avoid polluting an existing arduino-cli installation, the example
// client uses a temp folder to keep cores, libraries and the likes.
var homeDir, _ = os.UserHomeDir()
var ardiDir = fmt.Sprintf("%s/.ardi", homeDir)
var dataDir = fmt.Sprintf("%s/arduino-rpc-client", ardiDir)

type targetBoardInfo struct {
	FQBN   string
	Device string
}

type PlatformUpgradeMessage struct {
	platformPackage string
	architecture    string
	success         bool
}

type PlatformInstallMessage struct {
	platformPackage string
	architecture    string
	version         string
	success         bool
}

func filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func watchLogs(device string, baud int) {
	logFields := log.Fields{"baud": baud, "device": device}

	config := &serial.Config{Name: device, Baud: baud}
	stream, err := serial.OpenPort(config)
	if err != nil {
		logger.WithError(err).WithFields(logFields).Fatal("Failed to read from device")
		return
	}

	for {
		var buf = make([]byte, 128)
		n, err := stream.Read(buf)
		if err != nil {
			logger.WithError(err).WithFields(logFields).Fatal("Failed to read from serial port")
		}
		fmt.Printf("%s", buf[:n])
	}

}

func upload(client rpc.ArduinoCoreClient, instance *rpc.Instance, target targetBoardInfo, sketch string) {
	logger.Info("uploading...")

	uplRespStream, err := client.Upload(context.Background(),
		&rpc.UploadReq{
			Instance:   instance,
			Fqbn:       target.FQBN,
			SketchPath: sketch,
			Port:       target.Device,
			Verbose:    true,
		})

	if err != nil {
		logger.WithError(err).Fatal("Failed to upload")
	}

	for {
		uplResp, err := uplRespStream.Recv()
		if err == io.EOF {
			logger.Info("Upload done")
			break
		}

		if err != nil {
			logger.WithError(err).Fatal("Failed to upload")
			break
		}

		// When an operation is ongoing you can get its output
		if resp := uplResp.GetOutStream(); resp != nil {
			logger.Infof("STDOUT: %s", resp)
		}
		if resperr := uplResp.GetErrStream(); resperr != nil {
			logger.Infof("STDERR: %s", resperr)
		}
	}
}

func compile(client rpc.ArduinoCoreClient, instance *rpc.Instance, target targetBoardInfo, sketch string) {
	logger.Info("compiling...")

	compRespStream, err := client.Compile(context.Background(),
		&rpc.CompileReq{
			Instance:   instance,
			Fqbn:       target.FQBN,
			SketchPath: sketch,
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
			logger.Info("Compilation done")
			break
		}

		// There was an error.
		if err != nil {
			logger.WithError(err).Fatal("Failed to compile")
		}

		// When an operation is ongoing you can get its output
		if resp := compResp.GetOutStream(); resp != nil {
			logger.Infof("STDOUT: %s", resp)
		}
		if resperr := compResp.GetErrStream(); resperr != nil {
			logger.Infof("STDERR: %s", resperr)
		}
	}
}

func printBoardListWithIndices(list []targetBoardInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 5, 0, '\t', 0)
	defer w.Flush()
	fmt.Fprintln(w, "No.\tBoard\tDevice")
	for i, board := range list {
		fmt.Fprintf(w, "%d\t%s\t%s\n", i, board.FQBN, board.Device)
	}
}

func getTargetBoardInfo(list []targetBoardInfo) targetBoardInfo {
	var boardIndex int
	target := targetBoardInfo{}
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

func getBoardList(client rpc.ArduinoCoreClient, instance *rpc.Instance) []targetBoardInfo {
	logger.Info("Getting board list...")

	var boardList []targetBoardInfo

	boardListResp, err := client.BoardList(
		context.Background(),
		&rpc.BoardListReq{Instance: instance},
	)

	if err != nil {
		logger.Fatalf("Board list error: %s\n", err)
	}

	for _, port := range boardListResp.GetPorts() {
		for _, board := range port.GetBoards() {
			logger.Infof("port: %s, board: %+v\n", port.GetAddress(), board)
			target := targetBoardInfo{
				FQBN:   board.GetFQBN(),
				Device: port.GetAddress(),
			}
			boardList = append(boardList, target)
		}
	}

	return boardList
}

func platformList(client rpc.ArduinoCoreClient, instance *rpc.Instance) {
	listResp, err := client.PlatformList(context.Background(),
		&rpc.PlatformListReq{Instance: instance})

	if err != nil {
		logger.Fatalf("List error: %s", err)
	}

	logger.Info("------INSTALLED PLATFORMS------")
	for _, plat := range listResp.GetInstalledPlatform() {
		// We only print ID and version of the installed platforms but you can look
		// at the definition for the rpc.Platform struct for more fields.
		logger.Infof("Installed platform: %s - %s", plat.GetID(), plat.GetInstalled())
	}
	logger.Info("-------------------------------")
}

func platformUpgrade(client rpc.ArduinoCoreClient, instance *rpc.Instance, platPackage, arch string, done chan PlatformUpgradeMessage) {
	logger.Infof("Upgrading platform: %s:%s\n", platPackage, arch)

	upgradeRespStream, err := client.PlatformUpgrade(context.Background(),
		&rpc.PlatformUpgradeReq{
			Instance:        instance,
			PlatformPackage: platPackage,
			Architecture:    arch,
		})

	if err != nil {
		logger.WithError(err).Warn("Error upgrading platform")
	}

	message := PlatformUpgradeMessage{
		platformPackage: platPackage,
		architecture:    arch,
		success:         false,
	}

	// Loop and consume the server stream until all the operations are done.
	for {
		upgradeResp, err := upgradeRespStream.Recv()

		// The server is done.
		if err == io.EOF {
			logger.Infof("Upgrade done")
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
			logger.Infof("DOWNLOAD: %s", upgradeResp.GetProgress())
		}

		// When an overall task is ongoing, log the progress
		if upgradeResp.GetTaskProgress() != nil {
			logger.Infof("TASK: %s", upgradeResp.GetTaskProgress())
		}
	}
}

func upgradePlatforms(client rpc.ArduinoCoreClient, instance *rpc.Instance, platforms []*rpc.Platform) {
	count := len(platforms)
	completed := 0
	done := make(chan PlatformUpgradeMessage, count)
	for _, plat := range platforms {
		id := plat.GetID()
		idParts := strings.Split(id, ":")
		platPackage := idParts[0]
		arch := idParts[len(idParts)-1]
		go platformUpgrade(client, instance, platPackage, arch, done)
	}
	for message := range done {
		if message.success {
			logger.Infof("Successfully upgraded %s:s", message.platformPackage, message.architecture)
		}
		completed++
		if completed == count {
			close(done)
		}
	}
}

func platformInstall(client rpc.ArduinoCoreClient, instance *rpc.Instance, platPackage, arch, version string, done chan PlatformInstallMessage) {
	logger.Infof("Installing platform: %s:%s\n", arch, version)

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

	message := PlatformInstallMessage{
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
			logger.Info("Install done")
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
			logger.Infof("DOWNLOAD: %s", installResp.GetProgress())
		}

		// When an overall task is ongoing, log the progress
		if installResp.GetTaskProgress() != nil {
			logger.Infof("TASK: %s", installResp.GetTaskProgress())
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
	done := make(chan PlatformInstallMessage, count)
	for _, plat := range platforms {
		// We only print ID and version of the platforms found but you can look
		// at the definition for the rpc.Platform struct for more fields.
		id := plat.GetID()
		idParts := strings.Split(id, ":")
		platPackage := idParts[0]
		arch := idParts[len(idParts)-1]
		latest := plat.GetLatest()
		logger.Infof("Search result: %s: %s - %s", platPackage, id, latest)
		go platformInstall(client, instance, platPackage, arch, latest, done)
	}
	for message := range done {
		if message.success {
			logger.Infof("Successfully installed %s:%s - %s", message.platformPackage, message.architecture, message.version)

		}
		completed++
		if completed == count {
			close(done)
		}
	}
	upgradePlatforms(client, instance, platforms)
}

func updateIndex(client rpc.ArduinoCoreClient, instance *rpc.Instance) {
	logger.Info("updating index...")
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
			logger.Info("Update index done")
			break
		}

		// there was an error
		if err != nil {
			logger.Fatalf("Update error: %s", err)
		}

		// operations in progress
		if uiResp.GetDownloadProgress() != nil {
			logger.Infof("DOWNLOAD: %s", uiResp.GetDownloadProgress())
		}
	}
}

func getRpcInstance(client rpc.ArduinoCoreClient, dataDir string) *rpc.Instance {
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
			logger.Infof("Got a new instance with ID: %v", instance.GetId())
		}

		// When a download is ongoing, log the progress
		if initResp.GetDownloadProgress() != nil {
			logger.Infof("DOWNLOAD: %s", initResp.GetDownloadProgress())
		}

		// When an overall task is ongoing, log the progress
		if initResp.GetTaskProgress() != nil {
			logger.Infof("TASK: %s", initResp.GetTaskProgress())
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
	logger.Info("Starting daemon")
	cli.SetArgs([]string{"daemon"})
	if err := cli.Execute(); err != nil {
		logger.WithError(err).Fatal("Failed to start rpc server")
	}
	logger.Info("Daemon started")
}

func createDataDirIfNeeded() {
	logger.Info("Creating data directory if needed")
	_ = os.MkdirAll(dataDir, 0777)
}

func parseBaudRate(sketchPath string) int {
	var baud int
	rgx := regexp.MustCompile(`Serial\.begin\((\d+)\);`)
	sketchParts := strings.Split(sketchPath, "/")
	sketchName := sketchParts[len(sketchParts)-1]
	sketchFile := fmt.Sprintf("%s/%s.ino", sketchPath, sketchName)
	file, err := os.Open(sketchFile)
	if err != nil {
		// Log the error and return 0 for baud to let script continue
		// with either default value or value specified from command-line.
		logger.WithError(err).
			WithField("sketch", sketchPath).
			Info("Failed to read sketch")
		return baud
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := scanner.Text()
		if match := rgx.MatchString(text); match {
			stringBaud := strings.TrimSpace(rgx.ReplaceAllString(text, "$1"))
			if baud, err = strconv.Atoi(stringBaud); err != nil {
				// set baud to 0 and let script continue with either default
				// value or value specified from command-line.
				logger.WithError(err).Info("Failed to parse baud rate from sketch")
				baud = 0
			}
			break
		}
	}

	return baud
}

func getSketch() string {
	if len(os.Args) == 1 {
		return ""
	}

	sketch := os.Args[1]

	if !strings.Contains(sketch, "/") {
		return fmt.Sprintf("sketches/%s", sketch)
	}

	if strings.HasSuffix(sketch, "/") {
		sketch = strings.TrimSuffix(sketch, "/")
	}

	return sketch
}

func processSketch(baud int) (string, int) {
	sketch := getSketch()

	if sketch == "" {
		logger.WithError(errors.New("Missing sketch arguemnet")).Fatal("Must provide a sketch name as an argument to upload")
	}
	parsedBaud := parseBaudRate(sketch)

	if parsedBaud != 0 && parsedBaud != baud {
		fmt.Println("")
		logger.Info("Detected a different baud rate from sketch file.")
		logger.WithField("detected baud", parsedBaud).Info("Using detected baud rate")
		fmt.Println("")
		baud = parsedBaud
	}

	return sketch, baud
}

func initialize() (*grpc.ClientConn, rpc.ArduinoCoreClient, *rpc.Instance) {
	createDataDirIfNeeded()
	go startDaemon()

	conn := getServerConnection()
	client := rpc.NewArduinoCoreClient(conn)
	rpcInstance := getRpcInstance(client, dataDir)

	updateIndex(client, rpcInstance)
	loadPlatforms(client, rpcInstance)
	platformList(client, rpcInstance)
	return conn, client, rpcInstance
}

func process(baud int) {
	sketch, baud := processSketch(baud)

	logFields := log.Fields{"baud": baud, "sketch": sketch}
	logWithFields := logger.WithFields(logFields)

	conn, client, rpcInstance := initialize()
	defer conn.Close()

	list := getBoardList(client, rpcInstance)

	logWithFields.Info("Parsing target board")
	targetBoard := getTargetBoardInfo(list)

	logWithFields.WithField("target-board", targetBoard).Info("Found target")

	compile(client, rpcInstance, targetBoard, sketch)
	upload(client, rpcInstance, targetBoard, sketch)

	watchLogs(targetBoard.Device, baud)
}

func main() {
	var baud int
	rootCmd := &cobra.Command{
		Use:   "ardi [sketch]",
		Short: "Ardi uploads sketches and prints logs for a variety of arduino boards.",
		Long: "A light wrapper around arduino-cli that offers a quick way to upload\n" +
			"sketches and watch logs from command line for a variety of arduino boards.",
		Run: func(cmd *cobra.Command, args []string) {
			process(baud)
		},
	}

	initCommand := &cobra.Command{
		Use:   "init",
		Short: "Download and install platforms",
		Long: "Downloads and installs all available platforms to support\n" +
			"a maximum number of boards.",
		Run: func(cmd *cobra.Command, args []string) {
			conn, _, _ := initialize()
			defer conn.Close()
		},
	}

	cleanCommand := &cobra.Command{
		Use:   "clean",
		Short: "Delete all ardi data",
		Long:  "Removes all installed platforms from ~/.ardi",
		Run: func(cmd *cobra.Command, args []string) {
			if err := os.RemoveAll(ardiDir); err != nil {
				logger.WithError(err).Fatalf("Failed to clean ardi directory. You can manually clean all data by removing %s", ardiDir)
			}
			logger.Infof("Successfully removed all data from %s", ardiDir)
		},
	}

	rootCmd.Flags().IntVarP(&baud, "baud", "b", 9600, "specify sketch baud rate")
	rootCmd.AddCommand(initCommand)
	rootCmd.AddCommand(cleanCommand)
	rootCmd.Execute()
}
