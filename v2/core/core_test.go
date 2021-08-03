package core_test

import (
	"os"
	"path"
	"testing"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestArdiCoreGetCompileOptsFromArgs(t *testing.T) {
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

func TestArdiCoreGetBaudFromArgs(t *testing.T) {
	testutil.RunUnitTest("returns given baud", t, func(env *testutil.UnitTestEnv) {
		expectedBaud := 14400
		actualBaud := env.ArdiCore.GetBaudFromArgs(expectedBaud, []string{"doesn't matter"})
		assert.Equal(env.T, expectedBaud, actualBaud)
	})

	testutil.RunUnitTest("returns baud from single build", t, func(env *testutil.UnitTestEnv) {
		someBuild := "something"
		someSketchDir := testutil.PixieProjectDir()
		someSketch := path.Join(someSketchDir, "pixie.ino")
		someFQBN := "some:fq:bn"
		someBaud := 115200

		err := env.ArdiCore.Config.AddBuild(someBuild, someSketch, someFQBN, someBaud, []string{})
		assert.NoError(env.T, err)

		baud := env.ArdiCore.GetBaudFromArgs(0, []string{})
		assert.Equal(env.T, someBaud, baud)
	})

	testutil.RunUnitTest("returns baud from default build", t, func(env *testutil.UnitTestEnv) {
		defaultBuild := "default"
		defaultSketchDir := testutil.BlinkProjectDir()
		defaultSketch := path.Join(defaultSketchDir, "blink.ino")
		defaultFQBN := "someboardfqbn"
		defaultBaud := 14400

		anotherBuild := "another"
		anotherSketchDir := testutil.PixieProjectDir()
		anotherSketch := path.Join(anotherSketchDir, "pixie.ino")
		anotherFQBN := "anotherfqbn"
		anotherBaud := 115200

		err := env.ArdiCore.Config.AddBuild(defaultBuild, defaultSketch, defaultFQBN, defaultBaud, []string{})
		assert.NoError(env.T, err)

		err = env.ArdiCore.Config.AddBuild(anotherBuild, anotherSketch, anotherFQBN, anotherBaud, []string{})
		assert.NoError(env.T, err)

		baud := env.ArdiCore.GetBaudFromArgs(0, []string{})
		assert.Equal(env.T, defaultBaud, baud)
	})

	testutil.RunUnitTest("returns baud from current directory sketch", t, func(env *testutil.UnitTestEnv) {
		cwd, _ := os.Getwd()
		os.Chdir(testutil.Blink14400ProjectDir())
		defer os.Chdir(cwd)
		baud := env.ArdiCore.GetBaudFromArgs(0, []string{})
		assert.Equal(env.T, baud, 14400)
	})

	testutil.RunUnitTest("returns baud from defined build", t, func(env *testutil.UnitTestEnv) {
		myBuild := "mybuild"
		mySketchDir := testutil.PixieProjectDir()
		mySketch := path.Join(mySketchDir, "pixie.ino")
		myFQBN := "some:fq:bn"
		myBaud := 256000

		anotherBuild := "another"
		anotherSketchDir := testutil.PixieProjectDir()
		anotherSketch := path.Join(anotherSketchDir, "pixie.ino")
		anotherFQBN := "anotherfqbn"
		anotherBaud := 115200

		err := env.ArdiCore.Config.AddBuild(myBuild, mySketch, myFQBN, myBaud, []string{})
		assert.NoError(env.T, err)

		err = env.ArdiCore.Config.AddBuild(anotherBuild, anotherSketch, anotherFQBN, anotherBaud, []string{})
		assert.NoError(env.T, err)

		baud := env.ArdiCore.GetBaudFromArgs(0, []string{myBuild})
		assert.Equal(env.T, myBaud, baud)
	})

	testutil.RunUnitTest("returns baud from defined sketch", t, func(env *testutil.UnitTestEnv) {
		sketchDir := testutil.Blink14400ProjectDir()
		sketch := path.Join(sketchDir, "blink14400.ino")
		baud := env.ArdiCore.GetBaudFromArgs(0, []string{sketch})
		assert.Equal(env.T, 14400, baud)
	})

	testutil.RunUnitTest("returns default baud", t, func(env *testutil.UnitTestEnv) {
		baud := env.ArdiCore.GetBaudFromArgs(0, []string{})
		assert.Equal(env.T, 9600, baud)

		baud = env.ArdiCore.GetBaudFromArgs(0, []string{"noop1"})
		assert.Equal(env.T, 9600, baud)

		baud = env.ArdiCore.GetBaudFromArgs(0, []string{"noop1", "noop2"})
		assert.Equal(env.T, 9600, baud)
	})
}

func TestArdiCoreGetSketchPathsFromArgs(t *testing.T) {
	testutil.RunUnitTest("returns paths from default build", t, func(env *testutil.UnitTestEnv) {
		defaultBuild := "default"
		defaultSketchDir := testutil.BlinkProjectDir()
		defaultSketch := path.Join(defaultSketchDir, "blink.ino")
		defaultFQBN := "someboardfqbn"
		defaultBaud := 14400

		anotherBuild := "another"
		anotherSketchDir := testutil.PixieProjectDir()
		anotherSketch := path.Join(anotherSketchDir, "pixie.ino")
		anotherFQBN := "anotherfqbn"
		anotherBaud := 115200

		err := env.ArdiCore.Config.AddBuild(defaultBuild, defaultSketch, defaultFQBN, defaultBaud, []string{})
		assert.NoError(env.T, err)

		err = env.ArdiCore.Config.AddBuild(anotherBuild, anotherSketch, anotherFQBN, anotherBaud, []string{})
		assert.NoError(env.T, err)

		sketchDir, sketchPath, err := env.ArdiCore.GetSketchPathsFromArgs([]string{})
		assert.NoError(env.T, err)
		assert.Equal(env.T, defaultSketchDir, sketchDir)
		assert.Equal(env.T, defaultSketch, sketchPath)
	})

	testutil.RunUnitTest("returns paths from single build", t, func(env *testutil.UnitTestEnv) {
		someBuild := "something"
		someSketchDir := testutil.PixieProjectDir()
		someSketch := path.Join(someSketchDir, "pixie.ino")
		someFQBN := "some:fq:bn"
		someBaud := 115200

		err := env.ArdiCore.Config.AddBuild(someBuild, someSketch, someFQBN, someBaud, []string{})
		assert.NoError(env.T, err)

		sketchDir, sketchPath, err := env.ArdiCore.GetSketchPathsFromArgs([]string{})
		assert.NoError(env.T, err)
		assert.Equal(env.T, someSketchDir, sketchDir)
		assert.Equal(env.T, someSketch, sketchPath)
	})

	testutil.RunUnitTest("returns paths from directory sketch", t, func(env *testutil.UnitTestEnv) {
		blink14400Dir := "."
		blink14400Sketch := path.Join(blink14400Dir, "blink14400.ino")

		cwd, _ := os.Getwd()
		os.Chdir(testutil.Blink14400ProjectDir())
		defer os.Chdir(cwd)

		sketchDir, sketchPath, err := env.ArdiCore.GetSketchPathsFromArgs([]string{})
		assert.NoError(env.T, err)
		assert.Equal(env.T, blink14400Dir, sketchDir)
		assert.Equal(env.T, blink14400Sketch, sketchPath)
	})

	testutil.RunUnitTest("returns paths from defined build", t, func(env *testutil.UnitTestEnv) {
		myBuild := "mybuild"
		mySketchDir := testutil.PixieProjectDir()
		mySketch := path.Join(mySketchDir, "pixie.ino")
		myFQBN := "some:fq:bn"
		myBaud := 256000

		anotherBuild := "another"
		anotherSketchDir := testutil.PixieProjectDir()
		anotherSketch := path.Join(anotherSketchDir, "pixie.ino")
		anotherFQBN := "anotherfqbn"
		anotherBaud := 115200

		err := env.ArdiCore.Config.AddBuild(myBuild, mySketch, myFQBN, myBaud, []string{})
		assert.NoError(env.T, err)

		err = env.ArdiCore.Config.AddBuild(anotherBuild, anotherSketch, anotherFQBN, anotherBaud, []string{})
		assert.NoError(env.T, err)

		sketchDir, sketchPath, err := env.ArdiCore.GetSketchPathsFromArgs([]string{myBuild})
		assert.NoError(env.T, err)
		assert.Equal(env.T, mySketchDir, sketchDir)
		assert.Equal(env.T, mySketch, sketchPath)
	})

	testutil.RunUnitTest("returns paths from defined sketch", t, func(env *testutil.UnitTestEnv) {
		theSketchDir := testutil.PixieProjectDir()
		theSketch := path.Join(theSketchDir, "pixie.ino")
		sketchDir, sketchPath, err := env.ArdiCore.GetSketchPathsFromArgs([]string{theSketch})
		assert.NoError(env.T, err)
		assert.Equal(env.T, theSketchDir, sketchDir)
		assert.Equal(env.T, theSketch, sketchPath)
	})

	testutil.RunUnitTest("returns error if .ino does not exist in current directory", t, func(env *testutil.UnitTestEnv) {
		sketchDir, sketchPath, err := env.ArdiCore.GetSketchPathsFromArgs([]string{})
		assert.Error(env.T, err)
		assert.Empty(env.T, sketchDir)
		assert.Empty(env.T, sketchPath)
	})

	testutil.RunUnitTest("returns error if .ino does not exist", t, func(env *testutil.UnitTestEnv) {
		sketchDir, sketchPath, err := env.ArdiCore.GetSketchPathsFromArgs([]string{"../../noop.ino"})
		assert.Error(env.T, err)
		assert.Empty(env.T, sketchDir)
		assert.Empty(env.T, sketchPath)
	})
}
