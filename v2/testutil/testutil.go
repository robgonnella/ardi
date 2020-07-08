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

var port = 2221

func init() {
	logrus.SetOutput(ioutil.Discard)
}

func cleanCoreDir() {
	here, _ := filepath.Abs(".")
	dataDir := path.Join(here, "../core/.ardi")
	jsonFile := path.Join(here, "../core/ardi.json")
	os.RemoveAll(dataDir)
	os.Remove(jsonFile)
}

func cleanCommandsDir() {
	here, _ := filepath.Abs(".")
	dataDir := path.Join(here, "../commands/.ardi")
	jsonFile := path.Join(here, "../commands/ardi.json")
	os.RemoveAll(dataDir)
	os.Remove(jsonFile)
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

		cleanCoreDir()

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
		cleanCoreDir()
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
		cleanCommandsDir()

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
		cleanCommandsDir()
	})
}

// InstallAvrPlatform uses ardi command to install aruidno:avr platform
func (e *IntegrationTestEnv) InstallAvrPlatform() error {
	projectArgs := []string{"project", "init"}
	e.SetArgs(projectArgs)
	if err := e.RootCmd.ExecuteContext(e.Ctx); err != nil {
		return err
	}
	platformArgs := []string{"platform", "add", "arduino:avr"}
	e.SetArgs(platformArgs)
	if err := e.RootCmd.ExecuteContext(e.Ctx); err != nil {
		return err
	}
	return nil
}

// RunProjectInit initializes and ardi project directory
func (e *IntegrationTestEnv) RunProjectInit() error {
	projectInitArgs := []string{"project", "init"}
	e.SetArgs(projectInitArgs)
	return e.RootCmd.ExecuteContext(e.Ctx)
}

// AddLib adds an arduino library
func (e *IntegrationTestEnv) AddLib(lib string, global bool) error {
	args := []string{"lib", "add", lib}
	if global {
		args = append(args, "--global")
	}
	e.SetArgs(args)
	return e.RootCmd.ExecuteContext(e.Ctx)
}

// AddPlatform adds an arduino platform
func (e *IntegrationTestEnv) AddPlatform(platform string, global bool) error {
	args := []string{"platform", "add", platform}
	if global {
		args = append(args, "--global")
	}
	e.SetArgs(args)
	return e.RootCmd.ExecuteContext(e.Ctx)
}
