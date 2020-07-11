package commands_test

import (
	"testing"

	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCompileCommandGlobal(t *testing.T) {
	testutil.RunIntegrationTest("using globally installed platform", t, func(groupEnv *testutil.IntegrationTestEnv) {
		err := groupEnv.AddPlatform("arduino:avr", testutil.GlobalOpt{true})
		assert.NoError(groupEnv.T, err)

		groupEnv.T.Run("compiles project directory", func(st *testing.T) {
			testutil.CleanBuilds()
			blinkDir := testutil.BlinkProjectDir()
			args := []string{"compile", blinkDir, "--fqbn", testutil.ArduinoMegaFQBN(), "--global"}
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
			err := groupEnv.AddLib("Adafruit Pixie", testutil.GlobalOpt{true})
			assert.NoError(st, err)
			pixieDir := testutil.PixieProjectDir()
			args := []string{"compile", pixieDir, "--fqbn", testutil.ArduinoMegaFQBN(), "--global"}
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
		err = groupEnv.AddPlatform("arduino:avr", testutil.GlobalOpt{false})
		assert.NoError(groupEnv.T, err)

		groupEnv.T.Run("compiles project directory", func(st *testing.T) {
			testutil.CleanBuilds()
			blinkDir := testutil.BlinkProjectDir()
			args := []string{"compile", blinkDir, "--fqbn", testutil.ArduinoMegaFQBN()}
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
			args := []string{"compile", pixieDir, "--fqbn", testutil.ArduinoMegaFQBN()}
			err = groupEnv.Execute(args)
			assert.Error(st, err)
		})

		groupEnv.T.Run("compiles project that requires project library", func(st *testing.T) {
			testutil.CleanBuilds()
			err := groupEnv.AddLib("Adafruit Pixie", testutil.GlobalOpt{false})
			assert.NoError(st, err)
			pixieDir := testutil.PixieProjectDir()
			args := []string{"compile", pixieDir, "--fqbn", testutil.ArduinoMegaFQBN()}
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
