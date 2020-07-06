package commands_test

import (
	"testing"

	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCompileCommand(t *testing.T) {
	testutil.RunIntegrationTest("compiles project directory", t, func(env *testutil.IntegrationTestEnv) {
		err := env.InstallAvrPlatform()
		assert.NoError(env.T, err)
		projectDir := testutil.BlinkProjectDir()
		args := []string{"compile", projectDir, "--fqbn", testutil.ArduinoMegaFQBN()}
		env.SetArgs(args)
		err = env.RootCmd.ExecuteContext(env.Ctx)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("errors if fqbn is missing", t, func(env *testutil.IntegrationTestEnv) {
		err := env.InstallAvrPlatform()
		assert.NoError(env.T, err)
		projectDir := testutil.BlinkProjectDir()
		args := []string{"compile", projectDir}
		env.SetArgs(args)
		err = env.RootCmd.ExecuteContext(env.Ctx)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("errors if platform not installed", t, func(env *testutil.IntegrationTestEnv) {
		projectDir := testutil.BlinkProjectDir()
		args := []string{"compile", projectDir, "--fqbn", testutil.ArduinoMegaFQBN()}
		env.SetArgs(args)
		err := env.RootCmd.ExecuteContext(env.Ctx)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("errors if not a valid project directory", t, func(env *testutil.IntegrationTestEnv) {
		err := env.InstallAvrPlatform()
		assert.NoError(env.T, err)
		args := []string{"compile", ".", "--fqbn", testutil.ArduinoMegaFQBN()}
		env.SetArgs(args)
		err = env.RootCmd.ExecuteContext(env.Ctx)
		assert.Error(env.T, err)
	})
}
