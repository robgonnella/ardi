package core_test

import (
	"errors"
	"os/exec"
	"path"
	"testing"
	"time"

	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/mocks"
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestWatchCore(t *testing.T) {
	sketchDir := testutil.BlinkProjectDir()
	sketch := path.Join(sketchDir, "blink.ino")
	sketchCopy := path.Join(testutil.BlinkCopyProjectDir(), "blink2.ino")
	fqbn := testutil.ArduinoMegaFQBN()
	buildProps := []string{}
	connectedBoard := testutil.GenerateRPCBoard("arduino:avr:mega", fqbn)
	compileOpts := rpc.CompileOpts{
		FQBN:       fqbn,
		SketchDir:  sketchDir,
		SketchPath: sketch,
		BuildProps: buildProps,
		ShowProps:  false,
		ExportName: "",
	}

	testutil.RunUnitTest("recompiles and reuploads on file change", t, func(env *testutil.UnitTestEnv) {
		cpCmd := exec.Command("cp", sketchCopy, sketch)
		targetOpts := core.NewTargetOpts{
			ConnectedBoards: []*rpc.Board{connectedBoard},
			AllBoards:       []*rpc.Board{},
			OnlyConnected:   true,
			FQBN:            "",
			Logger:          env.Logger,
		}
		target, err := core.NewTarget(targetOpts)
		assert.NoError(env.T, err)

		env.ClearStdout()
		port := mocks.NewMockSerialPort(env.Ctrl)
		port.EXPECT().Stop().AnyTimes()
		port.EXPECT().Watch().AnyTimes()
		env.Client.EXPECT().Compile(compileOpts).AnyTimes().Return(nil)
		env.Client.EXPECT().Upload(fqbn, sketchDir, target.Board.Port).AnyTimes().Return(nil)

		go env.ArdiCore.Watcher.Watch(compileOpts, *target, 9600, port)

		time.Sleep(time.Second)
		err = cpCmd.Run()
		assert.NoError(env.T, err)

		// wait a second for watcher to trigger
		time.Sleep(time.Second)

		assert.Contains(env.T, env.Stdout.String(), "Reuploading")
		assert.Contains(env.T, env.Stdout.String(), "Upload successful")
		env.ArdiCore.Watcher.Stop()
	})

	testutil.RunUnitTest("does not reupload on compilation error", t, func(env *testutil.UnitTestEnv) {
		cpCmd := exec.Command("cp", sketchCopy, sketch)
		targetOpts := core.NewTargetOpts{
			ConnectedBoards: []*rpc.Board{connectedBoard},
			AllBoards:       []*rpc.Board{},
			OnlyConnected:   true,
			FQBN:            "",
			Logger:          env.Logger,
		}
		target, err := core.NewTarget(targetOpts)
		assert.NoError(env.T, err)

		env.ClearStdout()
		port := mocks.NewMockSerialPort(env.Ctrl)
		port.EXPECT().Stop().AnyTimes()
		port.EXPECT().Watch().Times(1)
		env.Client.EXPECT().Compile(compileOpts).AnyTimes().Return(errors.New("dummy errror"))

		go env.ArdiCore.Watcher.Watch(compileOpts, *target, 9600, port)

		time.Sleep(time.Second)
		err = cpCmd.Run()
		assert.NoError(env.T, err)

		// wait a second for watcher to trigger
		time.Sleep(time.Second)

		assert.NotContains(env.T, env.Stdout.String(), "Reuploading")
		assert.NotContains(env.T, env.Stdout.String(), "Upload successful")
		env.ArdiCore.Watcher.Stop()
	})
}
