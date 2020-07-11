package commands_test

import (
	"testing"

	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestLibAddCommand(t *testing.T) {
	testutil.RunIntegrationTest("globally adds a valid library", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"lib", "add", "Adafruit Pixie", "--global"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("errors when adding invalid library globally", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"lib", "add", "Noop", "--global"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("errors when adding a library to uninitialized project", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"lib", "add", "Adafruit Pixie"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("errors when adding an invalid library to ardi project", t, func(env *testutil.IntegrationTestEnv) {
		env.RunProjectInit()
		args := []string{"lib", "add", "Noop"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("adds valid library to ardi project", t, func(env *testutil.IntegrationTestEnv) {
		env.RunProjectInit()
		args := []string{"lib", "add", "Adafruit Pixie"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})
}

func TestLibRemoveCommand(t *testing.T) {
	testutil.RunIntegrationTest("globally removes a valid library", t, func(env *testutil.IntegrationTestEnv) {
		lib := "Adafruit Pixie"
		env.AddLib(lib, testutil.GlobalOpt{true})
		args := []string{"lib", "remove", lib, "--global"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("does not error when removing invalid library globally", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"lib", "remove", "Noop", "--global"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("errors when removing a library from uninitialized project", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"lib", "remove", "Adafruit Pixie"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("does not error when removing an invalid library from ardi project", t, func(env *testutil.IntegrationTestEnv) {
		env.RunProjectInit()
		args := []string{"lib", "remove", "Noop"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("removes valid library from ardi project", t, func(env *testutil.IntegrationTestEnv) {
		env.RunProjectInit()
		lib := "Adafruit Pixie"
		env.AddLib(lib, testutil.GlobalOpt{false})
		args := []string{"lib", "remove", "Adafruit Pixie"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})
}

func TestLibSearchCommand(t *testing.T) {
	testutil.RunIntegrationTest("finds a valid library", t, func(env *testutil.IntegrationTestEnv) {
		env.RunProjectInit()
		searchLib := "Adafruit Pixie"
		args := []string{"lib", "search", searchLib}
		err := env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), searchLib)
	})

	testutil.RunIntegrationTest("errors on invalid library", t, func(env *testutil.IntegrationTestEnv) {
		env.RunProjectInit()
		searchLib := "noop"
		args := []string{"lib", "search", searchLib}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})
}

func TestLibListCommand(t *testing.T) {
	testutil.RunIntegrationTest("lists globally installed library", t, func(env *testutil.IntegrationTestEnv) {
		lib := "Adafruit Pixie"
		env.AddLib(lib, testutil.GlobalOpt{true})
		args := []string{"lib", "list", "--global"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), lib)
	})

	testutil.RunIntegrationTest("does not error if no global libs found", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"lib", "list", "--global"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("lists project level installed library", t, func(env *testutil.IntegrationTestEnv) {
		env.RunProjectInit()
		lib := "Adafruit Pixie"
		env.AddLib(lib, testutil.GlobalOpt{false})
		args := []string{"lib", "list"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), lib)
	})

	testutil.RunIntegrationTest("does not error if no project libs found", t, func(env *testutil.IntegrationTestEnv) {
		env.RunProjectInit()
		args := []string{"lib", "list"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})
}
