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
	testutil.RunTest("init creates ardi.json", t, func(st *testing.T, env testutil.TestEnv) {
		defer env.Ctrl.Finish()
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(st, err)
		assert.FileExists(st, "ardi.json")
	})

	testutil.RunTest("init creates .ardi directory", t, func(st *testing.T, env testutil.TestEnv) {
		defer env.Ctrl.Finish()
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(st, err)
		assert.DirExists(st, ".ardi")
		assert.FileExists(st, ".ardi/arduino-cli.yaml")
	})

	testutil.RunTest("adds library to ardi.json", t, func(st *testing.T, env testutil.TestEnv) {
		defer env.Ctrl.Finish()
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(st, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(st, err)

		lib := "some-lib"
		vers := "1.0.0"

		err = env.ArdiCore.Project.AddLibrary(lib, vers)
		assert.NoError(st, err)

		libs := env.ArdiCore.Project.GetLibraries()
		assert.Contains(st, libs, lib)
		assert.Equal(st, libs[lib], vers)
	})

	testutil.RunTest("processes sketch", t, func(st *testing.T, env testutil.TestEnv) {
		defer env.Ctrl.Finish()
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(st, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(st, err)

		err = env.ArdiCore.Project.ProcessSketch(env.BlinkProjDir)
		assert.NoError(st, err)

		assert.NotEmpty(st, env.ArdiCore.Project.Directory)
		assert.NotEmpty(st, env.ArdiCore.Project.Sketch)
		assert.Equal(st, env.ArdiCore.Project.Baud, 9600)
	})

	testutil.RunTest("removes library from ardi.json", t, func(st *testing.T, env testutil.TestEnv) {
		defer env.Ctrl.Finish()
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(st, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(st, err)

		lib := "some-lib"
		vers := "1.0.0"

		err = env.ArdiCore.Project.AddLibrary(lib, vers)
		assert.NoError(st, err)

		libs := env.ArdiCore.Project.GetLibraries()
		assert.Contains(st, libs, lib)

		err = env.ArdiCore.Project.RemoveLibrary(lib)
		assert.NoError(st, err)

		libs = env.ArdiCore.Project.GetLibraries()
		assert.NotContains(st, libs, lib)
	})

	testutil.RunTest("lists libraries in ardi.json", t, func(st *testing.T, env testutil.TestEnv) {
		defer env.Ctrl.Finish()
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(st, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(st, err)

		lib := "some-lib"
		vers := "1.0.0"

		err = env.ArdiCore.Project.AddLibrary(lib, vers)
		assert.NoError(st, err)

		libs := env.ArdiCore.Project.GetLibraries()
		assert.Contains(st, libs, lib)
		assert.Equal(st, libs[lib], vers)
	})

	testutil.RunTest("adds build to ardi.json", t, func(st *testing.T, env testutil.TestEnv) {
		defer env.Ctrl.Finish()
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(st, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(st, err)

		buildName := "blink"
		platform := "arduino-platform"
		boardURL := "https://some-board-url.com"
		path := env.BlinkProjDir
		fqbn := "testboardfqbb"
		buildProp := "some_build_prop"
		buildPropVal := "DTest"
		buildProps := []string{fmt.Sprintf("%s=%s", buildProp, buildPropVal)}

		env.Client.EXPECT().InstallPlatform(platform).Times(1).Return(nil)
		env.ArdiCore.Project.AddBuild(buildName, platform, boardURL, path, fqbn, buildProps)
		builds := env.ArdiCore.Project.GetBuilds()

		assert.Contains(st, builds, buildName)
		assert.Equal(st, builds[buildName].Platform, platform)
		assert.Equal(st, builds[buildName].BoardURL, boardURL)
		assert.Equal(st, builds[buildName].Path, path)
		assert.Equal(st, builds[buildName].FQBN, fqbn)
		assert.Contains(st, builds[buildName].Props, buildProp)
		assert.Equal(st, builds[buildName].Props[buildProp], buildPropVal)
	})

	testutil.RunTest("removes build from ardi.json", t, func(st *testing.T, env testutil.TestEnv) {
		defer env.Ctrl.Finish()
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(st, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(st, err)

		buildName := "blink"
		platform := "arduino-platform"
		boardURL := "https://some-board-url.com"
		path := env.BlinkProjDir
		fqbn := "testboardfqbb"
		buildProp := "some_build_prop"
		buildPropVal := "DTest"
		buildProps := []string{fmt.Sprintf("%s=%s", buildProp, buildPropVal)}

		env.Client.EXPECT().InstallPlatform(platform).Times(1).Return(nil)
		env.ArdiCore.Project.AddBuild(buildName, platform, boardURL, path, fqbn, buildProps)
		builds := env.ArdiCore.Project.GetBuilds()

		assert.Contains(st, builds, buildName)
		assert.Equal(st, builds[buildName].Platform, platform)
		assert.Equal(st, builds[buildName].BoardURL, boardURL)
		assert.Equal(st, builds[buildName].Path, path)
		assert.Equal(st, builds[buildName].FQBN, fqbn)
		assert.Contains(st, builds[buildName].Props, buildProp)
		assert.Equal(st, builds[buildName].Props[buildProp], buildPropVal)

		env.ArdiCore.Project.RemoveBuild(buildName)
		builds = env.ArdiCore.Project.GetBuilds()
		assert.NotContains(st, builds, buildName)
	})

	testutil.RunTest("builds project specified ardi.json", t, func(st *testing.T, env testutil.TestEnv) {
		defer env.Ctrl.Finish()
		err := env.ArdiCore.Project.Init("2222")
		assert.NoError(st, err)
		err = env.ArdiCore.Project.SetConfigHelpers()
		assert.NoError(st, err)

		buildName := "blink"
		platform := "arduino-platform"
		boardURL := "https://some-board-url.com"
		sketchDir := env.BlinkProjDir
		fqbn := "testboardfqbb"
		buildProp := "some_build_prop"
		buildPropVal := "DTest"
		buildProps := []string{fmt.Sprintf("%s=%s", buildProp, buildPropVal)}

		env.Client.EXPECT().InstallPlatform(platform).Times(1).Return(nil)
		env.ArdiCore.Project.AddBuild(buildName, platform, boardURL, sketchDir, fqbn, buildProps)

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
		assert.NoError(st, err)
	})
}
