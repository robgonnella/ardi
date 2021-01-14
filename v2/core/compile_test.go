package core_test

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/testutil"
)

func TestCompileCore(t *testing.T) {
	testutil.RunUnitTest("returns nil on success", t, func(env *testutil.UnitTestEnv) {

		projectDir := testutil.BlinkProjectDir()
		expectedFqbn := "some-fqbb"
		expectedSketch := path.Join(projectDir, "blink.ino")
		expectedBuildProps := []string{"build.extra_flags='-DSOME_OPTION'"}
		expectedShowProps := false

		compileOpts := rpc.CompileOpts{
			FQBN:       expectedFqbn,
			SketchDir:  projectDir,
			SketchPath: expectedSketch,
			BuildProps: expectedBuildProps,
			ShowProps:  expectedShowProps,
		}
		env.Client.EXPECT().Compile(compileOpts).Times(1).Return(nil)

		err := env.ArdiCore.Compiler.Compile(compileOpts)
		assert.Nil(env.T, err)
	})

	testutil.RunUnitTest("returns compile error", t, func(env *testutil.UnitTestEnv) {
		errString := "dummy error"
		dummyErr := errors.New(errString)

		projectDir := testutil.BlinkProjectDir()
		expectedFqbn := "some-fqbb"
		expectedSketch := path.Join(projectDir, "blink.ino")
		expectedBuildProps := []string{"build.extra_flags='-DSOME_OPTION'"}
		expectedShowProps := false

		compileOpts := rpc.CompileOpts{
			FQBN:       expectedFqbn,
			SketchDir:  projectDir,
			SketchPath: expectedSketch,
			BuildProps: expectedBuildProps,
			ShowProps:  expectedShowProps,
		}

		env.Client.EXPECT().Compile(compileOpts).Times(1).Return(dummyErr)

		err := env.ArdiCore.Compiler.Compile(compileOpts)
		assert.Error(env.T, err)
		assert.EqualError(env.T, err, errString)
	})

	testutil.RunUnitTest("watches file for changes and recompiles", t, func(env *testutil.UnitTestEnv) {
		sketch, _ := filepath.Abs("test_compilation_file.ino")
		sketchDir := filepath.Dir(sketch)
		file, err := os.OpenFile(sketch, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		assert.NoError(env.T, err)
		defer func() {
			file.Close()
			os.RemoveAll(sketch)
		}()

		errString := "dummy error"
		dummyErr := errors.New(errString)

		projectDir := sketchDir
		expectedFqbn := "some-fqbb"
		expectedSketch := sketch
		expectedBuildProps := []string{"build.extra_flags='-DSOME_OPTION'"}
		expectedShowProps := false

		compileOpts := rpc.CompileOpts{
			FQBN:       expectedFqbn,
			SketchDir:  projectDir,
			SketchPath: expectedSketch,
			BuildProps: expectedBuildProps,
			ShowProps:  expectedShowProps,
		}

		env.Client.EXPECT().Compile(compileOpts).Times(1).Return(nil)

		err = env.ArdiCore.Compiler.Compile(compileOpts)
		assert.NoError(env.T, err)

		env.ClearStdout()
		env.Client.EXPECT().Compile(compileOpts).Times(1).Return(dummyErr)

		go env.ArdiCore.Compiler.WatchForChanges(compileOpts)

		time.Sleep(time.Second)
		_, err = file.WriteString("changes to file\n")
		assert.NoError(env.T, err)

		// wait a second for watcher to trigger
		time.Sleep(time.Second)

		assert.Contains(env.T, env.Stdout.String(), "Recompiling")
		assert.Contains(env.T, env.Stdout.String(), "Compilation failed")
	})
}
