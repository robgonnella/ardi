package core_test

import (
	"errors"
	"testing"

	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestUploadCore(t *testing.T) {
	testutil.RunUnitTest("returns nil on success ", t, func(env *testutil.UnitTestEnv) {
		connectedBoard := testutil.GenerateRPCBoard("someboard", "somefqbn")
		targetOpts := core.NewTargetOpts{
			ConnectedBoards: []*rpc.Board{connectedBoard},
			AllBoards:       []*rpc.Board{},
			OnlyConnected:   true,
			FQBN:            "",
			Logger:          env.Logger,
		}
		target, err := core.NewTarget(targetOpts)
		assert.NoError(env.T, err)

		env.Client.EXPECT().Upload(target.Board.FQBN, testutil.BlinkProjectDir(), target.Board.Port).Times(1).Return(nil)
		err = env.ArdiCore.Uploader.Upload(*target, testutil.BlinkProjectDir())
		assert.Nil(env.T, err)
	})

	testutil.RunUnitTest("returns upload error", t, func(env *testutil.UnitTestEnv) {
		dummyErr := errors.New("dummy error")
		connectedBoard := testutil.GenerateRPCBoard("someboard", "somefqbn")
		targetOpts := core.NewTargetOpts{
			ConnectedBoards: []*rpc.Board{connectedBoard},
			AllBoards:       []*rpc.Board{},
			OnlyConnected:   true,
			FQBN:            "",
			Logger:          env.Logger,
		}
		target, err := core.NewTarget(targetOpts)
		assert.NoError(env.T, err)

		env.Client.EXPECT().Upload(target.Board.FQBN, testutil.BlinkProjectDir(), target.Board.Port).Times(1).Return(dummyErr)
		err = env.ArdiCore.Uploader.Upload(*target, testutil.BlinkProjectDir())
		assert.EqualError(env.T, err, dummyErr.Error())
	})
}
