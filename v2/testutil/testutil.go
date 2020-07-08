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
	"github.com/spf13/cobra"

	"github.com/robgonnella/ardi/v2/commands"
	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/mocks"
)

var port = 2222
var here string
var userHome string

// GlobalOpt option to make a command act globally
type GlobalOpt struct {
	Global bool
}

func init() {
	here, _ = filepath.Abs(".")
	userHome, _ = os.UserHomeDir()
	logrus.SetOutput(ioutil.Discard)
}

func cleanCoreDir() {
	dataDir := path.Join(here, "../core/.ardi")
	jsonFile := path.Join(here, "../core/ardi.json")
	os.RemoveAll(dataDir)
	os.Remove(jsonFile)
}

func cleanCommandsDir() {
	projectDataDir := path.Join(here, "../commands/.ardi")
	projectJSONFile := path.Join(here, "../commands/ardi.json")
	os.RemoveAll(projectDataDir)
	os.Remove(projectJSONFile)
}

func cleanGlobalData() {
	globalDataDir := path.Join(userHome, ".ardi")
	os.RemoveAll(globalDataDir)
}

func cleanBuilds() {
	os.RemoveAll(path.Join(BlinkProjectDir(), "build"))
	os.RemoveAll(path.Join(PixieProjectDir(), "build"))
}

func cleanAll() {
	cleanCoreDir()
	cleanCommandsDir()
	cleanGlobalData()
	cleanBuilds()
}

// ArduinoMegaFQBN returns appropriate fqbn for arduino mega 2560
func ArduinoMegaFQBN() string {
	return "arduino:avr:mega"
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

// GenerateCmdBoard returns a single rpc Board
func GenerateCmdBoard(name, fqbn string) *rpccommands.Board {
	if fqbn == "" {
		fqbn = fmt.Sprintf("%s-fqbn", name)
	}
	return &rpccommands.Board{Name: name, Fqbn: fqbn}
}

// BlinkProjectDir returns path to blink project directory
func BlinkProjectDir() string {
	here, _ := filepath.Abs(".")
	return path.Join(here, "../test_projects/blink")
}

// PixieProjectDir returns path to blink project directory
func PixieProjectDir() string {
	here, _ := filepath.Abs(".")
	return path.Join(here, "../test_projects/pixie")
}

// GenerateCmdBoards generate a list of boards
func GenerateCmdBoards(n int) []*rpccommands.Board {
	var boards []*rpccommands.Board
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("test-board-%02d", i)
		b := GenerateCmdBoard(name, "")
		boards = append(boards, b)
	}
	return boards
}

// GenerateCmdPlatform generates a single named platform
func GenerateCmdPlatform(name string, boards []*rpccommands.Board) *rpccommands.Platform {
	return &rpccommands.Platform{
		ID:     name,
		Boards: boards,
	}
}

// RunUnitTest runs an ardi unit test
func RunUnitTest(name string, t *testing.T, f func(env UnitTestEnv)) {
	t.Run(name, func(st *testing.T) {
		ctrl := gomock.NewController(t)
		client := mocks.NewMockClient(ctrl)
		logger := log.New()

		cleanAll()

		var b bytes.Buffer
		logger.SetOutput(&b)
		logger.SetLevel(log.DebugLevel)

		ardiCore := core.NewArdiCore(client, logger)

		env := UnitTestEnv{
			T:        st,
			Ctrl:     ctrl,
			Logger:   logger,
			Client:   client,
			ArdiCore: ardiCore,
			Stdout:   &b,
		}

		f(env)
		cleanAll()
	})
}

// IntegrationTestEnv represents our integration test environment
type IntegrationTestEnv struct {
	Ctx     context.Context
	T       *testing.T
	Logger  *log.Logger
	RootCmd *cobra.Command
	Port    int
	SetArgs func(a []string)
	Stdout  *bytes.Buffer
}

// RunIntegrationTest runs an ardi integration test
func RunIntegrationTest(name string, t *testing.T, f func(env *IntegrationTestEnv)) {
	t.Run(name, func(st *testing.T) {
		port = port + 1
		cleanAll()

		ctx := context.Background()
		var b bytes.Buffer
		logger := log.New()
		logger.Out = &b
		logger.SetLevel(log.DebugLevel)

		rootCmd := commands.GetRootCmd(logger)
		rootCmd.SetOut(logger.Out)

		env := IntegrationTestEnv{
			Ctx:     ctx,
			T:       st,
			Logger:  logger,
			RootCmd: rootCmd,
			Port:    port,
			SetArgs: func(args []string) {
				args = append(args, "--verbose")
				args = append(args, "--port")
				args = append(args, fmt.Sprintf("%d", port))
				rootCmd.SetArgs(args)
			},
			Stdout: &b,
		}

		f(&env)
		cleanAll()
	})
}

// InstallAvrPlatform uses ardi command to install aruidno:avr platform
func (e *IntegrationTestEnv) InstallAvrPlatform(opt GlobalOpt) error {
	if !opt.Global {
		projectArgs := []string{"project", "init"}
		e.SetArgs(projectArgs)
		if err := e.RootCmd.ExecuteContext(e.Ctx); err != nil {
			return err
		}
	}
	platformArgs := []string{"platform", "add", "arduino:avr"}
	if opt.Global {
		platformArgs = append(platformArgs, "--global")
	}
	e.SetArgs(platformArgs)
	return e.RootCmd.ExecuteContext(e.Ctx)
}

// RunProjectInit initializes and ardi project directory
func (e *IntegrationTestEnv) RunProjectInit() error {
	projectInitArgs := []string{"project", "init"}
	e.SetArgs(projectInitArgs)
	return e.RootCmd.ExecuteContext(e.Ctx)
}

// AddLib adds an arduino library
func (e *IntegrationTestEnv) AddLib(lib string, opt GlobalOpt) error {
	args := []string{"lib", "add", lib}
	if opt.Global {
		args = append(args, "--global")
	}
	e.SetArgs(args)
	return e.RootCmd.ExecuteContext(e.Ctx)
}

// AddPlatform adds an arduino platform
func (e *IntegrationTestEnv) AddPlatform(platform string, opt GlobalOpt) error {
	args := []string{"platform", "add", platform}
	if opt.Global {
		args = append(args, "--global")
	}
	e.SetArgs(args)
	return e.RootCmd.ExecuteContext(e.Ctx)
}
