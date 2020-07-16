package commands_test

import (
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
