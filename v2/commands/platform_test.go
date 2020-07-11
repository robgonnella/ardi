package commands_test

import (
	"testing"

	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestPlatformAddCommand(t *testing.T) {
	testutil.RunIntegrationTest("adds a valid platform globally", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"platform", "add", "arduino:avr", "--global"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("errors when adding an invalid platform globally", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"platform", "add", "noop", "--global"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("adds a valid platform to project", t, func(env *testutil.IntegrationTestEnv) {
		env.RunProjectInit()
		args := []string{"platform", "add", "emoro:avr"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("errors when adding an invalid project platform", t, func(env *testutil.IntegrationTestEnv) {
		env.RunProjectInit()
		args := []string{"platform", "add", "noop"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})
}

func TestPlatformRemoveCommand(t *testing.T) {
	testutil.RunIntegrationTest("removes a valid platform globally", t, func(env *testutil.IntegrationTestEnv) {
		platform := "arduino:sam"
		env.AddPlatform(platform, testutil.GlobalOpt{true})
		args := []string{"platform", "remove", platform, "--global"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("errors when removing an invalid platform globally", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"platform", "remove", "noop", "--global"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("removes a valid platform from project", t, func(env *testutil.IntegrationTestEnv) {
		env.RunProjectInit()
		platform := "arduino:megaavr"
		env.AddPlatform(platform, testutil.GlobalOpt{false})
		args := []string{"platform", "remove", platform}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("errors when removing an invalid project platform", t, func(env *testutil.IntegrationTestEnv) {
		env.RunProjectInit()
		args := []string{"platform", "remove", "noop"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})
}

func TestPlatformListCommand(t *testing.T) {
	testutil.RunIntegrationTest("lists globally available platforms", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"platform", "list", "--global"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), "arduino:avr")
	})

	testutil.RunIntegrationTest("lists available platforms for project", t, func(env *testutil.IntegrationTestEnv) {
		env.RunProjectInit()
		args := []string{"platform", "list"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), "arduino:avr")
	})
}