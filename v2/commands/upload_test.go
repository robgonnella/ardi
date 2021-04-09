package commands_test

import (
	"errors"
	"path"
	"testing"

	rpc "github.com/arduino/arduino-cli/rpc/commands"
	"github.com/golang/mock/gomock"
	"github.com/robgonnella/ardi/v2/mocks"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestUploadCommand(t *testing.T) {
	testutil.RunIntegrationTest("Uploading", t, func(env *testutil.IntegrationTestEnv) {
		ctrl := gomock.NewController(env.T)
		inst := mocks.NewMockCli(ctrl)

		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		board := testutil.GenerateRPCBoard("Arduino Mega", "arduino:avr:mega")
		buildName := "blink"
		sketchDir := testutil.BlinkProjectDir()
		sketch := path.Join(sketchDir, "blink.ino")
		bogusSketch := "noop"
		fqbn := testutil.ArduinoMegaFQBN()

		args := []string{"add", "build", "--name", buildName, "--fqbn", fqbn, "--sketch", sketchDir}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		instance := &rpc.Instance{Id: int32(1)}
		req := &rpc.UploadReq{
			Instance:   instance,
			Fqbn:       fqbn,
			SketchPath: sketchDir,
			Port:       board.Port,
		}

		platformReq := &rpc.PlatformListReq{
			Instance:      instance,
			UpdatableOnly: false,
			All:           true,
		}

		boardItem := &rpc.BoardListItem{
			Name: board.Name,
			FQBN: board.FQBN,
		}
		port := &rpc.DetectedPort{
			Address: board.Port,
			Boards:  []*rpc.BoardListItem{boardItem},
		}
		detectedPorts := []*rpc.DetectedPort{port}

		env.T.Run("uploads a build", func(st *testing.T) {
			inst.EXPECT().CreateInstance().Return(instance, nil).AnyTimes()
			inst.EXPECT().ConnectedBoards(instance.GetId()).Return(detectedPorts, nil)
			inst.EXPECT().GetPlatforms(platformReq)
			inst.EXPECT().Upload(gomock.Any(), req, gomock.Any(), gomock.Any())

			args = []string{"upload", buildName}
			err = env.ExecuteWithMockCli(args, inst)
			assert.NoError(st, err)
		})

		env.T.Run("uploads a sketch", func(st *testing.T) {
			inst.EXPECT().CreateInstance().Return(instance, nil).AnyTimes()
			inst.EXPECT().ConnectedBoards(instance.GetId()).Return(detectedPorts, nil)
			inst.EXPECT().GetPlatforms(platformReq)
			inst.EXPECT().Upload(gomock.Any(), req, gomock.Any(), gomock.Any())

			args = []string{"upload", "--fqbn", fqbn, sketch}
			err = env.ExecuteWithMockCli(args, inst)
			assert.NoError(st, err)
		})

		env.T.Run("returns upload errors", func(st *testing.T) {
			dummyErr := errors.New("dummy")
			inst.EXPECT().CreateInstance().Return(instance, nil).AnyTimes()
			inst.EXPECT().ConnectedBoards(instance.GetId()).Return(detectedPorts, nil)
			inst.EXPECT().GetPlatforms(platformReq)
			inst.EXPECT().Upload(gomock.Any(), req, gomock.Any(), gomock.Any()).Return(nil, dummyErr)

			args = []string{"upload", "--fqbn", fqbn, sketch}
			err = env.ExecuteWithMockCli(args, inst)
			assert.Error(st, err)
			assert.EqualError(st, err, dummyErr.Error())
		})

		env.T.Run("errors if sketch not found", func(st *testing.T) {
			args = []string{"upload", "--fqbn", fqbn, bogusSketch}
			err = env.ExecuteWithMockCli(args, inst)
			assert.Error(st, err)
		})

		env.T.Run("errors if no board connected", func(st *testing.T) {
			inst.EXPECT().CreateInstance().Return(instance, nil).AnyTimes()
			inst.EXPECT().ConnectedBoards(instance.GetId()).Return([]*rpc.DetectedPort{}, nil)
			inst.EXPECT().GetPlatforms(platformReq)

			args = []string{"upload", buildName, "--attach"}
			err = env.ExecuteWithMockCli(args, inst)
			assert.Error(st, err)
		})
	})
}
