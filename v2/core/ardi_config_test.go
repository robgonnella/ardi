package core_test

import (
	"path"
	"testing"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/robgonnella/ardi/v2/util"
	"github.com/stretchr/testify/assert"
)

func TestArdiConfigBuilds(t *testing.T) {
	testutil.RunUnitTest("adds, lists, and removes builds", t, func(env *testutil.UnitTestEnv) {
		util.InitProjectDirectory()
		name1 := "somename"
		dir1 := testutil.BlinkProjectDir()
		name2 := "anothername"
		dir2 := testutil.Blink14400ProjectDir()
		fqbn := "somefqbn"
		buildProps := []string{"someprop=somevalue"}

		err := env.ArdiCore.Config.AddBuild(name1, dir1, fqbn, 0, buildProps)
		assert.NoError(env.T, err)

		err = env.ArdiCore.Config.AddBuild(name2, dir2, fqbn, 0, buildProps)
		assert.NoError(env.T, err)

		builds := env.ArdiCore.Config.GetBuilds()
		build, ok := builds[name1]
		assert.True(env.T, ok)
		assert.Equal(env.T, dir1, build.Directory)
		assert.Equal(env.T, fqbn, build.FQBN)
		assert.Equal(env.T, 9600, build.Baud)

		assert.Contains(env.T, build.Props, "someprop")
		assert.Equal(env.T, build.Props["someprop"], "somevalue")

		build, ok = builds[name2]
		assert.True(env.T, ok)
		assert.Equal(env.T, dir2, build.Directory)
		assert.Equal(env.T, fqbn, build.FQBN)
		assert.Equal(env.T, 14400, build.Baud)
		assert.Contains(env.T, build.Props, "someprop")
		assert.Equal(env.T, build.Props["someprop"], "somevalue")

		env.ClearStdout()
		env.ArdiCore.Config.ListBuilds([]string{})
		out := env.Stdout.String()
		assert.Contains(env.T, out, name1)
		assert.Contains(env.T, out, name2)
		assert.Contains(env.T, out, dir1)
		assert.Contains(env.T, out, dir2)
		assert.Contains(env.T, out, fqbn)
		assert.Contains(env.T, out, "someprop")
		assert.Contains(env.T, out, "somevalue")

		env.ClearStdout()
		env.ArdiCore.Config.ListBuilds([]string{name1})
		out = env.Stdout.String()
		assert.Contains(env.T, out, name1)
		assert.NotContains(env.T, out, name2)
		assert.Contains(env.T, out, dir1)
		assert.NotContains(env.T, out, dir2)
		assert.Contains(env.T, out, fqbn)
		assert.Contains(env.T, out, "someprop")
		assert.Contains(env.T, out, "somevalue")

		err = env.ArdiCore.Config.RemoveBuild(name2)
		assert.NoError(env.T, err)
		builds = env.ArdiCore.Config.GetBuilds()
		_, ok = builds[name2]
		assert.False(env.T, ok)
	})

	testutil.RunUnitTest("overrides baud", t, func(env *testutil.UnitTestEnv) {
		util.InitProjectDirectory()
		name := "somename"
		dir := testutil.Blink14400ProjectDir()
		fqbn := "somefqbn"
		buildProps := []string{}
		baud := 9600

		err := env.ArdiCore.Config.AddBuild(name, dir, fqbn, baud, buildProps)
		assert.NoError(env.T, err)

		builds := env.ArdiCore.Config.GetBuilds()
		build, ok := builds[name]
		assert.True(env.T, ok)
		assert.Equal(env.T, dir, build.Directory)
		assert.Equal(env.T, fqbn, build.FQBN)
		assert.Equal(env.T, baud, build.Baud)
	})

	testutil.RunUnitTest("errors if sketch not found", t, func(env *testutil.UnitTestEnv) {
		util.InitProjectDirectory()
		name := "somename"
		dir := "noop"
		fqbn := "somefqbn"
		buildProps := []string{"someprop=somevalue"}

		err := env.ArdiCore.Config.AddBuild(name, dir, fqbn, 0, buildProps)
		assert.Error(env.T, err)
	})

	testutil.RunUnitTest("allows relative paths", t, func(env *testutil.UnitTestEnv) {
		util.InitProjectDirectory()
		name := "somename"
		dir := "../test_projects/pixie"
		fqbn := "somefqbn"
		buildProps := []string{}
		baud := 9600

		err := env.ArdiCore.Config.AddBuild(name, dir, fqbn, baud, buildProps)
		assert.NoError(env.T, err)

		builds := env.ArdiCore.Config.GetBuilds()
		build, ok := builds[name]
		assert.True(env.T, ok)
		assert.Equal(env.T, dir, build.Directory)
		assert.Equal(env.T, fqbn, build.FQBN)
		assert.Equal(env.T, baud, build.Baud)
	})
}

func TestArdiConfigBoardURLS(t *testing.T) {
	testutil.RunUnitTest("adds, lists, and removes board urls", t, func(env *testutil.UnitTestEnv) {
		util.InitProjectDirectory()
		url := "https://someboardurl.com"

		err := env.ArdiCore.Config.AddBoardURL(url)
		assert.NoError(env.T, err)

		urls := env.ArdiCore.Config.GetBoardURLS()
		assert.Contains(env.T, urls, url)

		env.ClearStdout()
		env.ArdiCore.Config.ListBoardURLS()
		assert.Contains(env.T, env.Stdout.String(), url)

		err = env.ArdiCore.Config.RemoveBoardURL(url)
		assert.NoError(env.T, err)
		urls = env.ArdiCore.Config.GetBoardURLS()
		assert.NotContains(env.T, urls, url)
	})
}

func TestArdiConfigPlatform(t *testing.T) {
	testutil.RunUnitTest("adds, lists, and removes platforms", t, func(env *testutil.UnitTestEnv) {
		util.InitProjectDirectory()
		platform := "someplatform"
		vers := "1.4.3"

		err := env.ArdiCore.Config.AddPlatform(platform, vers)
		assert.NoError(env.T, err)

		platforms := env.ArdiCore.Config.GetPlatforms()
		assert.Contains(env.T, platforms, platform)
		assert.Equal(env.T, platforms[platform], vers)

		env.ClearStdout()
		env.ArdiCore.Config.ListPlatforms()
		out := env.Stdout.String()
		assert.Contains(env.T, out, platform)
		assert.Contains(env.T, out, vers)

		err = env.ArdiCore.Config.RemovePlatform(platform)
		assert.NoError(env.T, err)
		platforms = env.ArdiCore.Config.GetPlatforms()
		assert.NotContains(env.T, platforms, platform)
	})
}

func TestArdiConfigLibraries(t *testing.T) {
	testutil.RunUnitTest("adds, lists, and removes libraries", t, func(env *testutil.UnitTestEnv) {
		util.InitProjectDirectory()
		lib := "somelibrary"
		vers := "1.2.3"

		err := env.ArdiCore.Config.AddLibrary(lib, vers)
		assert.NoError(env.T, err)

		libraries := env.ArdiCore.Config.GetLibraries()
		assert.Contains(env.T, libraries, lib)
		assert.Equal(env.T, libraries[lib], vers)

		env.ClearStdout()
		env.ArdiCore.Config.ListLibraries()
		out := env.Stdout.String()
		assert.Contains(env.T, out, lib)
		assert.Contains(env.T, out, vers)

		err = env.ArdiCore.Config.RemoveLibrary(lib)
		assert.NoError(env.T, err)
		libraries = env.ArdiCore.Config.GetLibraries()
		assert.NotContains(env.T, libraries, lib)
	})
}

func TestArdiConfigCompileOpts(t *testing.T) {
	testutil.RunUnitTest("returns compile options for build", t, func(env *testutil.UnitTestEnv) {
		util.InitProjectDirectory()
		name := "somename"
		dir := testutil.BlinkProjectDir()
		fqbn := "somefqbn"
		buildProps := []string{"someprop=somevalue"}
		expectedOpts := &cli.CompileOpts{
			SketchDir:  dir,
			SketchPath: path.Join(dir, "blink.ino"),
			FQBN:       fqbn,
			BuildProps: buildProps,
		}

		err := env.ArdiCore.Config.AddBuild(name, dir, fqbn, 0, buildProps)
		assert.NoError(env.T, err)

		compileOpts, err := env.ArdiCore.Config.GetCompileOpts(name)
		assert.NoError(env.T, err)
		assert.Equal(env.T, expectedOpts, compileOpts)
	})

	testutil.RunUnitTest("errors if build not found", t, func(env *testutil.UnitTestEnv) {
		util.InitProjectDirectory()
		buildName := "noop"
		compileOpts, err := env.ArdiCore.Config.GetCompileOpts(buildName)
		assert.Error(env.T, err)
		assert.Nil(env.T, compileOpts)
	})
}
