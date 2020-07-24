package core_test

import (
	"path"
	"testing"

	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestArdiCore(t *testing.T) {
	testutil.RunUnitTest("returns target if fqbn provided and onlyConnected false", t, func(env *testutil.UnitTestEnv) {
		connectedBoards := []*rpc.Board{}
		allBoards := []*rpc.Board{}
		fqbn := "someboardfqbn"

		env.Client.EXPECT().ConnectedBoards().Times(1).Return(connectedBoards)
		env.Client.EXPECT().AllBoards().Times(1).Return(allBoards)

		board, err := env.ArdiCore.GetTargetBoard(fqbn, false)
		assert.NoError(env.T, err)
		assert.Equal(env.T, board.FQBN, fqbn)
	})

	testutil.RunUnitTest("returns target if fqbn provided and onlyConnected true", t, func(env *testutil.UnitTestEnv) {
		boardName := "someboardname"
		fqbn := "someboardfqbn"
		connectedBoard := testutil.GenerateRPCBoard(boardName, fqbn)
		connectedBoards := []*rpc.Board{connectedBoard}
		allBoards := []*rpc.Board{}

		env.Client.EXPECT().ConnectedBoards().Times(1).Return(connectedBoards)
		env.Client.EXPECT().AllBoards().Times(1).Return(allBoards)

		board, err := env.ArdiCore.GetTargetBoard(fqbn, true)
		assert.NoError(env.T, err)
		assert.Equal(env.T, board.FQBN, fqbn)
	})

	testutil.RunUnitTest("errors if fqbn provided and onlyConnected true and board not connected", t, func(env *testutil.UnitTestEnv) {
		fqbn := "someboardfqbn"
		connectedBoards := []*rpc.Board{}
		allBoards := []*rpc.Board{}

		env.Client.EXPECT().ConnectedBoards().Times(1).Return(connectedBoards)
		env.Client.EXPECT().AllBoards().Times(1).Return(allBoards)

		board, err := env.ArdiCore.GetTargetBoard(fqbn, true)
		assert.Error(env.T, err)
		assert.Nil(env.T, board)
	})

	testutil.RunUnitTest("returns target if 1 connected board", t, func(env *testutil.UnitTestEnv) {
		boardName := "somboardname"
		boardFQBN := "someboardfqbn"
		connectedBoard := testutil.GenerateRPCBoard(boardName, boardFQBN)
		connectedBoards := []*rpc.Board{connectedBoard}
		allBoards := []*rpc.Board{}

		env.Client.EXPECT().ConnectedBoards().Times(1).Return(connectedBoards)
		env.Client.EXPECT().AllBoards().Times(1).Return(allBoards)

		board, err := env.ArdiCore.GetTargetBoard("", false)
		assert.NoError(env.T, err)
		assert.Equal(env.T, board.Name, boardName)
		assert.Equal(env.T, board.FQBN, boardFQBN)
	})

	testutil.RunUnitTest("returns error and prints connected boards if more than one found", t, func(env *testutil.UnitTestEnv) {
		board1Name := "somboardname"
		board1FQBN := "someboardfqbn"
		board2Name := "someotherboardname"
		board2FQBN := "someotherboardfqbn"
		connectedBoard1 := testutil.GenerateRPCBoard(board1Name, board1FQBN)
		connectedBoard2 := testutil.GenerateRPCBoard(board2Name, board2FQBN)
		connectedBoards := []*rpc.Board{connectedBoard1, connectedBoard2}
		allBoards := []*rpc.Board{}

		env.Client.EXPECT().ConnectedBoards().Times(1).Return(connectedBoards)
		env.Client.EXPECT().AllBoards().Times(1).Return(allBoards)

		env.ClearStdout()
		board, err := env.ArdiCore.GetTargetBoard("", false)
		assert.Error(env.T, err)
		assert.Nil(env.T, board)

		out := env.Stdout.String()
		assert.Contains(env.T, out, board1Name)
		assert.Contains(env.T, out, board1FQBN)
		assert.Contains(env.T, out, board2Name)
		assert.Contains(env.T, out, board2FQBN)
	})

	testutil.RunUnitTest("returns error and prints all available boards if no connected boards found", t, func(env *testutil.UnitTestEnv) {
		otherBoardName := "someotherboardname"
		otherBoardFQBN := "someotherboardfqbn"
		connectedBoards := []*rpc.Board{}

		otherBoard := testutil.GenerateRPCBoard(otherBoardName, otherBoardFQBN)
		allBoards := []*rpc.Board{otherBoard}

		env.Client.EXPECT().ConnectedBoards().Times(1).Return(connectedBoards)
		env.Client.EXPECT().AllBoards().Times(1).Return(allBoards)

		env.ClearStdout()
		board, err := env.ArdiCore.GetTargetBoard("", false)
		assert.Error(env.T, err)
		assert.Nil(env.T, board)

		out := env.Stdout.String()
		assert.Contains(env.T, out, otherBoardName)
		assert.Contains(env.T, out, otherBoardFQBN)
	})

	testutil.RunUnitTest("returns error and does not print available boards if only connected specified", t, func(env *testutil.UnitTestEnv) {
		connectedBoards := []*rpc.Board{}

		otherBoardName := "someotherboardname"
		otherBoardFQBN := "someotherboardfqbn"
		otherBoard := testutil.GenerateRPCBoard(otherBoardName, otherBoardFQBN)
		allBoards := []*rpc.Board{otherBoard}

		env.Client.EXPECT().ConnectedBoards().Times(1).Return(connectedBoards)
		env.Client.EXPECT().AllBoards().Times(1).Return(allBoards)

		env.ClearStdout()
		board, err := env.ArdiCore.GetTargetBoard("", true)
		assert.Error(env.T, err)
		assert.Nil(env.T, board)

		out := env.Stdout.String()
		assert.NotContains(env.T, out, otherBoardName)
		assert.NotContains(env.T, out, otherBoardFQBN)
	})

	testutil.RunUnitTest("compiles ardi build", t, func(env *testutil.UnitTestEnv) {
		connectedBoards := []*rpc.Board{}
		allBoards := []*rpc.Board{}
		buildName := "somebuild"
		sketch := path.Join(testutil.BlinkProjectDir(), "blink.ino")
		fqbn := "someboardfqbn"

		err := env.ArdiCore.Config.AddBuild(buildName, sketch, fqbn, []string{})
		assert.NoError(env.T, err)

		env.Client.EXPECT().ConnectedBoards().Times(1).Return(connectedBoards)
		env.Client.EXPECT().AllBoards().Times(1).Return(allBoards)

		buildOpts := core.CompileArdiBuildOpts{
			BuildName:           "somebuild",
			OnlyConnectedBoards: false,
		}

		expectedCompileOpts := rpc.CompileOpts{
			ExportName: buildName,
			FQBN:       fqbn,
			SketchDir:  testutil.BlinkProjectDir(),
			SketchPath: sketch,
			BuildProps: []string{},
			ShowProps:  false,
		}

		env.Client.EXPECT().Compile(expectedCompileOpts).Times(1).Return(nil)

		compileOpts, board, err := env.ArdiCore.CompileArdiBuild(buildOpts)
		assert.NoError(env.T, err)
		assert.Equal(env.T, board.FQBN, fqbn)
		assert.Equal(env.T, &expectedCompileOpts, compileOpts)
	})

	testutil.RunUnitTest("errors compiling ardi build when onlyConnectedBoards is true", t, func(env *testutil.UnitTestEnv) {
		connectedBoards := []*rpc.Board{}
		allBoards := []*rpc.Board{}
		buildName := "somebuild"
		sketch := path.Join(testutil.BlinkProjectDir(), "blink.ino")
		fqbn := "someboardfqbn"

		err := env.ArdiCore.Config.AddBuild(buildName, sketch, fqbn, []string{})
		assert.NoError(env.T, err)

		env.Client.EXPECT().ConnectedBoards().Times(1).Return(connectedBoards)
		env.Client.EXPECT().AllBoards().Times(1).Return(allBoards)

		buildOpts := core.CompileArdiBuildOpts{
			BuildName:           "somebuild",
			OnlyConnectedBoards: true,
		}

		compileOpts, board, err := env.ArdiCore.CompileArdiBuild(buildOpts)
		assert.Error(env.T, err)
		assert.Nil(env.T, compileOpts)
		assert.Nil(env.T, board)
	})

	testutil.RunUnitTest("compiles sketch", t, func(env *testutil.UnitTestEnv) {
		connectedBoards := []*rpc.Board{}
		allBoards := []*rpc.Board{}
		buildName := "somebuild"
		sketch := path.Join(testutil.BlinkProjectDir(), "blink.ino")
		fqbn := "someboardfqbn"

		err := env.ArdiCore.Config.AddBuild(buildName, sketch, fqbn, []string{})
		assert.NoError(env.T, err)

		env.Client.EXPECT().ConnectedBoards().Times(1).Return(connectedBoards)
		env.Client.EXPECT().AllBoards().Times(1).Return(allBoards)

		sketchOpts := core.CompileSketchOpts{
			Sketch:              sketch,
			FQBN:                fqbn,
			BuildPros:           []string{},
			ShowProps:           false,
			OnlyConnectedBoards: false,
		}

		expectedCompileOpts := rpc.CompileOpts{
			ExportName: "",
			FQBN:       fqbn,
			SketchDir:  testutil.BlinkProjectDir(),
			SketchPath: sketch,
			BuildProps: []string{},
			ShowProps:  false,
		}

		env.Client.EXPECT().Compile(expectedCompileOpts).Times(1).Return(nil)

		compileOpts, board, err := env.ArdiCore.CompileSketch(sketchOpts)
		assert.NoError(env.T, err)
		assert.Equal(env.T, board.FQBN, fqbn)
		assert.Equal(env.T, &expectedCompileOpts, compileOpts)
	})

	testutil.RunUnitTest("errors compiling sketch when onlyConnectedBoards is true", t, func(env *testutil.UnitTestEnv) {
		connectedBoards := []*rpc.Board{}
		allBoards := []*rpc.Board{}
		buildName := "somebuild"
		sketch := path.Join(testutil.BlinkProjectDir(), "blink.ino")
		fqbn := "someboardfqbn"

		err := env.ArdiCore.Config.AddBuild(buildName, sketch, fqbn, []string{})
		assert.NoError(env.T, err)

		env.Client.EXPECT().ConnectedBoards().Times(1).Return(connectedBoards)
		env.Client.EXPECT().AllBoards().Times(1).Return(allBoards)

		sketchOpts := core.CompileSketchOpts{
			Sketch:              sketch,
			FQBN:                fqbn,
			BuildPros:           []string{},
			ShowProps:           false,
			OnlyConnectedBoards: true,
		}

		compileOpts, board, err := env.ArdiCore.CompileSketch(sketchOpts)
		assert.Error(env.T, err)
		assert.Nil(env.T, compileOpts)
		assert.Nil(env.T, board)
	})
}
