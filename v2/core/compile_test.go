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
	testutil.RunUnitTest("errors when compiling directory with no sketch", t, func(env testutil.UnitTestEnv) {
		defer env.Ctrl.Finish()
		err := env.ArdiCore.Compiler.Compile(".", "some-fqbn", []string{}, false)
		assert.Error(env.T, err)
	})

	testutil.RunUnitTest("succeeds when compiling directory with .ino file", t, func(env testutil.UnitTestEnv) {
		defer env.Ctrl.Finish()

		expectedFqbn := "some-fqbb"
		expectedSketch := path.Join(env.BlinkProjDir, "blink.ino")
		expectedBuildProps := []string{"build.extra_flags='-DSOME_OPTION'"}
		expectedShowProps := false

		compileOpts := rpc.CompileOpts{
			FQBN:       expectedFqbn,
			SketchDir:  env.BlinkProjDir,
			SketchPath: expectedSketch,
			BuildProps: expectedBuildProps,
			ShowProps:  expectedShowProps,
		}

		env.Client.EXPECT().ConnectedBoards().Times(1).Return([]*rpc.Board{})
		env.Client.EXPECT().AllBoards().Times(1).Return([]*rpc.Board{})
		env.Client.EXPECT().Compile(compileOpts).Times(1)

		err := env.ArdiCore.Compiler.Compile(env.BlinkProjDir, expectedFqbn, expectedBuildProps, expectedShowProps)
		assert.NoError(env.T, err)
	})

	testutil.RunUnitTest("returns compile error", t, func(env testutil.UnitTestEnv) {
		defer env.Ctrl.Finish()
		errString := "dummy error"
		dummyErr := errors.New(errString)

		expectedFqbn := "some-fqbb"
		expectedSketch := path.Join(env.BlinkProjDir, "blink.ino")
		expectedBuildProps := []string{"build.extra_flags='-DSOME_OPTION'"}
		expectedShowProps := false

		compileOpts := rpc.CompileOpts{
			FQBN:       expectedFqbn,
			SketchDir:  env.BlinkProjDir,
			SketchPath: expectedSketch,
			BuildProps: expectedBuildProps,
			ShowProps:  expectedShowProps,
		}

		env.Client.EXPECT().ConnectedBoards().Times(1).Return([]*rpc.Board{})
		env.Client.EXPECT().AllBoards().Times(1).Return([]*rpc.Board{})
		env.Client.EXPECT().Compile(compileOpts).Times(1).Return(dummyErr)

		err := env.ArdiCore.Compiler.Compile(env.BlinkProjDir, expectedFqbn, expectedBuildProps, expectedShowProps)
		assert.Error(env.T, err)
		assert.EqualError(env.T, err, errString)
	})
}
