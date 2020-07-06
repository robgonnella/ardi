package commands_test

import (
	"testing"

	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func runProjectInit(env *testutil.IntegrationTestEnv) {
	projectInitArgs := []string{"project", "init"}
	env.SetArgs(projectInitArgs)
	env.RootCmd.ExecuteContext(env.Ctx)
}

func TestLibAddCommand(t *testing.T) {
	testutil.RunIntegrationTest("globally adds a valid library", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"lib", "add", "Adafruit Pixie", "--global"}
		env.SetArgs(args)
		err := env.RootCmd.ExecuteContext(env.Ctx)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("errors when adding invalid library globally", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"lib", "add", "Noop", "--global"}
		env.SetArgs(args)
		err := env.RootCmd.ExecuteContext(env.Ctx)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("errors when adding a library to uninitialized project", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"lib", "add", "Adafruit Pixie"}
		env.SetArgs(args)
		err := env.RootCmd.ExecuteContext(env.Ctx)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("errors when adding an invalid library to ardi project", t, func(env *testutil.IntegrationTestEnv) {
		runProjectInit(env)
		args := []string{"lib", "add", "Noop"}
		env.SetArgs(args)
		err := env.RootCmd.ExecuteContext(env.Ctx)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("adds valid library to ardi project", t, func(env *testutil.IntegrationTestEnv) {
		runProjectInit(env)
		args := []string{"lib", "add", "Adafruit Pixie"}
		env.SetArgs(args)
		err := env.RootCmd.ExecuteContext(env.Ctx)
		assert.NoError(env.T, err)
	})
}
