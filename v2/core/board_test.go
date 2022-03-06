package core_test

import (
	"errors"
	"testing"

	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/golang/mock/gomock"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestBoardCore(t *testing.T) {
	testutil.RunUnitTest("returns target if fqbn provided and onlyConnected false", t, func(env *testutil.UnitTestEnv) {
		fqbn := "someboardfqbn"

		instance := &rpc.Instance{Id: int32(1)}
		platformReq := &rpc.PlatformListRequest{
			Instance: instance,
			All:      true,
		}
		boardReq := &rpc.BoardListRequest{
			Instance: instance,
		}
		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().GetPlatforms(platformReq)
		env.ArduinoCli.EXPECT().ConnectedBoards(boardReq).Times(1)

		board, err := env.ArdiCore.Cli.GetTargetBoard(fqbn, "", false)
		assert.NoError(env.T, err)
		assert.Equal(env.T, board.FQBN, fqbn)
	})

	testutil.RunUnitTest("returns target if fqbn provided and onlyConnected true", t, func(env *testutil.UnitTestEnv) {
		boardName := "someboardname"
		fqbn := "someboardfqbn"
		connectedBoard := testutil.GenerateRPCBoard(boardName, fqbn)
		rpcPort := &rpc.Port{
			Address: connectedBoard.Port,
		}

		instance := &rpc.Instance{Id: int32(1)}
		boardReq := &rpc.BoardListRequest{
			Instance: instance,
		}
		platformReq := &rpc.PlatformListRequest{
			Instance: instance,
			All:      true,
		}
		detectedPorts := []*rpc.DetectedPort{
			{
				Port: rpcPort,
				MatchingBoards: []*rpc.BoardListItem{
					{
						Name: connectedBoard.Name,
						Fqbn: connectedBoard.FQBN,
					},
				},
			},
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().GetPlatforms(platformReq)
		env.ArduinoCli.EXPECT().ConnectedBoards(boardReq).Return(detectedPorts, nil)

		board, err := env.ArdiCore.Cli.GetTargetBoard(fqbn, "", true)
		assert.NoError(env.T, err)
		assert.Equal(env.T, board.FQBN, fqbn)
	})

	testutil.RunUnitTest("errors if fqbn provided and onlyConnected true and board not connected", t, func(env *testutil.UnitTestEnv) {
		fqbn := "someboardfqbn"

		instance := &rpc.Instance{Id: int32(1)}
		platformReq := &rpc.PlatformListRequest{
			Instance: instance,
			All:      true,
		}
		boardReq := &rpc.BoardListRequest{
			Instance: instance,
		}
		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().GetPlatforms(platformReq)
		env.ArduinoCli.EXPECT().ConnectedBoards(boardReq)

		board, err := env.ArdiCore.Cli.GetTargetBoard(fqbn, "", true)
		assert.Error(env.T, err)
		assert.Nil(env.T, board)
	})

	testutil.RunUnitTest("returns target if 1 connected board", t, func(env *testutil.UnitTestEnv) {
		boardName := "somboardname"
		boardFQBN := "someboardfqbn"
		connectedBoard := testutil.GenerateRPCBoard(boardName, boardFQBN)
		rpcPort := &rpc.Port{
			Address: connectedBoard.Port,
		}

		instance := &rpc.Instance{Id: int32(1)}
		platformReq := &rpc.PlatformListRequest{
			Instance: instance,
			All:      true,
		}
		boardReq := &rpc.BoardListRequest{
			Instance: instance,
		}
		detectedPorts := []*rpc.DetectedPort{
			{
				Port: rpcPort,
				MatchingBoards: []*rpc.BoardListItem{
					{
						Name: connectedBoard.Name,
						Fqbn: connectedBoard.FQBN,
					},
				},
			},
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().GetPlatforms(platformReq)
		env.ArduinoCli.EXPECT().ConnectedBoards(boardReq).Return(detectedPorts, nil)

		board, err := env.ArdiCore.Cli.GetTargetBoard("", "", false)
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

		rpcPort1 := &rpc.Port{
			Address: connectedBoard1.Port,
		}

		rpcPort2 := &rpc.Port{
			Address: connectedBoard2.Port,
		}

		instance := &rpc.Instance{Id: int32(1)}
		platformReq := &rpc.PlatformListRequest{
			Instance: instance,
			All:      true,
		}
		boardReq := &rpc.BoardListRequest{
			Instance: instance,
		}
		detectedPorts := []*rpc.DetectedPort{
			{
				Port: rpcPort1,
				MatchingBoards: []*rpc.BoardListItem{
					{
						Name: connectedBoard1.Name,
						Fqbn: connectedBoard1.FQBN,
					},
				},
			},
			{
				Port: rpcPort2,
				MatchingBoards: []*rpc.BoardListItem{
					{
						Name: connectedBoard2.Name,
						Fqbn: connectedBoard2.FQBN,
					},
				},
			},
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().GetPlatforms(platformReq)
		env.ArduinoCli.EXPECT().ConnectedBoards(boardReq).Return(detectedPorts, nil)

		env.ClearStdout()
		board, err := env.ArdiCore.Cli.GetTargetBoard("", "", false)
		assert.Error(env.T, err)
		assert.Nil(env.T, board)

		out := env.Stdout.String()
		assert.Contains(env.T, out, board1Name)
		assert.Contains(env.T, out, board1FQBN)
		assert.Contains(env.T, out, board2Name)
		assert.Contains(env.T, out, board2FQBN)
	})

	testutil.RunUnitTest("returns error and prints all available boards if no connected boards found", t, func(env *testutil.UnitTestEnv) {
		platformBoard := testutil.GenerateRPCBoard("board-name", "board-fqbn")

		instance := &rpc.Instance{Id: int32(1)}
		platformReq := &rpc.PlatformListRequest{
			Instance: instance,
			All:      true,
		}
		boardReq := &rpc.BoardListRequest{
			Instance: instance,
		}
		platforms := []*rpc.Platform{
			{
				Id: "test:platform",
				Boards: []*rpc.Board{
					{
						Name: platformBoard.Name,
						Fqbn: platformBoard.FQBN,
					},
				},
			},
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().GetPlatforms(platformReq).Return(platforms, nil)
		env.ArduinoCli.EXPECT().ConnectedBoards(boardReq)

		env.ClearStdout()
		board, err := env.ArdiCore.Cli.GetTargetBoard("", "", false)
		assert.Error(env.T, err)
		assert.Nil(env.T, board)

		out := env.Stdout.String()
		assert.Contains(env.T, out, platformBoard.Name)
		assert.Contains(env.T, out, platformBoard.FQBN)
	})

	testutil.RunUnitTest("returns error and does not print available boards if only connected specified", t, func(env *testutil.UnitTestEnv) {
		platformBoard1 := testutil.GenerateRPCBoard("board-name", "board-fqbn")
		platformBoard2 := testutil.GenerateRPCBoard("board2-name", "board2-fqbn")

		instance := &rpc.Instance{Id: int32(1)}
		platformReq := &rpc.PlatformListRequest{
			Instance: instance,
			All:      true,
		}
		boardReq := &rpc.BoardListRequest{
			Instance: instance,
		}
		platforms := []*rpc.Platform{
			{
				Id: "test:platform",
				Boards: []*rpc.Board{
					{
						Name: platformBoard1.Name,
						Fqbn: platformBoard1.FQBN,
					},
					{
						Name: platformBoard2.Name,
						Fqbn: platformBoard2.FQBN,
					},
				},
			},
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().GetPlatforms(platformReq).Return(platforms, nil)
		env.ArduinoCli.EXPECT().ConnectedBoards(boardReq)

		env.ClearStdout()
		board, err := env.ArdiCore.Cli.GetTargetBoard("", "", true)
		assert.Error(env.T, err)
		assert.Nil(env.T, board)

		out := env.Stdout.String()
		assert.NotContains(env.T, out, platformBoard1.Name)
		assert.NotContains(env.T, out, platformBoard1.FQBN)
		assert.NotContains(env.T, out, platformBoard2.Name)
		assert.NotContains(env.T, out, platformBoard2.FQBN)
	})

	testutil.RunUnitTest("prints fqbns", t, func(env *testutil.UnitTestEnv) {
		boards := testutil.GenerateCmdBoards(10)
		platform := testutil.GenerateCmdPlatform("test-platform", boards)
		platforms := []*rpc.Platform{platform}

		instance := &rpc.Instance{Id: int32(1)}
		req := &rpc.PlatformSearchRequest{
			Instance:    instance,
			AllVersions: false,
		}
		resp := &rpc.PlatformSearchResponse{SearchOutput: platforms}

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().UpdateIndex(gomock.Any(), gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().PlatformSearch(req).Return(resp, nil)

		env.ArdiCore.Board.FQBNS("")

		for _, b := range boards {
			assert.Contains(env.T, env.Stdout.String(), b.GetName())
			assert.Contains(env.T, env.Stdout.String(), b.GetFqbn())
		}
	})

	testutil.RunUnitTest("returns fqbn error", t, func(env *testutil.UnitTestEnv) {
		errString := "dummy error"
		dummyErr := errors.New(errString)
		instance := &rpc.Instance{Id: int32(1)}
		req := &rpc.PlatformSearchRequest{
			Instance:    instance,
			AllVersions: false,
		}
		var resp *rpc.PlatformSearchResponse

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().UpdateIndex(gomock.Any(), gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().PlatformSearch(req).Return(resp, dummyErr)

		err := env.ArdiCore.Board.FQBNS("")
		assert.Error(env.T, err)
		assert.EqualError(env.T, err, errString)
	})

	testutil.RunUnitTest("prints platforms", t, func(env *testutil.UnitTestEnv) {
		boards := testutil.GenerateCmdBoards(10)
		platform := testutil.GenerateCmdPlatform("test-platform", boards)
		platforms := []*rpc.Platform{platform}

		instance := &rpc.Instance{Id: int32(1)}
		req := &rpc.PlatformSearchRequest{
			Instance:    instance,
			AllVersions: false,
		}
		resp := &rpc.PlatformSearchResponse{SearchOutput: platforms}

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().UpdateIndex(gomock.Any(), gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().PlatformSearch(req).Return(resp, nil)

		env.ArdiCore.Board.Platforms("")

		for _, b := range boards {
			assert.Contains(env.T, env.Stdout.String(), b.GetName())
			assert.Contains(env.T, env.Stdout.String(), platform.GetId())
		}
	})

	testutil.RunUnitTest("returns platform error", t, func(env *testutil.UnitTestEnv) {
		errString := "dummy error"
		dummyErr := errors.New(errString)
		instance := &rpc.Instance{Id: int32(1)}
		req := &rpc.PlatformSearchRequest{
			Instance:    instance,
			AllVersions: false,
		}
		var resp *rpc.PlatformSearchResponse

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().UpdateIndex(gomock.Any(), gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().PlatformSearch(req).Return(resp, dummyErr)

		err := env.ArdiCore.Board.Platforms("")
		assert.Error(env.T, err)
		assert.EqualError(env.T, err, errString)
	})
}
