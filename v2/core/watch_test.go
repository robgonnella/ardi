package core_test

import (
	"path"
	"testing"

	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

// @todo: check that list is actually sorted
func TestWatchCore(t *testing.T) {
	testutil.RunUnitTest("returns error if no boards connected", t, func(env testutil.UnitTestEnv) {
		defer env.Ctrl.Finish()
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(env.T, err)

		port := "/dev/null"
		sketchDir := env.BlinkProjDir
		props := []string{}
		connectedBoards := []*rpc.Board{}
		allBoards := []*rpc.Board{}

		env.Client.EXPECT().ConnectedBoards().Times(1).Return(connectedBoards)
		env.Client.EXPECT().AllBoards().Times(1).Return(allBoards)

		err = env.ArdiCore.Watch.Init(port, sketchDir, props)
		assert.Error(env.T, err)
	})

	testutil.RunUnitTest("succeeds when board is connected", t, func(env testutil.UnitTestEnv) {
		defer env.Ctrl.Finish()
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(env.T, err)

		board := rpc.Board{
			FQBN: "some-fqbn",
			Name: "board-name",
			Port: "/dev/null",
		}
		sketchDir := env.BlinkProjDir
		props := []string{}
		connectedBoards := []*rpc.Board{&board}
		allBoards := []*rpc.Board{}

		env.Client.EXPECT().ConnectedBoards().Times(1).Return(connectedBoards)
		env.Client.EXPECT().AllBoards().Times(1).Return(allBoards)

		err = env.ArdiCore.Watch.Init(board.Port, sketchDir, props)
		assert.NoError(env.T, err)
	})

	testutil.RunUnitTest("errors if not a valid sketch directory", t, func(env testutil.UnitTestEnv) {
		defer env.Ctrl.Finish()
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(env.T, err)

		board := rpc.Board{
			FQBN: "some-fqbn",
			Name: "board-name",
			Port: "/dev/null",
		}
		sketchDir := "."
		props := []string{}

		err = env.ArdiCore.Watch.Init(board.Port, sketchDir, props)
		assert.Error(env.T, err)
	})

	testutil.RunUnitTest("compiles sketch", t, func(env testutil.UnitTestEnv) {
		defer env.Ctrl.Finish()
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(env.T, err)

		fqbn := "some-board-fqbn"
		boardName := "some-board-name"
		port := "/dev/null"
		board := rpc.Board{
			FQBN: fqbn,
			Name: boardName,
			Port: port,
		}

		props := []string{"somebuild_prop=DTest"}
		sketchPath := path.Join(env.BlinkProjDir, "blink.ino")
		showProps := false
		exportName := ""
		compileOpts := rpc.CompileOpts{
			FQBN:       fqbn,
			BuildProps: props,
			ShowProps:  showProps,
			SketchDir:  env.BlinkProjDir,
			ExportName: exportName,
			SketchPath: sketchPath,
		}

		connectedBoards := []*rpc.Board{&board}
		allBoards := []*rpc.Board{}

		env.Client.EXPECT().ConnectedBoards().Times(1).Return(connectedBoards)
		env.Client.EXPECT().AllBoards().Times(1).Return(allBoards)
		env.ArdiCore.Watch.Init(board.Port, env.BlinkProjDir, props)

		env.Client.EXPECT().Compile(compileOpts).Times(1).Return(nil)

		err = env.ArdiCore.Watch.Compile()
		assert.NoError(env.T, err)
	})

	testutil.RunUnitTest("uploads sketch", t, func(env testutil.UnitTestEnv) {
		defer env.Ctrl.Finish()
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(env.T, err)

		fqbn := "some-board-fqbn"
		boardName := "some-board-name"
		port := "/dev/null"
		board := rpc.Board{
			FQBN: fqbn,
			Name: boardName,
			Port: port,
		}
		props := []string{}
		connectedBoards := []*rpc.Board{&board}
		allBoards := []*rpc.Board{}

		env.Client.EXPECT().ConnectedBoards().Times(1).Return(connectedBoards)
		env.Client.EXPECT().AllBoards().Times(1).Return(allBoards)

		env.ArdiCore.Watch.Init(board.Port, env.BlinkProjDir, props)

		env.Client.EXPECT().Upload(fqbn, env.BlinkProjDir, board.Port).Times(1).Return(nil)

		err = env.ArdiCore.Watch.Upload()
		assert.NoError(env.T, err)
	})
}
