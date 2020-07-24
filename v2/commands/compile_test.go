package commands_test

import (
	"path"
	"testing"

	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCompileCommandGlobal(t *testing.T) {
	testutil.RunIntegrationTest("using globally installed platform from external url", t, func(groupEnv *testutil.IntegrationTestEnv) {
		boardURL := testutil.Esp8266BoardURL()
		platform := testutil.Esp8266Platform()
		fqbn := testutil.Esp8266WifiduinoFQBN()

		args := []string{"add", "board-url", boardURL, "--global"}
		err := groupEnv.Execute(args)
		assert.NoError(groupEnv.T, err)

		args = []string{"add", "platform", platform, "--global"}
		err = groupEnv.Execute(args)
		assert.NoError(groupEnv.T, err)

		groupEnv.T.Run("compiles project directory", func(st *testing.T) {
			testutil.CleanBuilds()
			blinkDir := testutil.BlinkProjectDir()
			args := []string{"compile", blinkDir, "--fqbn", fqbn, "--global"}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
		})

		groupEnv.T.Run("compiles ardi.json build", func(st *testing.T) {
			testutil.CleanBuilds()
			buildName := "blink"
			sketchDir := testutil.BlinkProjectDir()
			buildDir := path.Join(sketchDir, "build")

			args := []string{"add", "build", "-n", buildName, "-f", fqbn, "-s", sketchDir, "--global"}
			err := groupEnv.Execute(args)
			assert.NoError(st, err)

			args = []string{"compile", buildName, "--global"}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
			assert.DirExists(st, buildDir)
		})

		groupEnv.T.Run("compiles multiple ardi.json builds", func(st *testing.T) {
			testutil.CleanBuilds()
			buildName1 := "blink"
			sketchDir1 := testutil.BlinkProjectDir()
			buildDir1 := path.Join(sketchDir1, "build")

			args := []string{"add", "build", "-n", buildName1, "-f", fqbn, "-s", sketchDir1, "--global"}
			err := groupEnv.Execute(args)
			assert.NoError(st, err)

			buildName2 := "blink2"
			sketchDir2 := testutil.BlinkCopyProjectDir()
			buildDir2 := path.Join(sketchDir2, "build")

			args = []string{"add", "build", "-n", buildName2, "-f", fqbn, "-s", sketchDir2, "--global"}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)

			args = []string{"compile", buildName1, buildName2, "--global"}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
			assert.DirExists(st, buildDir1)
			assert.DirExists(st, buildDir2)
		})

		groupEnv.T.Run("errors if attempt to watch multiple builds", func(st *testing.T) {
			testutil.CleanBuilds()
			buildName1 := "blink"
			sketchDir1 := testutil.BlinkProjectDir()

			args := []string{"add", "build", "-n", buildName1, "-f", fqbn, "-s", sketchDir1, "--global"}
			err := groupEnv.Execute(args)
			assert.NoError(st, err)

			buildName2 := "blink2"
			sketchDir2 := testutil.BlinkCopyProjectDir()

			args = []string{"add", "build", "-n", buildName2, "-f", fqbn, "-s", sketchDir2, "--global"}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)

			args = []string{"compile", buildName1, buildName2, "--watch", "--global"}
			err = groupEnv.Execute(args)
			assert.Error(st, err)
		})

		groupEnv.T.Run("compiles all ardi.json builds", func(st *testing.T) {
			testutil.CleanBuilds()
			buildName := "blink"
			sketchDir := testutil.BlinkProjectDir()
			buildDir := path.Join(sketchDir, "build")

			args := []string{"add", "build", "-n", buildName, "-f", fqbn, "-s", sketchDir, "--global"}
			err := groupEnv.Execute(args)
			assert.NoError(st, err)

			args = []string{"compile", "--all", "--global"}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
			assert.DirExists(st, buildDir)
		})

		groupEnv.T.Run("errors if attempting to watch all builds", func(st *testing.T) {
			testutil.CleanBuilds()
			buildName := "blink"
			sketchDir := testutil.BlinkProjectDir()

			args := []string{"add", "build", "-n", buildName, "-f", fqbn, "-s", sketchDir, "--global"}
			err := groupEnv.Execute(args)
			assert.NoError(st, err)

			args = []string{"compile", "--all", "--watch", "--global"}
			err = groupEnv.Execute(args)
			assert.Error(st, err)
		})

		groupEnv.T.Run("errors if fqbn is missing", func(st *testing.T) {
			testutil.CleanBuilds()
			blinkDir := testutil.BlinkProjectDir()
			args := []string{"compile", blinkDir, "--global"}
			err = groupEnv.Execute(args)
			assert.Error(st, err)
		})

		groupEnv.T.Run("errors if global lib required", func(st *testing.T) {
			testutil.CleanBuilds()
			pixieDir := testutil.PixieProjectDir()
			args := []string{"compile", pixieDir, "--fqbn", testutil.ArduinoMegaFQBN(), "--global"}
			err = groupEnv.Execute(args)
			assert.Error(st, err)
		})

		groupEnv.T.Run("compiles project that requires globally installed library", func(st *testing.T) {
			testutil.CleanBuilds()
			args := []string{"add", "lib", "Adafruit Pixie", "--global"}
			err := groupEnv.Execute(args)
			assert.NoError(st, err)

			pixieDir := testutil.PixieProjectDir()
			args = []string{"compile", pixieDir, "--fqbn", fqbn, "--global"}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
		})
	})

	testutil.RunIntegrationTest("errors if platform not installed globally", t, func(env *testutil.IntegrationTestEnv) {
		blinkDir := testutil.BlinkProjectDir()
		args := []string{"compile", blinkDir, "--fqbn", testutil.ArduinoMegaFQBN(), "--global"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})
}

func TestCompileCommandProject(t *testing.T) {
	testutil.RunIntegrationTest("using installed project platform", t, func(groupEnv *testutil.IntegrationTestEnv) {
		err := groupEnv.RunProjectInit()
		assert.NoError(groupEnv.T, err)

		boardURL := testutil.Esp8266BoardURL()
		platform := testutil.Esp8266Platform()
		fqbn := testutil.Esp8266WifiduinoFQBN()

		args := []string{"add", "board-url", boardURL}
		err = groupEnv.Execute(args)
		assert.NoError(groupEnv.T, err)

		args = []string{"add", "platform", platform}
		err = groupEnv.Execute(args)
		assert.NoError(groupEnv.T, err)

		groupEnv.T.Run("compiles project directory", func(st *testing.T) {
			testutil.CleanBuilds()
			blinkDir := testutil.BlinkProjectDir()
			args := []string{"compile", blinkDir, "--fqbn", fqbn}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
		})

		groupEnv.T.Run("compiles ardi.json build", func(st *testing.T) {
			testutil.CleanBuilds()
			buildName := "blink"
			sketchDir := testutil.BlinkProjectDir()
			buildDir := path.Join(sketchDir, "build")

			args := []string{"add", "build", "-n", buildName, "-f", fqbn, "-s", sketchDir}
			err := groupEnv.Execute(args)
			assert.NoError(st, err)

			args = []string{"compile", buildName}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
			assert.DirExists(st, buildDir)
		})

		groupEnv.T.Run("compiles multiple ardi.json builds", func(st *testing.T) {
			testutil.CleanBuilds()
			buildName1 := "blink"
			sketchDir1 := testutil.BlinkProjectDir()
			buildDir1 := path.Join(sketchDir1, "build")

			args := []string{"add", "build", "-n", buildName1, "-f", fqbn, "-s", sketchDir1}
			err := groupEnv.Execute(args)
			assert.NoError(st, err)

			buildName2 := "blink2"
			sketchDir2 := testutil.BlinkCopyProjectDir()
			buildDir2 := path.Join(sketchDir2, "build")

			args = []string{"add", "build", "-n", buildName2, "-f", fqbn, "-s", sketchDir2}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)

			args = []string{"compile", buildName1, buildName2}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
			assert.DirExists(st, buildDir1)
			assert.DirExists(st, buildDir2)
		})

		groupEnv.T.Run("errors if attempt to watch multiple builds", func(st *testing.T) {
			testutil.CleanBuilds()
			buildName1 := "blink"
			sketchDir1 := testutil.BlinkProjectDir()

			args := []string{"add", "build", "-n", buildName1, "-f", fqbn, "-s", sketchDir1}
			err := groupEnv.Execute(args)
			assert.NoError(st, err)

			buildName2 := "blink2"
			sketchDir2 := testutil.BlinkCopyProjectDir()

			args = []string{"add", "build", "-n", buildName2, "-f", fqbn, "-s", sketchDir2}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)

			args = []string{"compile", buildName1, buildName2, "--watch"}
			err = groupEnv.Execute(args)
			assert.Error(st, err)
		})

		groupEnv.T.Run("compiles all ardi.json builds", func(st *testing.T) {
			testutil.CleanBuilds()
			buildName := "blink"
			sketchDir := testutil.BlinkProjectDir()
			buildDir := path.Join(sketchDir, "build")

			args := []string{"add", "build", "-n", buildName, "-f", fqbn, "-s", sketchDir}
			err := groupEnv.Execute(args)
			assert.NoError(st, err)

			args = []string{"compile", "--all"}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
			assert.DirExists(st, buildDir)
		})

		groupEnv.T.Run("errors if attempting to watch all builds", func(st *testing.T) {
			testutil.CleanBuilds()
			buildName := "blink"
			sketchDir := testutil.BlinkProjectDir()

			args := []string{"add", "build", "-n", buildName, "-f", fqbn, "-s", sketchDir}
			err := groupEnv.Execute(args)
			assert.NoError(st, err)

			args = []string{"compile", "--all", "--watch"}
			err = groupEnv.Execute(args)
			assert.Error(st, err)
		})

		groupEnv.T.Run("errors if fqbn is missing", func(st *testing.T) {
			testutil.CleanBuilds()
			blinkDir := testutil.BlinkProjectDir()
			args := []string{"compile", blinkDir}
			err = groupEnv.Execute(args)
			assert.Error(st, err)
		})

		groupEnv.T.Run("errors if project library required", func(st *testing.T) {
			testutil.CleanBuilds()
			pixieDir := testutil.PixieProjectDir()
			args := []string{"compile", pixieDir, "--fqbn", fqbn}
			err = groupEnv.Execute(args)
			assert.Error(st, err)
		})

		groupEnv.T.Run("compiles project that requires project library", func(st *testing.T) {
			testutil.CleanBuilds()
			args := []string{"add", "lib", "Adafruit Pixie"}
			err := groupEnv.Execute(args)
			assert.NoError(st, err)

			pixieDir := testutil.PixieProjectDir()
			args = []string{"compile", pixieDir, "--fqbn", fqbn}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
		})
	})

	testutil.RunIntegrationTest("errors if platform not installed for project", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)
		blinkDir := testutil.BlinkProjectDir()
		args := []string{"compile", blinkDir, "--fqbn", testutil.ArduinoMegaFQBN()}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("errors if not a valid project directory", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"compile", ".", "--fqbn", testutil.ArduinoMegaFQBN()}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})
}
