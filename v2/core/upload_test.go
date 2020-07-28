package core_test

import (
	"errors"
	"testing"

	"github.com/robgonnella/ardi/v2/mocks"
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestUploadCore(t *testing.T) {
	testutil.RunUnitTest("returns nil on success ", t, func(env *testutil.UnitTestEnv) {
		connectedBoard := testutil.GenerateRPCBoard("someboard", "somefqbn")
		connectedBoards := []*rpc.Board{connectedBoard}

		env.Client.EXPECT().ConnectedBoards().Times(1).Return(connectedBoards)
		env.Client.EXPECT().AllBoards().Times(1).Return([]*rpc.Board{})

		env.ClearStdout()
		board, err := env.ArdiCore.GetTargetBoard("", true)
		assert.NoError(env.T, err)

		env.Client.EXPECT().Upload(board.FQBN, testutil.BlinkProjectDir(), board.Port).Times(1).Return(nil)

		err = env.ArdiCore.Uploader.Upload(board, testutil.BlinkProjectDir())
		assert.Nil(env.T, err)
	})

	testutil.RunUnitTest("returns upload error", t, func(env *testutil.UnitTestEnv) {
		dummyErr := errors.New("dummy error")
		connectedBoard := testutil.GenerateRPCBoard("someboard", "somefqbn")
		connectedBoards := []*rpc.Board{connectedBoard}

		env.Client.EXPECT().ConnectedBoards().Times(1).Return(connectedBoards)
		env.Client.EXPECT().AllBoards().Times(1).Return([]*rpc.Board{})

		env.ClearStdout()
		board, err := env.ArdiCore.GetTargetBoard("", true)
		assert.NoError(env.T, err)

		env.Client.EXPECT().Upload(board.FQBN, testutil.BlinkProjectDir(), board.Port).Times(1).Return(dummyErr)

		err = env.ArdiCore.Uploader.Upload(board, testutil.BlinkProjectDir())
		assert.EqualError(env.T, err, dummyErr.Error())
	})

	testutil.RunUnitTest("attaches to board port to print logs", t, func(env *testutil.UnitTestEnv) {
		device := "/dev/null"
		baud := 9600
		port := mocks.NewMockSerialPort(env.Ctrl)

		port.EXPECT().Stop().Times(1)
		port.EXPECT().Watch().Times(1)
		env.ArdiCore.Uploader.Attach(device, baud, port)
	})
}