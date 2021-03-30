package commands_test

import (
	"testing"

	rpc "github.com/arduino/arduino-cli/rpc/commands"
	"github.com/golang/mock/gomock"
	"github.com/robgonnella/ardi/v2/mocks"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestUploadCommand(t *testing.T) {
	testutil.RunIntegrationTest("uploads a build", t, func(env *testutil.IntegrationTestEnv) {
		ctrl := gomock.NewController(env.T)
		inst := mocks.NewMockCli(ctrl)

		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		board := testutil.GenerateRPCBoard("Arduino Mega", "arduino:avr:mega")
		buildName := "blink"
		sketchDir := testutil.BlinkProjectDir()
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

		inst.EXPECT().CreateInstance().Return(instance, nil).AnyTimes()
		inst.EXPECT().ConnectedBoards(instance.GetId()).Return(detectedPorts, nil)
		inst.EXPECT().GetPlatforms(platformReq)
		inst.EXPECT().Upload(gomock.Any(), req, gomock.Any(), gomock.Any())

		args = []string{"upload", buildName}
		err = env.ExecuteWithMockCli(args, inst)
		assert.NoError(env.T, err)
	})
}
