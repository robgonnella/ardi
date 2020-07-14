package commands_test

import (
	"testing"

	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestPlatform(t *testing.T) {
	testutil.RunIntegrationTest("adding, listing, and removing platform globally", t, func(env *testutil.IntegrationTestEnv) {
		platform := "arduino:avr"
		args := []string{"platform", "add", platform, "--global"}
		err := env.Execute(args)
		assert.NoError(env.T, err)

		env.ClearStdout()
		args = []string{"platform", "list", "--installed", "--global"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), platform)

		args = []string{"platform", "remove", platform, "--global"}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		env.ClearStdout()
		args = []string{"platform", "list", "--installed", "--global"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.NotContains(env.T, env.Stdout.String(), platform)
	})

	testutil.RunIntegrationTest("adding, listing, and removing platform in project", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		platform := "arduino:avr"
		args := []string{"platform", "add", platform}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		env.ClearStdout()
		args = []string{"platform", "list", "--installed"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), platform)

		args = []string{"platform", "remove", platform}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		env.ClearStdout()
		args = []string{"platform", "list", "--installed"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.NotContains(env.T, env.Stdout.String(), platform)
	})
}

func TestPlatformAddCommand(t *testing.T) {
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
	testutil.RunIntegrationTest("errors when removing an invalid platform globally", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"platform", "remove", "noop", "--global"}
		err := env.Execute(args)
		assert.Error(env.T, err)
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
