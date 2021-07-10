package commands_test

import (
	"os"
	"os/exec"
	"path"
	"testing"
	"time"

	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/golang/mock/gomock"
	"github.com/robgonnella/ardi/v2/mocks"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAttachAndWatchCommand(t *testing.T) {
	testutil.RunIntegrationTest("attaches and watches saved build", t, func(env *testutil.IntegrationTestEnv) {
		ctrl := gomock.NewController(env.T)
		inst := mocks.NewMockCli(ctrl)

		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		board := testutil.GenerateRPCBoard("Arduino Mega", "arduino:avr:mega")
		buildName := "blink"
		sketchDir := testutil.BlinkProjectDir()
		sketch := path.Join(sketchDir, "blink.ino")
		sketchCopy := path.Join(testutil.BlinkCopyProjectDir(), "blink2.ino")
		fqbn := testutil.ArduinoMegaFQBN()
		cpCmd := exec.Command("cp", sketchCopy, sketch)

		args := []string{"add", "build", "--name", buildName, "--fqbn", fqbn, "--sketch", sketchDir}
		err = env.Execute(args)
		assert.NoError(env.T, err)

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

		inst.EXPECT().CreateInstance().Return(instance).AnyTimes()
		inst.EXPECT().Compile(gomock.Any(), compileReq, gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		inst.EXPECT().ConnectedBoards(instance.GetId()).Return(detectedPorts, nil)
		inst.EXPECT().GetPlatforms(platformReq)
		inst.EXPECT().Upload(gomock.Any(), uploadReq, gomock.Any(), gomock.Any()).AnyTimes()

		args = []string{"attach-and-watch", buildName}
		go env.ExecuteWithMockCli(args, inst)

		time.Sleep(time.Second * 5)

		env.ClearStdout()
		err = cpCmd.Run()
		assert.NoError(env.T, err)

		// wait a second for watcher to trigger
		time.Sleep(time.Second)

		assert.Contains(env.T, env.Stdout.String(), "Reuploading")
		assert.Contains(env.T, env.Stdout.String(), "Upload successful")
	})

	testutil.RunIntegrationTest("attaches and watches directory sketch", t, func(env *testutil.IntegrationTestEnv) {
		ctrl := gomock.NewController(env.T)
		inst := mocks.NewMockCli(ctrl)

		board := testutil.GenerateRPCBoard("Arduino Mega", "arduino:avr:mega")

		sketchDir := testutil.BlinkProjectDir()
		sketch := path.Join(sketchDir, "blink.ino")
		sketchCopy := path.Join(testutil.BlinkCopyProjectDir(), "blink2.ino")
		fqbn := testutil.ArduinoMegaFQBN()
		cpCmd := exec.Command("cp", sketchCopy, sketch)

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

		cwd, _ := os.Getwd()
		os.Chdir(testutil.BlinkProjectDir())
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		inst.EXPECT().CreateInstance().Return(instance).AnyTimes()
		inst.EXPECT().Compile(gomock.Any(), compileReq, gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		inst.EXPECT().ConnectedBoards(instance.GetId()).Return(detectedPorts, nil)
		inst.EXPECT().GetPlatforms(platformReq)
		inst.EXPECT().Upload(gomock.Any(), uploadReq, gomock.Any(), gomock.Any()).AnyTimes()

		args := []string{"attach-and-watch", "--fqbn", fqbn, sketch}
		go env.ExecuteWithMockCli(args, inst)

		time.Sleep(time.Second * 5)

		env.ClearStdout()
		err = cpCmd.Run()
		assert.NoError(env.T, err)

		// wait a second for watcher to trigger
		time.Sleep(time.Second)

		assert.Contains(env.T, env.Stdout.String(), "Reuploading")
		assert.Contains(env.T, env.Stdout.String(), "Upload successful")
		os.Chdir(cwd)
	})

	testutil.RunIntegrationTest("attaches and watches using auto detected values", t, func(env *testutil.IntegrationTestEnv) {
		ctrl := gomock.NewController(env.T)
		inst := mocks.NewMockCli(ctrl)

		board := testutil.GenerateRPCBoard("Arduino Mega", "arduino:avr:mega")

		sketchDir := testutil.BlinkProjectDir()
		sketch := path.Join(sketchDir, "blink.ino")
		sketchCopy := path.Join(testutil.BlinkCopyProjectDir(), "blink2.ino")
		fqbn := testutil.ArduinoMegaFQBN()
		cpCmd := exec.Command("cp", sketchCopy, sketch)

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

		cwd, _ := os.Getwd()
		os.Chdir(testutil.BlinkProjectDir())
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		inst.EXPECT().CreateInstance().Return(instance).AnyTimes()
		inst.EXPECT().Compile(gomock.Any(), compileReq, gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		inst.EXPECT().ConnectedBoards(instance.GetId()).Return(detectedPorts, nil)
		inst.EXPECT().GetPlatforms(platformReq)
		inst.EXPECT().Upload(gomock.Any(), uploadReq, gomock.Any(), gomock.Any()).AnyTimes()

		args := []string{"attach-and-watch"}
		go env.ExecuteWithMockCli(args, inst)

		time.Sleep(time.Second * 5)

		env.ClearStdout()
		err = cpCmd.Run()
		assert.NoError(env.T, err)

		// wait a second for watcher to trigger
		time.Sleep(time.Second)

		assert.Contains(env.T, env.Stdout.String(), "Reuploading")
		assert.Contains(env.T, env.Stdout.String(), "Upload successful")
		os.Chdir(cwd)
	})
}
