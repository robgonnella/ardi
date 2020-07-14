package core_test

import (
	"fmt"
	"path"
	"testing"

	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

// @todo: check that list is actually sorted
func TestProjectCore(t *testing.T) {
	testutil.RunUnitTest("init creates ardi.json", t, func(env *testutil.UnitTestEnv) {
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(env.T, err)
		assert.FileExists(env.T, "ardi.json")
	})

	testutil.RunUnitTest("init creates .ardi directory", t, func(env *testutil.UnitTestEnv) {
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(env.T, err)
		assert.DirExists(env.T, ".ardi")
		assert.FileExists(env.T, ".ardi/arduino-cli.yaml")
	})

	testutil.RunUnitTest("adds library to ardi.json", t, func(env *testutil.UnitTestEnv) {
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(env.T, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(env.T, err)

		lib := "some-lib"
		vers := "1.0.0"

		err = env.ArdiCore.Project.AddLibrary(lib, vers)
		assert.NoError(env.T, err)

		libs := env.ArdiCore.Project.GetLibraries()
		assert.Contains(env.T, libs, lib)
		assert.Equal(env.T, libs[lib], vers)
	})

	testutil.RunUnitTest("removes library from ardi.json", t, func(env *testutil.UnitTestEnv) {
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(env.T, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(env.T, err)

		lib := "some-lib"
		vers := "1.0.0"

		err = env.ArdiCore.Project.AddLibrary(lib, vers)
		assert.NoError(env.T, err)

		libs := env.ArdiCore.Project.GetLibraries()
		assert.Contains(env.T, libs, lib)

		err = env.ArdiCore.Project.RemoveLibrary(lib)
		assert.NoError(env.T, err)

		libs = env.ArdiCore.Project.GetLibraries()
		assert.NotContains(env.T, libs, lib)
	})

	testutil.RunUnitTest("lists libraries in ardi.json", t, func(env *testutil.UnitTestEnv) {
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(env.T, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(env.T, err)

		lib := "some-lib"
		vers := "1.0.0"

		err = env.ArdiCore.Project.AddLibrary(lib, vers)
		assert.NoError(env.T, err)

		env.ClearStdout()
		env.ArdiCore.Project.ListLibraries()
		assert.Contains(env.T, env.Stdout.String(), lib)
	})

	testutil.RunUnitTest("adds platform to ardi.json", t, func(env *testutil.UnitTestEnv) {
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(env.T, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(env.T, err)

		platform := "some-platform"
		vers := "1.0.0"

		err = env.ArdiCore.Project.AddPlatform(platform, vers)
		assert.NoError(env.T, err)

		plats := env.ArdiCore.Project.GetPlatforms()
		assert.Contains(env.T, plats, platform)
		assert.Equal(env.T, plats[platform], vers)
	})

	testutil.RunUnitTest("removes platform from ardi.json", t, func(env *testutil.UnitTestEnv) {
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(env.T, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(env.T, err)

		platform := "some-platform"
		vers := "1.0.0"

		err = env.ArdiCore.Project.AddPlatform(platform, vers)
		assert.NoError(env.T, err)

		plats := env.ArdiCore.Project.GetPlatforms()
		assert.Contains(env.T, plats, platform)
		assert.Equal(env.T, plats[platform], vers)

		err = env.ArdiCore.Project.RemovePlatform(platform)
		assert.NoError(env.T, err)

		plats = env.ArdiCore.Project.GetPlatforms()
		assert.NotContains(env.T, plats, platform)
	})

	testutil.RunUnitTest("lists platforms in ardi.json", t, func(env *testutil.UnitTestEnv) {
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(env.T, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(env.T, err)

		platform := "some-platform"
		vers := "1.0.0"

		err = env.ArdiCore.Project.AddPlatform(platform, vers)
		assert.NoError(env.T, err)

		env.ClearStdout()
		env.ArdiCore.Project.ListPlatforms()
		assert.Contains(env.T, env.Stdout.String(), platform)
	})

	testutil.RunUnitTest("adds board url", t, func(env *testutil.UnitTestEnv) {
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(env.T, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(env.T, err)

		boardURL := "https://someboardurl.com"

		err = env.ArdiCore.Project.AddBoardURL(boardURL)
		assert.NoError(env.T, err)

		configURLS := env.ArdiCore.Project.GetBoardURLS()
		dataURLS := env.ArdiCore.Project.GetDataConfig().BoardManager.AdditionalUrls

		assert.Contains(env.T, configURLS, boardURL)
		assert.Contains(env.T, dataURLS, boardURL)
	})

	testutil.RunUnitTest("removes board url", t, func(env *testutil.UnitTestEnv) {
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(env.T, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(env.T, err)

		boardURL := "https://someboardurl.com"

		err = env.ArdiCore.Project.AddBoardURL(boardURL)
		assert.NoError(env.T, err)

		configURLS := env.ArdiCore.Project.GetBoardURLS()
		dataURLS := env.ArdiCore.Project.GetDataConfig().BoardManager.AdditionalUrls

		assert.Contains(env.T, configURLS, boardURL)
		assert.Contains(env.T, dataURLS, boardURL)

		err = env.ArdiCore.Project.RemoveBoardURL(boardURL)
		assert.NoError(env.T, err)

		dataURLS = env.ArdiCore.Project.GetDataConfig().BoardManager.AdditionalUrls
		configURLS = env.ArdiCore.Project.GetBoardURLS()
		assert.NotContains(env.T, dataURLS, boardURL)
		assert.NotContains(env.T, configURLS, boardURL)
	})

	testutil.RunUnitTest("lists board urls", t, func(env *testutil.UnitTestEnv) {
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(env.T, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(env.T, err)

		boardURL := "https://someboardurl.com"

		err = env.ArdiCore.Project.AddBoardURL(boardURL)
		assert.NoError(env.T, err)

		env.ClearStdout()
		env.ArdiCore.Project.ListBoardURLS()
		assert.Contains(env.T, env.Stdout.String(), boardURL)
	})

	testutil.RunUnitTest("adds build to ardi.json", t, func(env *testutil.UnitTestEnv) {
		port := "2222"
		buildName := "blink"
		platform := "arduino-platform"
		boardURL := "https://some-board-url.com"
		projectPath := testutil.BlinkProjectDir()
		fqbn := "testboardfqbb"
		buildProp := "some_build_prop"
		buildPropVal := "DTest"
		buildProps := []string{fmt.Sprintf("%s=%s", buildProp, buildPropVal)}

		err := env.ArdiCore.Project.Init(port)
		assert.NoError(env.T, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(env.T, err)

		env.Client.EXPECT().InstallPlatform(platform).Times(1).Return(platform, "1.1.1", nil)
		err = env.ArdiCore.Project.AddBuild(buildName, platform, boardURL, projectPath, fqbn, buildProps)
		assert.NoError(env.T, err)

		builds := env.ArdiCore.Project.GetBuilds()
		dataConfig := env.ArdiCore.Project.GetDataConfig()

		assert.Contains(env.T, builds, buildName)
		assert.Equal(env.T, builds[buildName].Path, projectPath)
		assert.Equal(env.T, builds[buildName].FQBN, fqbn)
		assert.Contains(env.T, builds[buildName].Props, buildProp)
		assert.Equal(env.T, builds[buildName].Props[buildProp], buildPropVal)
		assert.Contains(env.T, dataConfig.BoardManager.AdditionalUrls, boardURL)
		assert.Equal(env.T, dataConfig.Daemon.Port, port)
	})

	testutil.RunUnitTest("removes build from ardi.json", t, func(env *testutil.UnitTestEnv) {
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(env.T, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(env.T, err)

		buildName := "blink"
		platform := "arduino-platform"
		boardURL := "https://some-board-url.com"
		path := testutil.BlinkProjectDir()
		fqbn := "testboardfqbb"
		buildProp := "some_build_prop"
		buildPropVal := "DTest"
		buildProps := []string{fmt.Sprintf("%s=%s", buildProp, buildPropVal)}

		env.Client.EXPECT().InstallPlatform(platform).Times(1).Return(platform, "1.1.1", nil)
		err = env.ArdiCore.Project.AddBuild(buildName, platform, boardURL, path, fqbn, buildProps)
		assert.NoError(env.T, err)

		builds := env.ArdiCore.Project.GetBuilds()
		assert.Contains(env.T, builds, buildName)
		assert.Equal(env.T, builds[buildName].Path, path)
		assert.Equal(env.T, builds[buildName].FQBN, fqbn)
		assert.Contains(env.T, builds[buildName].Props, buildProp)
		assert.Equal(env.T, builds[buildName].Props[buildProp], buildPropVal)

		env.ArdiCore.Project.RemoveBuild(buildName)
		builds = env.ArdiCore.Project.GetBuilds()
		assert.NotContains(env.T, builds, buildName)
	})

	testutil.RunUnitTest("lists builds in ardi.json", t, func(env *testutil.UnitTestEnv) {
		port := "2222"
		buildName := "blink"
		platform := "arduino-platform"
		boardURL := "https://some-board-url.com"
		projectPath := testutil.BlinkProjectDir()
		fqbn := "testboardfqbb"
		buildProp := "some_build_prop"
		buildPropVal := "DTest"
		buildProps := []string{fmt.Sprintf("%s=%s", buildProp, buildPropVal)}

		err := env.ArdiCore.Project.Init(port)
		assert.NoError(env.T, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(env.T, err)

		env.Client.EXPECT().InstallPlatform(platform).Times(1).Return(platform, "1.1.1", nil)
		err = env.ArdiCore.Project.AddBuild(buildName, platform, boardURL, projectPath, fqbn, buildProps)
		assert.NoError(env.T, err)

		env.ClearStdout()
		env.ArdiCore.Project.ListBuilds([]string{})
		out := env.Stdout.String()
		assert.Contains(env.T, out, buildName)
		assert.Contains(env.T, out, projectPath)
		assert.Contains(env.T, out, fqbn)
		assert.Contains(env.T, out, buildProp)
		assert.Contains(env.T, out, buildPropVal)
	})

	testutil.RunUnitTest("errors if not initialized", t, func(env *testutil.UnitTestEnv) {
		err := env.ArdiCore.Project.SetConfigHelpers()
		assert.Error(env.T, err)

		err = env.ArdiCore.Project.AddLibrary("somelib", "1.0.0")
		assert.Error(env.T, err)

		err = env.ArdiCore.Project.RemoveLibrary("somelib")
		assert.Error(env.T, err)

		err = env.ArdiCore.Project.AddPlatform("someplatform", "1.2.3")
		assert.Error(env.T, err)

		err = env.ArdiCore.Project.RemovePlatform("someplatform")
		assert.Error(env.T, err)

		err = env.ArdiCore.Project.AddBoardURL("url")
		assert.Error(env.T, err)

		err = env.ArdiCore.Project.RemoveBoardURL("url")
		assert.Error(env.T, err)
	})

	testutil.RunUnitTest("builds project specified ardi.json", t, func(env *testutil.UnitTestEnv) {
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(env.T, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(env.T, err)

		buildName := "blink"
		platform := "arduino-platform"
		boardURL := "https://some-board-url.com"
		sketchDir := testutil.BlinkProjectDir()
		fqbn := "testboardfqbb"
		buildProp := "some_build_prop"
		buildPropVal := "DTest"
		buildProps := []string{fmt.Sprintf("%s=%s", buildProp, buildPropVal)}

		env.Client.EXPECT().InstallPlatform(platform).Times(1).Return(platform, "1.1.1", nil)
		err = env.ArdiCore.Project.AddBuild(buildName, platform, boardURL, sketchDir, fqbn, buildProps)
		assert.NoError(env.T, err)

		compileOpts := rpc.CompileOpts{
			SketchDir:  sketchDir,
			SketchPath: path.Join(sketchDir, "blink.ino"),
			FQBN:       fqbn,
			BuildProps: buildProps,
			ShowProps:  false,
			ExportName: buildName,
		}

		env.Client.EXPECT().Compile(compileOpts).Times(1).Return(nil)
		err = env.ArdiCore.Project.BuildAll()
		assert.NoError(env.T, err)
	})

	testutil.RunUnitTest("builds project by name", t, func(env *testutil.UnitTestEnv) {
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(env.T, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(env.T, err)

		buildName := "blink"
		platform := "arduino-platform"
		boardURL := "https://some-board-url.com"
		sketchDir := testutil.BlinkProjectDir()
		fqbn := "testboardfqbb"
		buildProp := "some_build_prop"
		buildPropVal := "DTest"
		buildProps := []string{fmt.Sprintf("%s=%s", buildProp, buildPropVal)}

		env.Client.EXPECT().InstallPlatform(platform).Times(1).Return(platform, "1.1.1", nil)
		err = env.ArdiCore.Project.AddBuild(buildName, platform, boardURL, sketchDir, fqbn, buildProps)
		assert.NoError(env.T, err)

		compileOpts := rpc.CompileOpts{
			SketchDir:  sketchDir,
			SketchPath: path.Join(sketchDir, "blink.ino"),
			FQBN:       fqbn,
			BuildProps: buildProps,
			ShowProps:  false,
			ExportName: buildName,
		}

		env.Client.EXPECT().Compile(compileOpts).Times(1).Return(nil)
		err = env.ArdiCore.Project.Build(buildName)
		assert.NoError(env.T, err)
	})

	testutil.RunUnitTest("errors if build doesn't exist", t, func(env *testutil.UnitTestEnv) {
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(env.T, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(env.T, err)
		err = env.ArdiCore.Project.Build("noop")
		assert.Error(env.T, err)
	})

	testutil.RunUnitTest("processes sketch", t, func(env *testutil.UnitTestEnv) {
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(env.T, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(env.T, err)

		err = env.ArdiCore.Project.ProcessSketch(testutil.BlinkProjectDir())
		assert.NoError(env.T, err)

		assert.NotEmpty(env.T, env.ArdiCore.Project.Directory)
		assert.NotEmpty(env.T, env.ArdiCore.Project.Sketch)
		assert.Equal(env.T, env.ArdiCore.Project.Baud, 9600)
	})
}
