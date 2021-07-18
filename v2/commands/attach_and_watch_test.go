package commands_test

import (
	"os"
	"os/exec"
	"path"
	"testing"
	"time"

	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/golang/mock/gomock"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAttachAndWatchCommand(t *testing.T) {
	board := testutil.GenerateRPCBoard("Arduino Mega", "arduino:avr:mega")
	buildName := "blink"
	sketchDir := testutil.BlinkProjectDir()
	sketch := path.Join(sketchDir, "blink.ino")
	sketchCopy := path.Join(testutil.BlinkCopyProjectDir(), "blink2.ino")
	fqbn := testutil.ArduinoMegaFQBN()

	instance := &rpc.Instance{Id: int32(1)}

	uploadReq := &rpc.UploadRequest{
		Instance:   instance,
		Fqbn:       fqbn,
		SketchPath: sketchDir,
		Port:       board.Port,
	}

	platformReq := &rpc.PlatformListRequest{
		Instance:      instance,
		UpdatableOnly: false,
		All:           true,
	}

	compileReq := &rpc.CompileRequest{
		Instance:        instance,
		Fqbn:            board.FQBN,
		SketchPath:      sketch,
		ExportDir:       path.Join(sketchDir, "build"),
		BuildProperties: []string{},
		ShowProperties:  false,
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

	testutil.RunMockIntegrationTest("attaches and watches saved build", t, func(env *testutil.MockIntegrationTestEnv) {
		cpCmd := exec.Command("cp", sketchCopy, sketch)

		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		args := []string{"add", "build", "--name", buildName, "--fqbn", fqbn, "--sketch", sketchDir}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		env.SerialPort.EXPECT().SetTargets(board.Port, 9600).AnyTimes()
		env.SerialPort.EXPECT().SetTargets("", 0).AnyTimes()
		env.SerialPort.EXPECT().Watch().AnyTimes()
		env.SerialPort.EXPECT().Close().AnyTimes()

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), compileReq, gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		env.ArduinoCli.EXPECT().ConnectedBoards(instance.GetId()).Return(detectedPorts, nil)
		env.ArduinoCli.EXPECT().GetPlatforms(platformReq)
		env.ArduinoCli.EXPECT().Upload(gomock.Any(), uploadReq, gomock.Any(), gomock.Any()).AnyTimes()

		args = []string{"attach-and-watch", buildName}
		go env.Execute(args)

		time.Sleep(time.Second * 5)

		env.ClearStdout()
		err = cpCmd.Run()
		assert.NoError(env.T, err)

		// wait a second for watcher to trigger
		time.Sleep(time.Second)

		assert.Contains(env.T, env.Stdout.String(), "Uploading...")
		assert.Contains(env.T, env.Stdout.String(), "Upload successful")
	})

	testutil.RunMockIntegrationTest("attaches and watches directory sketch", t, func(env *testutil.MockIntegrationTestEnv) {
		cpCmd := exec.Command("cp", sketchCopy, sketch)

		cwd, _ := os.Getwd()
		os.Chdir(testutil.BlinkProjectDir())
		defer os.Chdir(cwd)

		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		env.SerialPort.EXPECT().SetTargets(board.Port, 9600).AnyTimes()
		env.SerialPort.EXPECT().SetTargets("", 0).AnyTimes()
		env.SerialPort.EXPECT().Watch().AnyTimes()
		env.SerialPort.EXPECT().Close().AnyTimes()

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), compileReq, gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		env.ArduinoCli.EXPECT().ConnectedBoards(instance.GetId()).Return(detectedPorts, nil)
		env.ArduinoCli.EXPECT().GetPlatforms(platformReq)
		env.ArduinoCli.EXPECT().Upload(gomock.Any(), uploadReq, gomock.Any(), gomock.Any()).AnyTimes()

		args := []string{"attach-and-watch", "--fqbn", fqbn, sketch}
		go env.Execute(args)

		time.Sleep(time.Second * 5)

		env.ClearStdout()
		err = cpCmd.Run()
		assert.NoError(env.T, err)

		// wait a second for watcher to trigger
		time.Sleep(time.Second)

		assert.Contains(env.T, env.Stdout.String(), "Uploading...")
		assert.Contains(env.T, env.Stdout.String(), "Upload successful")
	})

	testutil.RunMockIntegrationTest("attaches and watches using auto detected values", t, func(env *testutil.MockIntegrationTestEnv) {
		cpCmd := exec.Command("cp", sketchCopy, sketch)

		cwd, _ := os.Getwd()
		os.Chdir(testutil.BlinkProjectDir())

		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		env.SerialPort.EXPECT().SetTargets(board.Port, 9600).AnyTimes()
		env.SerialPort.EXPECT().SetTargets("", 0).AnyTimes()
		env.SerialPort.EXPECT().Watch().AnyTimes()
		env.SerialPort.EXPECT().Close().AnyTimes()

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), compileReq, gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		env.ArduinoCli.EXPECT().ConnectedBoards(instance.GetId()).Return(detectedPorts, nil)
		env.ArduinoCli.EXPECT().GetPlatforms(platformReq)
		env.ArduinoCli.EXPECT().Upload(gomock.Any(), uploadReq, gomock.Any(), gomock.Any()).AnyTimes()

		args := []string{"attach-and-watch"}
		go env.Execute(args)

		time.Sleep(time.Second * 5)

		env.ClearStdout()
		err = cpCmd.Run()
		assert.NoError(env.T, err)

		// wait a second for watcher to trigger
		time.Sleep(time.Second)

		assert.Contains(env.T, env.Stdout.String(), "Uploading...")
		assert.Contains(env.T, env.Stdout.String(), "Upload successful")
		os.Chdir(cwd)
	})

	testutil.RunMockIntegrationTest("returns error if no sketch found in current directory", t, func(env *testutil.MockIntegrationTestEnv) {
		env.RunProjectInit()
		args := []string{"attach-and-watch"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})
}
