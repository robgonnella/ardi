package commands_test

import (
	"testing"

	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCompileCommandGlobal(t *testing.T) {
	testutil.RunIntegrationTest("compiles project directory using global platform", t, func(env *testutil.IntegrationTestEnv) {
		err := env.InstallAvrPlatform(testutil.GlobalOpt{true})
		assert.NoError(env.T, err)
		blinkDir := testutil.BlinkProjectDir()
		args := []string{"compile", blinkDir, "--fqbn", testutil.ArduinoMegaFQBN(), "--global"}
		env.SetArgs(args)
		err = env.RootCmd.ExecuteContext(env.Ctx)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("errors if fqbn is missing when using global platform", t, func(env *testutil.IntegrationTestEnv) {
		err := env.InstallAvrPlatform(testutil.GlobalOpt{true})
		assert.NoError(env.T, err)
		blinkDir := testutil.BlinkProjectDir()
		args := []string{"compile", blinkDir, "--global"}
		env.SetArgs(args)
		err = env.RootCmd.ExecuteContext(env.Ctx)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("errors if platform not installed globally", t, func(env *testutil.IntegrationTestEnv) {
		blinkDir := testutil.BlinkProjectDir()
		args := []string{"compile", blinkDir, "--fqbn", testutil.ArduinoMegaFQBN(), "--global"}
		env.SetArgs(args)
		err := env.RootCmd.ExecuteContext(env.Ctx)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("compiles project that requires globally installed library", t, func(env *testutil.IntegrationTestEnv) {
		err := env.InstallAvrPlatform(testutil.GlobalOpt{true})
		env.AddLib("Adafruit Pixie", testutil.GlobalOpt{true})
		pixieDir := testutil.PixieProjectDir()
		args := []string{"compile", pixieDir, "--fqbn", testutil.ArduinoMegaFQBN(), "--global"}
		env.SetArgs(args)
		err = env.RootCmd.ExecuteContext(env.Ctx)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("errors when compiling project that requires a globally installed library", t, func(env *testutil.IntegrationTestEnv) {
		err := env.InstallAvrPlatform(testutil.GlobalOpt{true})
		pixieDir := testutil.PixieProjectDir()
		args := []string{"compile", pixieDir, "--fqbn", testutil.ArduinoMegaFQBN(), "--global"}
		env.SetArgs(args)
		err = env.RootCmd.ExecuteContext(env.Ctx)
		assert.Error(env.T, err)
	})
}

func TestCompileCommandProject(t *testing.T) {
	testutil.RunIntegrationTest("compiles project directory using project platform", t, func(env *testutil.IntegrationTestEnv) {
		err := env.InstallAvrPlatform(testutil.GlobalOpt{false})
		assert.NoError(env.T, err)
		blinkDir := testutil.BlinkProjectDir()
		args := []string{"compile", blinkDir, "--fqbn", testutil.ArduinoMegaFQBN()}
		env.SetArgs(args)
		err = env.RootCmd.ExecuteContext(env.Ctx)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("errors if fqbn is missing when using project platform", t, func(env *testutil.IntegrationTestEnv) {
		err := env.InstallAvrPlatform(testutil.GlobalOpt{false})
		assert.NoError(env.T, err)
		blinkDir := testutil.BlinkProjectDir()
		args := []string{"compile", blinkDir}
		env.SetArgs(args)
		err = env.RootCmd.ExecuteContext(env.Ctx)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("errors if platform not installed for project", t, func(env *testutil.IntegrationTestEnv) {
		env.RunProjectInit()
		blinkDir := testutil.BlinkProjectDir()
		args := []string{"compile", blinkDir, "--fqbn", testutil.ArduinoMegaFQBN()}
		env.SetArgs(args)
		err := env.RootCmd.ExecuteContext(env.Ctx)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("errors if not a valid project directory", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"compile", ".", "--fqbn", testutil.ArduinoMegaFQBN()}
		env.SetArgs(args)
		err := env.RootCmd.ExecuteContext(env.Ctx)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("compiles project that requires project library", t, func(env *testutil.IntegrationTestEnv) {
		err := env.InstallAvrPlatform(testutil.GlobalOpt{false})
		env.AddLib("Adafruit Pixie", testutil.GlobalOpt{false})
		pixieDir := testutil.PixieProjectDir()
		args := []string{"compile", pixieDir, "--fqbn", testutil.ArduinoMegaFQBN()}
		env.SetArgs(args)
		err = env.RootCmd.ExecuteContext(env.Ctx)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("errors when compiling project that requires project library", t, func(env *testutil.IntegrationTestEnv) {
		err := env.InstallAvrPlatform(testutil.GlobalOpt{false})
		pixieDir := testutil.PixieProjectDir()
		args := []string{"compile", pixieDir, "--fqbn", testutil.ArduinoMegaFQBN()}
		env.SetArgs(args)
		err = env.RootCmd.ExecuteContext(env.Ctx)
		assert.Error(env.T, err)
	})
}
