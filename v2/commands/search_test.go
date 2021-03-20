package commands_test

import (
	"testing"

	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestSearchLibCommandGlobal(t *testing.T) {
	testutil.RunIntegrationTest("searches a valid library", t, func(env *testutil.IntegrationTestEnv) {
		searchLib := "Adafruit Pixie"
		args := []string{"search", "libraries", searchLib, "--global"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), searchLib)
	})

	testutil.RunIntegrationTest("errors on invalid library", t, func(env *testutil.IntegrationTestEnv) {
		searchLib := "noop"
		args := []string{"search", "libraries", searchLib, "--global"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("does not error if search arg not provided", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"search", "libraries", "--global"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})
}

func TestSearchLibCommand(t *testing.T) {
	testutil.RunIntegrationTest("searches a valid library", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		searchLib := "Adafruit Pixie"
		args := []string{"search", "libs", searchLib}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), searchLib)
	})

	testutil.RunIntegrationTest("errors on invalid library", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		searchLib := "noop"
		args := []string{"search", "libs", searchLib}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("does not error if search arg not provided", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		args := []string{"search", "libs"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("errors if project not initialized", t, func(env *testutil.IntegrationTestEnv) {
		searchLib := "Adafruit Pixie"
		args := []string{"search", "libs", searchLib}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})
}

func TestSearchPlatformCommand(t *testing.T) {
	testutil.RunIntegrationTest("searches globally available platforms", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"search", "platforms", "--global"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), "arduino:avr")
	})

	testutil.RunIntegrationTest("searches platforms available to project", t, func(env *testutil.IntegrationTestEnv) {
		env.RunProjectInit()
		args := []string{"search", "platforms"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), "arduino:avr")
	})
}
