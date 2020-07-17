package core_test

import (
	"testing"

	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestTarget(t *testing.T) {
	testutil.RunUnitTest("returns target if fqbn provided", t, func(env *testutil.UnitTestEnv) {
		connectedBoards := []*rpc.Board{}
		allBoards := []*rpc.Board{}
		fqbn := "someboardfqbn"

		target, err := core.NewTarget(connectedBoards, allBoards, fqbn, false, env.Logger)
		assert.NoError(env.T, err)
		assert.Equal(env.T, target.Board.FQBN, fqbn)
	})

	testutil.RunUnitTest("returns target if 1 connected board", t, func(env *testutil.UnitTestEnv) {
		boardName := "somboardname"
		boardFQBN := "someboardfqbn"
		connectedBoard := testutil.GenerateRPCBoard(boardName, boardFQBN)
		connectedBoards := []*rpc.Board{connectedBoard}
		allBoards := []*rpc.Board{}

		target, err := core.NewTarget(connectedBoards, allBoards, "", false, env.Logger)
		assert.NoError(env.T, err)
		assert.Equal(env.T, target.Board.Name, boardName)
		assert.Equal(env.T, target.Board.FQBN, boardFQBN)
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

		env.ClearStdout()
		target, err := core.NewTarget(connectedBoards, allBoards, "", false, env.Logger)
		assert.Error(env.T, err)
		assert.Nil(env.T, target)
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

		env.ClearStdout()
		target, err := core.NewTarget(connectedBoards, allBoards, "", false, env.Logger)
		assert.Error(env.T, err)
		assert.Nil(env.T, target)
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

		env.ClearStdout()
		target, err := core.NewTarget(connectedBoards, allBoards, "", true, env.Logger)
		assert.Error(env.T, err)
		assert.Nil(env.T, target)
		out := env.Stdout.String()
		assert.NotContains(env.T, out, otherBoardName)
		assert.NotContains(env.T, out, otherBoardFQBN)
	})
}
