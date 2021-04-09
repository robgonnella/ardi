package core_test

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
	"github.com/robgonnella/ardi/v2/testutil"
)

func TestCompileCore(t *testing.T) {
	testutil.RunUnitTest("returns nil on success", t, func(env *testutil.UnitTestEnv) {

		projectDir := testutil.BlinkProjectDir()
		expectedFqbn := "some-fqbb"
		expectedSketch := path.Join(projectDir, "blink.ino")
		expectedBuildProps := []string{"build.extra_flags='-DSOME_OPTION'"}
		expectedShowProps := false
		expectedExportDir := path.Join(projectDir, "build")

		compileOpts := cli.CompileOpts{
			FQBN:       expectedFqbn,
			SketchDir:  projectDir,
			SketchPath: expectedSketch,
			BuildProps: expectedBuildProps,
			ShowProps:  expectedShowProps,
		}

		instance := &rpc.Instance{Id: int32(1)}
		req := &rpc.CompileRequest{
			Instance:        instance,
			Fqbn:            expectedFqbn,
			SketchPath:      expectedSketch,
			ExportDir:       expectedExportDir,
			BuildProperties: expectedBuildProps,
			ShowProperties:  expectedShowProps,
			Verbose:         true,
		}

		env.Cli.EXPECT().CreateInstance().Return(instance, nil)
		env.Cli.EXPECT().Compile(gomock.Any(), req, gomock.Any(), gomock.Any(), gomock.Any())

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
		expectedExportDir := path.Join(projectDir, "build")

		compileOpts := cli.CompileOpts{
			FQBN:       expectedFqbn,
			SketchDir:  projectDir,
			SketchPath: expectedSketch,
			BuildProps: expectedBuildProps,
			ShowProps:  expectedShowProps,
		}

		instance := &rpc.Instance{Id: int32(1)}
		req := &rpc.CompileRequest{
			Instance:        instance,
			Fqbn:            expectedFqbn,
			SketchPath:      expectedSketch,
			ExportDir:       expectedExportDir,
			BuildProperties: expectedBuildProps,
			ShowProperties:  expectedShowProps,
			Verbose:         true,
		}

		env.Cli.EXPECT().CreateInstance().Return(instance, nil)
		env.Cli.EXPECT().Compile(gomock.Any(), req, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dummyErr)

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
		expectedExportDir := path.Join(sketchDir, "build")

		compileOpts := cli.CompileOpts{
			FQBN:       expectedFqbn,
			SketchDir:  projectDir,
			SketchPath: expectedSketch,
			BuildProps: expectedBuildProps,
			ShowProps:  expectedShowProps,
		}

		instance := &rpc.Instance{Id: int32(1)}
		req := &rpc.CompileRequest{
			Instance:        instance,
			Fqbn:            expectedFqbn,
			SketchPath:      expectedSketch,
			ExportDir:       expectedExportDir,
			BuildProperties: expectedBuildProps,
			ShowProperties:  expectedShowProps,
			Verbose:         true,
		}

		env.Cli.EXPECT().CreateInstance().Return(instance, nil)
		env.Cli.EXPECT().Compile(gomock.Any(), req, gomock.Any(), gomock.Any(), gomock.Any())

		err = env.ArdiCore.Compiler.Compile(compileOpts)
		assert.NoError(env.T, err)

		env.ClearStdout()
		env.Cli.EXPECT().CreateInstance().Return(instance, nil)
		env.Cli.EXPECT().Compile(gomock.Any(), req, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dummyErr)

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
