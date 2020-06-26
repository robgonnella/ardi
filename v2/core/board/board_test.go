package board_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/arduino/arduino-cli/rpc/commands"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/robgonnella/ardi/v2/core/board"
	"github.com/robgonnella/ardi/v2/mocks"
)

type TestEnv struct {
	ctrl      *gomock.Controller
	logger    *log.Logger
	client    *mocks.MockClient
	boardCore *board.Board
	boards    []*commands.Board
	platforms []*commands.Platform
	stdout    *bytes.Buffer
}

func runTest(name string, t *testing.T, f func(t *testing.T, env TestEnv)) {
	t.Run(name, func(st *testing.T) {
		ctrl := gomock.NewController(t)
		client := mocks.NewMockClient(ctrl)
		logger := log.New()

		var b bytes.Buffer
		logger.SetOutput(&b)
		logger.SetLevel(log.DebugLevel)

		var boards []*commands.Board
		for i := 1; i <= 10; i++ {
			b := &commands.Board{
				Name: fmt.Sprintf("test-board-%d", i),
				Fqbn: fmt.Sprintf("board-fqbn-%d", i),
			}
			boards = append(boards, b)
		}

		platform := &commands.Platform{Boards: boards, ID: "test-platform"}
		platforms := []*commands.Platform{platform}

		env := TestEnv{
			ctrl:      ctrl,
			logger:    logger,
			client:    client,
			boardCore: board.New(client, logger),
			boards:    boards,
			platforms: platforms,
			stdout:    &b,
		}

		f(st, env)
	})
}

func TestBoardCore(t *testing.T) {
	runTest("prints fqbns", t, func(st *testing.T, env TestEnv) {
		defer env.ctrl.Finish()
		env.client.EXPECT().GetPlatforms().Return(env.platforms, nil)
		env.boardCore.FQBNS("")
		for _, b := range env.boards {
			assert.Contains(t, env.stdout.String(), b.GetName())
			assert.Contains(t, env.stdout.String(), b.GetFqbn())
		}
	})

	runTest("prints platforms", t, func(st *testing.T, env TestEnv) {
		defer env.ctrl.Finish()
		env.client.EXPECT().GetPlatforms().Return(env.platforms, nil)
		env.boardCore.Platforms("")
		platform := env.platforms[0]
		for _, b := range env.boards {
			assert.Contains(t, env.stdout.String(), b.GetName())
			assert.Contains(t, env.stdout.String(), platform.GetID())
		}
	})
}
