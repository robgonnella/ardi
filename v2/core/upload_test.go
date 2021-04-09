package core_test

import (
	"errors"
	"testing"

	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/golang/mock/gomock"
	"github.com/robgonnella/ardi/v2/mocks"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestUploadCore(t *testing.T) {
	testutil.RunUnitTest("returns nil on success ", t, func(env *testutil.UnitTestEnv) {
		connectedBoard := testutil.GenerateRPCBoard("someboard", "somefqbn")
		projectDir := testutil.BlinkProjectDir()

		instance := &rpc.Instance{Id: int32(1)}
		req := &rpc.UploadRequest{
			Instance:   instance,
			Fqbn:       connectedBoard.FQBN,
			Port:       connectedBoard.Port,
			SketchPath: projectDir,
			Verbose:    true,
		}

		env.Cli.EXPECT().CreateInstance().Return(instance, nil)
		env.Cli.EXPECT().Upload(gomock.Any(), req, gomock.Any(), gomock.Any())

		err := env.ArdiCore.Uploader.Upload(connectedBoard, projectDir)
		assert.Nil(env.T, err)
	})

	testutil.RunUnitTest("returns upload error", t, func(env *testutil.UnitTestEnv) {
		board := testutil.GenerateRPCBoard("whatever", "fqbn")
		dummyErr := errors.New("dummy error")
		projectDir := testutil.BlinkProjectDir()

		instance := &rpc.Instance{Id: int32(1)}
		req := &rpc.UploadRequest{
			Instance:   instance,
			Fqbn:       board.FQBN,
			Port:       board.Port,
			SketchPath: projectDir,
			Verbose:    true,
		}

		env.Cli.EXPECT().CreateInstance().Return(instance, nil)
		env.Cli.EXPECT().Upload(gomock.Any(), req, gomock.Any(), gomock.Any()).Return(nil, dummyErr)

		err := env.ArdiCore.Uploader.Upload(board, projectDir)
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
