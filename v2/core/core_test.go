package core_test

import (
	"path"
	"testing"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestArdiCore(t *testing.T) {
	testutil.RunUnitTest("returns default build compile opts", t, func(env *testutil.UnitTestEnv) {
		defaultBuild := "default"
		defaultSketchDir := testutil.BlinkProjectDir()
		defaultSketch := path.Join(defaultSketchDir, "blink.ino")
		defaultFQBN := "someboardfqbn"

		anotherBuild := "another"
		anotherSketchDir := testutil.PixieProjectDir()
		anotherSketch := path.Join(anotherSketchDir, "pixie.ino")
		anotherFQBN := "anotherfqbn"

		expectedCompileOpts := &cli.CompileOpts{
			FQBN:       defaultFQBN,
			SketchDir:  defaultSketchDir,
			SketchPath: defaultSketch,
			BuildProps: []string{},
			ShowProps:  false,
		}

		err := env.ArdiCore.Config.AddBuild(defaultBuild, defaultSketch, defaultFQBN, 0, []string{})
		assert.NoError(env.T, err)

		err = env.ArdiCore.Config.AddBuild(anotherBuild, anotherSketch, anotherFQBN, 0, []string{})
		assert.NoError(env.T, err)

		opts, err := env.ArdiCore.GetCompileOptsFromArgs("", []string{}, false, []string{})
		assert.NoError(env.T, err)
		assert.Equal(env.T, expectedCompileOpts, opts[0])
	})

	testutil.RunUnitTest("returns specified build build compile opts", t, func(env *testutil.UnitTestEnv) {
		defaultBuild := "default"
		defaultSketchDir := testutil.BlinkProjectDir()
		defaultSketch := path.Join(defaultSketchDir, "blink.ino")
		defaultFQBN := "someboardfqbn"

		anotherBuild := "another"
		anotherSketchDir := testutil.PixieProjectDir()
		anotherSketch := path.Join(anotherSketchDir, "pixie.ino")
		anotherFQBN := "anotherfqbn"

		expectedCompileOpts := &cli.CompileOpts{
			FQBN:       anotherFQBN,
			SketchDir:  anotherSketchDir,
			SketchPath: anotherSketch,
			BuildProps: []string{},
			ShowProps:  false,
		}

		err := env.ArdiCore.Config.AddBuild(defaultBuild, defaultSketch, defaultFQBN, 0, []string{})
		assert.NoError(env.T, err)

		err = env.ArdiCore.Config.AddBuild(anotherBuild, anotherSketch, anotherFQBN, 0, []string{})
		assert.NoError(env.T, err)

		opts, err := env.ArdiCore.GetCompileOptsFromArgs("", []string{}, false, []string{"another"})
		assert.NoError(env.T, err)
		assert.Equal(env.T, expectedCompileOpts, opts[0])
	})

	testutil.RunUnitTest("returns single build compile opts", t, func(env *testutil.UnitTestEnv) {
		build := "default"
		sketchDir := testutil.BlinkProjectDir()
		sketch := path.Join(sketchDir, "blink.ino")
		fqbn := "someboardfqbn"

		expectedCompileOpts := &cli.CompileOpts{
			FQBN:       fqbn,
			SketchDir:  sketchDir,
			SketchPath: sketch,
			BuildProps: []string{},
			ShowProps:  false,
		}

		err := env.ArdiCore.Config.AddBuild(build, sketch, fqbn, 0, []string{})
		assert.NoError(env.T, err)

		opts, err := env.ArdiCore.GetCompileOptsFromArgs("", []string{}, false, []string{})
		assert.NoError(env.T, err)
		assert.Equal(env.T, expectedCompileOpts, opts[0])
	})

	testutil.RunUnitTest("returns .ino compile opts", t, func(env *testutil.UnitTestEnv) {
		sketchDir := testutil.BlinkProjectDir()
		sketch := path.Join(sketchDir, "blink.ino")
		fqbn := "someboardfqbn"

		expectedCompileOpts := &cli.CompileOpts{
			FQBN:       fqbn,
			SketchDir:  sketchDir,
			SketchPath: sketch,
			BuildProps: []string{},
			ShowProps:  false,
		}

		opts, err := env.ArdiCore.GetCompileOptsFromArgs(fqbn, []string{}, false, []string{sketch})
		assert.NoError(env.T, err)
		assert.Equal(env.T, expectedCompileOpts, opts[0])
	})

	testutil.RunUnitTest("returns error if .ino does not exist in current directory", t, func(env *testutil.UnitTestEnv) {
		opts, err := env.ArdiCore.GetCompileOptsFromArgs("", []string{}, false, []string{})
		assert.Error(env.T, err)
		assert.Nil(env.T, opts)
	})

	testutil.RunUnitTest("returns error if .ino does not exist", t, func(env *testutil.UnitTestEnv) {
		opts, err := env.ArdiCore.GetCompileOptsFromArgs("", []string{}, false, []string{"noop.ino"})
		assert.Error(env.T, err)
		assert.Nil(env.T, opts)
	})
}
