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
	rpcPort := &rpc.Port{
		Address: board.Port,
	}
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
			Port:       rpcPort,
			Verbose:    true,
		}

		targets := core.WatchCoreTargets{
			Board:       board,
			CompileOpts: &compileOpts,
			Baud:        9600,
		}

		env.SerialPort.EXPECT().SetTargets(board.Port, targets.Baud).AnyTimes()
		env.SerialPort.EXPECT().SetTargets("", 0).AnyTimes()
		env.SerialPort.EXPECT().Close().AnyTimes()
		env.SerialPort.EXPECT().Watch().AnyTimes()

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), compileReq, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		env.ArduinoCli.EXPECT().Upload(gomock.Any(), uploadReq, gomock.Any(), gomock.Any()).AnyTimes()

		env.ArdiCore.Watcher.SetTargets(targets)
		go env.ArdiCore.Watcher.Watch()

		time.Sleep(time.Second)
		env.ClearStdout()
		err := cpCmd.Run()
		assert.NoError(env.T, err)

		// wait a second for watcher to trigger
		time.Sleep(time.Second)

		assert.Contains(env.T, env.Stdout.String(), "Uploading...")
		assert.Contains(env.T, env.Stdout.String(), "Upload successful")
	})

	testutil.RunUnitTest("does not reupload on compilation error", t, func(env *testutil.UnitTestEnv) {
		cpCmd := exec.Command("cp", sketchCopy, sketch)
		env.ClearStdout()

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
			Port:       rpcPort,
			Verbose:    true,
		}

		targets := core.WatchCoreTargets{
			Board:       board,
			CompileOpts: &compileOpts,
			Baud:        9600,
		}

		env.SerialPort.EXPECT().SetTargets(board.Port, targets.Baud).AnyTimes()
		env.SerialPort.EXPECT().SetTargets("", 0).AnyTimes()
		env.SerialPort.EXPECT().Close().AnyTimes()
		env.SerialPort.EXPECT().Watch().AnyTimes()

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), compileReq, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, dummyErr)
		env.ArduinoCli.EXPECT().Upload(gomock.Any(), uploadReq, gomock.Any(), gomock.Any()).AnyTimes()

		env.ArdiCore.Watcher.SetTargets(targets)
		go env.ArdiCore.Watcher.Watch()

		time.Sleep(time.Second)
		env.ClearStdout()
		err := cpCmd.Run()
		assert.NoError(env.T, err)

		// wait a second for watcher to trigger
		time.Sleep(time.Second)

		assert.NotContains(env.T, env.Stdout.String(), "Uploading...")
		assert.NotContains(env.T, env.Stdout.String(), "Upload successful")
	})
}
