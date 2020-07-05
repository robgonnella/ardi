package testutil

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	rpccommands "github.com/arduino/arduino-cli/rpc/commands"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/robgonnella/ardi/v2/commands"
	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/mocks"
)

func cleanCoreDir() {
	here, _ := filepath.Abs(".")
	dataDir := path.Join(here, "../core/.ardi")
	jsonFile := path.Join(here, "../core/ardi.json")
	os.RemoveAll(dataDir)
	os.Remove(jsonFile)
}

// UnitTestEnv represents our unit test environment
type UnitTestEnv struct {
	T            *testing.T
	Ctrl         *gomock.Controller
	Logger       *log.Logger
	Client       *mocks.MockClient
	ArdiCore     *core.ArdiCore
	Stdout       *bytes.Buffer
	BlinkProjDir string
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
		here, _ := filepath.Abs(".")

		cleanCoreDir()

		var b bytes.Buffer
		logger.SetOutput(&b)
		logger.SetLevel(log.DebugLevel)

		ardiCore := core.NewArdiCore(client, logger)

		env := UnitTestEnv{
			T:            st,
			Ctrl:         ctrl,
			Logger:       logger,
			Client:       client,
			ArdiCore:     ardiCore,
			Stdout:       &b,
			BlinkProjDir: path.Join(here, "../test_projects/blink"),
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
	SetArgs func(a []string)
	Stdout  *bytes.Buffer
}

// RunIntegrationTest runs an ardi integration test
func RunIntegrationTest(name string, t *testing.T, f func(env IntegrationTestEnv)) {
	t.Run(name, func(st *testing.T) {
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
			Stdout:  &b,
		}

		f(env)
	})
}
