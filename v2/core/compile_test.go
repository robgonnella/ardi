package core_test

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/testutil"
)

func TestCompileCore(t *testing.T) {
	testutil.RunTest("errors when compiling directory with no sketch", t, func(st *testing.T, env testutil.TestEnv) {
		defer env.Ctrl.Finish()
		err := env.ArdiCore.Compiler.Compile(".", "some-fqbn", []string{}, false)
		assert.Error(st, err)
	})

	testutil.RunTest("succeeds when compiling directory with .ino file", t, func(st *testing.T, env testutil.TestEnv) {
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
		assert.NoError(st, err)
	})
}
