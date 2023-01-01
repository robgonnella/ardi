package core_test

import (
	"errors"
	"path"
	"testing"

	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	cli "github.com/robgonnella/ardi/v3/cli-wrapper"
	"github.com/robgonnella/ardi/v3/testutil"
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

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), req, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

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

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), req, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dummyErr)

		err := env.ArdiCore.Compiler.Compile(compileOpts)
		assert.Error(env.T, err)
		assert.EqualError(env.T, err, errString)
	})

}
