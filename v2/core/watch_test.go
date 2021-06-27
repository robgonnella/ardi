package core_test

import (
	"errors"
	"os/exec"
	"path"
	"testing"
	"time"

	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/golang/mock/gomock"
	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/mocks"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestWatchCore(t *testing.T) {
	sketchDir := testutil.BlinkProjectDir()
	sketch := path.Join(sketchDir, "blink.ino")
	sketchCopy := path.Join(testutil.BlinkCopyProjectDir(), "blink2.ino")
	fqbn := testutil.ArduinoMegaFQBN()
	buildProps := []string{}
	board := testutil.GenerateRPCBoard("arduino:avr:mega", fqbn)
	compileOpts := cli.CompileOpts{
		FQBN:       fqbn,
		SketchDir:  sketchDir,
		SketchPath: sketch,
		BuildProps: buildProps,
		ShowProps:  false,
	}

	testutil.RunUnitTest("recompiles and reuploads on file change", t, func(env *testutil.UnitTestEnv) {
		cpCmd := exec.Command("cp", sketchCopy, sketch)
		env.ClearStdout()
		port := mocks.NewMockSerialPort(env.Ctrl)
		port.EXPECT().Stop().AnyTimes()
		port.EXPECT().Watch().AnyTimes()

		instance := &rpc.Instance{Id: int32(1)}
		compileReq := &rpc.CompileRequest{
			Instance:        instance,
			Fqbn:            fqbn,
			SketchPath:      sketch,
			ExportDir:       path.Join(sketchDir, "build"),
			BuildProperties: compileOpts.BuildProps,
			ShowProperties:  compileOpts.ShowProps,
			Verbose:         true,
		}

		uploadReq := &rpc.UploadRequest{
			Instance:   instance,
			Fqbn:       fqbn,
			SketchPath: sketchDir,
			Port:       board.Port,
			Verbose:    true,
		}

		env.Cli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.Cli.EXPECT().Compile(gomock.Any(), compileReq, gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		env.Cli.EXPECT().Upload(gomock.Any(), uploadReq, gomock.Any(), gomock.Any()).AnyTimes()

		targets := core.WatchCoreTargets{
			Board:       board,
			CompileOpts: &compileOpts,
			Baud:        9600,
			Port:        port,
		}
		env.ArdiCore.Watcher.SetTargets(targets)
		go env.ArdiCore.Watcher.Watch()

		time.Sleep(time.Second)
		err := cpCmd.Run()
		assert.NoError(env.T, err)

		// wait a second for watcher to trigger
		time.Sleep(time.Second)

		assert.Contains(env.T, env.Stdout.String(), "Reuploading")
		assert.Contains(env.T, env.Stdout.String(), "Upload successful")
		env.ArdiCore.Watcher.Stop()
	})

	testutil.RunUnitTest("does not reupload on compilation error", t, func(env *testutil.UnitTestEnv) {
		cpCmd := exec.Command("cp", sketchCopy, sketch)
		env.ClearStdout()
		port := mocks.NewMockSerialPort(env.Ctrl)
		port.EXPECT().Stop().AnyTimes()
		port.EXPECT().Watch().Times(1)

		dummyErr := errors.New("dummy errror")
		instance := &rpc.Instance{Id: int32(1)}
		compileReq := &rpc.CompileRequest{
			Instance:        instance,
			Fqbn:            fqbn,
			SketchPath:      sketch,
			ExportDir:       path.Join(sketchDir, "build"),
			BuildProperties: compileOpts.BuildProps,
			ShowProperties:  compileOpts.ShowProps,
			Verbose:         true,
		}

		uploadReq := &rpc.UploadRequest{
			Instance:   instance,
			Fqbn:       fqbn,
			SketchPath: sketchDir,
			Port:       board.Port,
			Verbose:    true,
		}

		env.Cli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.Cli.EXPECT().Compile(gomock.Any(), compileReq, gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, dummyErr)
		env.Cli.EXPECT().Upload(gomock.Any(), uploadReq, gomock.Any(), gomock.Any()).AnyTimes()

		targets := core.WatchCoreTargets{
			Board:       board,
			CompileOpts: &compileOpts,
			Baud:        9600,
			Port:        port,
		}
		env.ArdiCore.Watcher.SetTargets(targets)
		go env.ArdiCore.Watcher.Watch()

		time.Sleep(time.Second)
		err := cpCmd.Run()
		assert.NoError(env.T, err)

		// wait a second for watcher to trigger
		time.Sleep(time.Second)

		assert.NotContains(env.T, env.Stdout.String(), "Reuploading")
		assert.NotContains(env.T, env.Stdout.String(), "Upload successful")
		env.ArdiCore.Watcher.Stop()
	})
}
