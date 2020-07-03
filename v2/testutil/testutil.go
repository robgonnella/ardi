package testutil

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/arduino/arduino-cli/rpc/commands"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"

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

// TestEnv represents our test environment
type TestEnv struct {
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
func GenerateCmdBoard(name, fqbn string) *commands.Board {
	if fqbn == "" {
		fqbn = fmt.Sprintf("%s-fqbn", name)
	}
	return &commands.Board{Name: name, Fqbn: fqbn}
}

// GenerateCmdBoards generate a list of boards
func GenerateCmdBoards(n int) []*commands.Board {
	var boards []*commands.Board
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("test-board-%02d", i)
		b := GenerateCmdBoard(name, "")
		boards = append(boards, b)
	}
	return boards
}

// GenerateCmdPlatform generates a single named platform
func GenerateCmdPlatform(name string, boards []*commands.Board) *commands.Platform {
	return &commands.Platform{
		ID:     name,
		Boards: boards,
	}
}

// RunTest runs an ardi unit test
func RunTest(name string, t *testing.T, f func(t *testing.T, env TestEnv)) {
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

		env := TestEnv{
			Ctrl:         ctrl,
			Logger:       logger,
			Client:       client,
			ArdiCore:     ardiCore,
			Stdout:       &b,
			BlinkProjDir: path.Join(here, "../test_projects/blink"),
		}

		f(st, env)
		cleanCoreDir()
	})
}
