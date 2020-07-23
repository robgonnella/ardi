package testutil

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	rpccommands "github.com/arduino/arduino-cli/rpc/commands"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/robgonnella/ardi/v2/commands"
	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/mocks"
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/util"
)

var port = 3000
var here string
var userHome string

func init() {
	here, _ = filepath.Abs(".")
	userHome, _ = os.UserHomeDir()
	logrus.SetOutput(ioutil.Discard)
}

// CleanCoreDir removes test data from core directory
func CleanCoreDir() {
	dataDir := path.Join(here, "../core/.ardi")
	jsonFile := path.Join(here, "../core/ardi.json")
	os.RemoveAll(dataDir)
	os.Remove(jsonFile)
}

// CleanCommandsDir removes project data from commands directory
func CleanCommandsDir() {
	projectDataDir := path.Join(here, "../commands/.ardi")
	projectJSONFile := path.Join(here, "../commands/ardi.json")
	os.RemoveAll(projectDataDir)
	os.Remove(projectJSONFile)
}

// CleanGlobalData removes global data directory
func CleanGlobalData() {
	globalDataDir := path.Join(userHome, ".ardi")
	os.RemoveAll(globalDataDir)
}

// CleanBuilds removes compiled test project builds
func CleanBuilds() {
	os.RemoveAll(path.Join(BlinkProjectDir(), "build"))
	os.RemoveAll(path.Join(PixieProjectDir(), "build"))
}

// CleanAll removes all test data
func CleanAll() {
	CleanCoreDir()
	CleanCommandsDir()
	CleanGlobalData()
	CleanBuilds()
}

// ArduinoMegaFQBN returns appropriate fqbn for arduino mega 2560
func ArduinoMegaFQBN() string {
	return "arduino:avr:mega"
}

// Esp8266Platform returns appropriate platform for esp8266
func Esp8266Platform() string {
	return "esp8266:esp8266"
}

// Esp8266WifiduinoFQBN returns appropriate fqbn for esp8266 board
func Esp8266WifiduinoFQBN() string {
	return "esp8266:esp8266:wifiduino"
}

// Esp8266BoardURL returns appropriate board url for esp8266 board
func Esp8266BoardURL() string {
	return "https://arduino.esp8266.com/stable/package_esp8266com_index.json"
}

// UnitTestEnv represents our unit test environment
type UnitTestEnv struct {
	T            *testing.T
	Ctrl         *gomock.Controller
	Logger       *log.Logger
	Client       *mocks.MockClient
	ArdiCore     *core.ArdiCore
	Stdout       *bytes.Buffer
	PixieProjDir string
	EmptyProjDIr string
}

// ClearStdout clears stdout for unit test
func (e *UnitTestEnv) ClearStdout() {
	var b bytes.Buffer
	e.Logger.SetOutput(&b)
	e.Stdout = &b
}

// GenerateCmdBoard returns a single arduino-cli command Board
func GenerateCmdBoard(name, fqbn string) *rpccommands.Board {
	if fqbn == "" {
		fqbn = fmt.Sprintf("%s-fqbn", name)
	}
	return &rpccommands.Board{Name: name, Fqbn: fqbn}
}

// GenerateCmdBoards generate a list of arduino-cli command boards
func GenerateCmdBoards(n int) []*rpccommands.Board {
	var boards []*rpccommands.Board
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("test-board-%02d", i)
		b := GenerateCmdBoard(name, "")
		boards = append(boards, b)
	}
	return boards
}

// GenerateCmdPlatform generates a single named arduino-cli command platform
func GenerateCmdPlatform(name string, boards []*rpccommands.Board) *rpccommands.Platform {
	return &rpccommands.Platform{
		ID:     name,
		Boards: boards,
	}
}

// GenerateRPCBoard returns a single ardi rpc Board
func GenerateRPCBoard(name, fqbn string) *rpc.Board {
	if fqbn == "" {
		fqbn = fmt.Sprintf("%s-fqbn", name)
	}
	return &rpc.Board{
		Name: name,
		FQBN: fqbn,
		Port: "/dev/null",
	}
}

// GenerateRPCBoards generate a list of ardi rpc boards
func GenerateRPCBoards(n int) []*rpc.Board {
	var boards []*rpc.Board
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("test-board-%02d", i)
		b := GenerateRPCBoard(name, "")
		boards = append(boards, b)
	}
	return boards
}

// BlinkProjectDir returns path to blink project directory
func BlinkProjectDir() string {
	return path.Join(here, "../test_projects/blink")
}

// BlinkCopyProjectDir returns path to blink project directory
func BlinkCopyProjectDir() string {
	return path.Join(here, "../test_projects/blink2")
}

// Blink14400ProjectDir returns path to blink14400 project directory
func Blink14400ProjectDir() string {
	return path.Join(here, "../test_projects/blink14400")
}

// PixieProjectDir returns path to blink project directory
func PixieProjectDir() string {
	return path.Join(here, "../test_projects/pixie")
}

// RunUnitTest runs an ardi unit test
func RunUnitTest(name string, t *testing.T, f func(env *UnitTestEnv)) {
	t.Run(name, func(st *testing.T) {
		ctrl := gomock.NewController(st)
		defer ctrl.Finish()

		client := mocks.NewMockClient(ctrl)
		logger := log.New()

		CleanAll()

		var b bytes.Buffer
		logger.SetOutput(&b)
		logger.SetLevel(log.DebugLevel)

		opts := util.GetAllSettingsOpts{
			Global:   false,
			LogLevel: "debug",
			Port:     "2222",
		}
		ardiConfig, svrSettings := util.GetAllSettings(opts)

		coreOpts := core.NewArdiCoreOpts{
			Global:             false,
			Logger:             logger,
			Client:             client,
			ArdiConfig:         *ardiConfig,
			ArduinoCliSettings: *svrSettings,
		}
		ardiCore := core.NewArdiCore(coreOpts)

		env := UnitTestEnv{
			T:        st,
			Ctrl:     ctrl,
			Logger:   logger,
			Client:   client,
			ArdiCore: ardiCore,
			Stdout:   &b,
		}

		f(&env)
		CleanAll()
	})
}

// IntegrationTestEnv represents our integration test environment
type IntegrationTestEnv struct {
	T      *testing.T
	Stdout *bytes.Buffer
	logger *log.Logger
}

// RunIntegrationTest runs an ardi integration test
func RunIntegrationTest(name string, t *testing.T, f func(env *IntegrationTestEnv)) {
	t.Run(name, func(st *testing.T) {
		CleanAll()

		var b bytes.Buffer
		logger := log.New()
		logger.Out = &b
		logger.SetLevel(log.DebugLevel)

		env := IntegrationTestEnv{
			T:      st,
			Stdout: &b,
			logger: logger,
		}

		f(&env)
		CleanAll()
	})
}

// RunProjectInit initializes and ardi project directory
func (e *IntegrationTestEnv) RunProjectInit() error {
	projectInitArgs := []string{"project-init"}
	return e.Execute(projectInitArgs)
}

// Execute executes the root command with given arguments
func (e *IntegrationTestEnv) Execute(args []string) error {
	rootCmd := commands.GetRootCmd(e.logger)
	rootCmd.SetOut(e.logger.Out)

	port = port + 1
	args = append(args, "--verbose")
	args = append(args, "--port")
	args = append(args, fmt.Sprintf("%d", port))
	rootCmd.SetArgs(args)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return rootCmd.ExecuteContext(ctx)
}

// ClearStdout clears integration test env stdout
func (e *IntegrationTestEnv) ClearStdout() {
	var b bytes.Buffer
	e.logger.SetOutput(&b)
	e.Stdout = &b
}
