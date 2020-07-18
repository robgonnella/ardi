package core_test

import (
	"errors"
	"path"
	"testing"

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
			ExportName: "",
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
			ExportName: "",
		}

		env.Client.EXPECT().Compile(compileOpts).Times(1).Return(dummyErr)

		err := env.ArdiCore.Compiler.Compile(compileOpts)
		assert.Error(env.T, err)
		assert.EqualError(env.T, err, errString)
	})
}
