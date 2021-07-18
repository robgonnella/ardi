package commands_test

import (
	"errors"
	"os"
	"path"
	"testing"

	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/golang/mock/gomock"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestUploadCommand(t *testing.T) {
	board := testutil.GenerateRPCBoard("Arduino Mega", "arduino:avr:mega")
	buildName := "blink"
	sketchDir := testutil.BlinkProjectDir()
	sketch := path.Join(sketchDir, "blink.ino")
	bogusSketch := "noop"
	fqbn := testutil.ArduinoMegaFQBN()

	instance := &rpc.Instance{Id: int32(1)}

	platformReq := &rpc.PlatformListRequest{
		Instance:      instance,
		UpdatableOnly: false,
		All:           true,
	}

	boardItem := &rpc.BoardListItem{
		Name: board.Name,
		Fqbn: board.FQBN,
	}

	port := &rpc.DetectedPort{
		Address: board.Port,
		Boards:  []*rpc.BoardListItem{boardItem},
	}

	detectedPorts := []*rpc.DetectedPort{port}

	addBuild := func(e *testutil.MockCliIntegrationTestEnv) {
		err := e.RunProjectInit()
		assert.NoError(e.T, err)

		args := []string{"add", "build", "--name", buildName, "--fqbn", fqbn, "--sketch", sketchDir}
		err = e.Execute(args)
		assert.NoError(e.T, err)
	}

	expectUsuals := func(e *testutil.MockCliIntegrationTestEnv) {
		e.Cli.EXPECT().CreateInstance().Return(instance).MaxTimes(1)
		e.Cli.EXPECT().GetPlatforms(platformReq).MaxTimes(1)
	}

	testutil.RunMockCliIntegrationTest("uploads a build", t, func(env *testutil.MockCliIntegrationTestEnv) {
		addBuild(env)
		expectUsuals(env)
		req := &rpc.UploadRequest{
			Instance:   instance,
			SketchPath: sketchDir,
			Fqbn:       fqbn,
			Port:       board.Port,
		}
		env.Cli.EXPECT().ConnectedBoards(instance.GetId()).Return(detectedPorts, nil).MaxTimes(1)
		env.Cli.EXPECT().Upload(gomock.Any(), req, gomock.Any(), gomock.Any()).MaxTimes(1)
		args := []string{"upload", buildName}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunMockCliIntegrationTest("uploads a sketch", t, func(env *testutil.MockCliIntegrationTestEnv) {
		addBuild(env)
		expectUsuals(env)
		req := &rpc.UploadRequest{
			Instance:   instance,
			SketchPath: sketchDir,
			Fqbn:       fqbn,
			Port:       board.Port,
		}
		env.Cli.EXPECT().ConnectedBoards(instance.GetId()).Return(detectedPorts, nil).MaxTimes(1)
		env.Cli.EXPECT().Upload(gomock.Any(), req, gomock.Any(), gomock.Any()).MaxTimes(1)
		args := []string{"upload", "--fqbn", fqbn, sketch}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunMockCliIntegrationTest("uploads a sketch using default build auto-detected values", t, func(env *testutil.MockCliIntegrationTestEnv) {
		addBuild(env)
		expectUsuals(env)

		env.RunProjectInit()

		defaultBuild := "default"
		defaultSketchDir := testutil.PixieProjectDir()
		defaultSketch := path.Join(defaultSketchDir, "pixie.ino")
		defaultFQBN := testutil.ArduinoMegaFQBN()

		args := []string{"add", "build", "--name", defaultBuild, "--fqbn", defaultFQBN, "--sketch", defaultSketch}

		err := env.Execute(args)
		assert.NoError(env.T, err)

		defaultUploadReq := &rpc.UploadRequest{
			Instance:   instance,
			Fqbn:       defaultFQBN,
			SketchPath: defaultSketchDir,
			Port:       board.Port,
		}

		env.Cli.EXPECT().ConnectedBoards(instance.GetId()).Return(detectedPorts, nil).MaxTimes(1)
		env.Cli.EXPECT().Upload(gomock.Any(), defaultUploadReq, gomock.Any(), gomock.Any()).MaxTimes(1)

		args = []string{"upload"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunMockCliIntegrationTest("uploads a sketch using single build values", t, func(env *testutil.MockCliIntegrationTestEnv) {
		addBuild(env)
		expectUsuals(env)
		req := &rpc.UploadRequest{
			Instance:   instance,
			SketchPath: sketchDir,
			Fqbn:       fqbn,
			Port:       board.Port,
		}
		env.Cli.EXPECT().ConnectedBoards(instance.GetId()).Return(detectedPorts, nil).MaxTimes(1)
		env.Cli.EXPECT().Upload(gomock.Any(), req, gomock.Any(), gomock.Any()).MaxTimes(1)

		args := []string{"upload"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunMockCliIntegrationTest("uploads a sketch using auto detected .ino values", t, func(env *testutil.MockCliIntegrationTestEnv) {
		currentDir, _ := os.Getwd()
		pixieDir := testutil.PixieProjectDir()

		os.Chdir(pixieDir)

		defer func() {
			testutil.CleanPixieDir()
			os.Chdir(currentDir)
		}()

		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		expectedReq := &rpc.UploadRequest{
			Instance:   instance,
			SketchPath: pixieDir,
			Fqbn:       board.FQBN,
			Port:       board.Port,
		}

		expectUsuals(env)
		env.Cli.EXPECT().ConnectedBoards(instance.GetId()).Return(detectedPorts, nil).MaxTimes(1)
		env.Cli.EXPECT().Upload(gomock.Any(), expectedReq, gomock.Any(), gomock.Any()).MaxTimes(1)

		args := []string{"upload"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunMockCliIntegrationTest("returns upload errors", t, func(env *testutil.MockCliIntegrationTestEnv) {
		addBuild(env)
		expectUsuals(env)
		dummyErr := errors.New("dummy")
		req := &rpc.UploadRequest{
			Instance:   instance,
			SketchPath: sketchDir,
			Fqbn:       fqbn,
			Port:       board.Port,
		}
		env.Cli.EXPECT().ConnectedBoards(instance.GetId()).Return(detectedPorts, nil).MaxTimes(1)
		env.Cli.EXPECT().Upload(gomock.Any(), req, gomock.Any(), gomock.Any()).Return(nil, dummyErr).MaxTimes(1)
		args := []string{"upload", "--fqbn", fqbn, sketch}
		err := env.Execute(args)
		assert.Error(env.T, err)
		assert.EqualError(env.T, err, dummyErr.Error())
	})

	testutil.RunMockCliIntegrationTest("errors if sketch not found", t, func(env *testutil.MockCliIntegrationTestEnv) {
		addBuild(env)
		expectUsuals(env)
		env.Cli.EXPECT().ConnectedBoards(instance.GetId()).Return([]*rpc.DetectedPort{}, nil).MaxTimes(1)
		args := []string{"upload", "--fqbn", fqbn, bogusSketch}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockCliIntegrationTest("errors if no board connected", t, func(env *testutil.MockCliIntegrationTestEnv) {
		addBuild(env)
		expectUsuals(env)
		env.Cli.EXPECT().ConnectedBoards(instance.GetId()).Return([]*rpc.DetectedPort{}, nil).MaxTimes(1)
		args := []string{"upload", buildName, "--attach"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})
}
