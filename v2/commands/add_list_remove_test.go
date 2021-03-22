package commands_test

import (
	"testing"

	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAddListRemovePlatform(t *testing.T) {
	testutil.RunIntegrationTest("errors if project not initialized", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"add", "platforms", "arduino:avr"}
		err := env.Execute(args)
		assert.Error(env.T, err)

		args = []string{"list", "platforms"}

		err = env.Execute(args)
		assert.Error(env.T, err)

		args = []string{"remove", "platforms", "arduino:avr"}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("adding, listing, and removing platform", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		platform := "arduino:avr"
		args := []string{"add", "platform", platform}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		env.ClearStdout()
		args = []string{"list", "platforms"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), platform)

		env.ClearStdout()
		args = []string{"list", "board-fqbns"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), "arduino:avr:mega")

		env.ClearStdout()
		args = []string{"list", "board-platforms"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), "arduino:avr")

		args = []string{"remove", "platform", platform}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		env.ClearStdout()
		args = []string{"list", "platforms"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.NotContains(env.T, env.Stdout.String(), platform)
	})

	testutil.RunIntegrationTest("does not error listing if no platforms installed", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		args := []string{"list", "platforms"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("errors when adding invalid platform", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		args := []string{"add", "platform", "noop"}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("errors removing when platform not installed", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		args := []string{"remove", "platform", "arduino:avr"}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})
}

func TestAddListRemoveLibrary(t *testing.T) {
	testutil.RunIntegrationTest("errors if project not initialized", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"add", "libraries", "Adafruit Pixie"}
		err := env.Execute(args)
		assert.Error(env.T, err)

		args = []string{"list", "libraries"}

		err = env.Execute(args)
		assert.Error(env.T, err)

		args = []string{"remove", "libraries", "Adafruit_Pixie"}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("adding, listing, and removing library", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		library := "Adafruit Pixie"
		installedLib := "Adafruit_Pixie"

		args := []string{"add", "library", library}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		env.ClearStdout()
		args = []string{"list", "libraries"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), installedLib)

		args = []string{"remove", "library", installedLib}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		env.ClearStdout()
		args = []string{"list", "libraries"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.NotContains(env.T, env.Stdout.String(), installedLib)
	})

	testutil.RunIntegrationTest("does not error listing if no libraries installed", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		args := []string{"list", "libs"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("errors when adding invalid library", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		args := []string{"add", "lib", "noop"}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("does not error removing if library not installed", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		args := []string{"remove", "lib", "Adafruit_Pixie"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})
}

func TestAddListRemoveBoardURL(t *testing.T) {
	testutil.RunIntegrationTest("errors if project not initialized", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"add", "board-urls", "https://someboardurl.com"}
		err := env.Execute(args)
		assert.Error(env.T, err)

		args = []string{"list", "board-urls"}

		err = env.Execute(args)
		assert.Error(env.T, err)

		args = []string{"remove", "board-urls", "https://someboardurl.com"}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("adding, listing, and removing board urls", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		url := "https://someboardurl.com"

		args := []string{"add", "board-urls", url}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		env.ClearStdout()
		args = []string{"list", "board-urls"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), url)

		args = []string{"remove", "board-urls", url}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		env.ClearStdout()
		args = []string{"list", "board-urls"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.NotContains(env.T, env.Stdout.String(), url)
	})

	testutil.RunIntegrationTest("does not error listing if no board urls added", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		args := []string{"list", "board-urls"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("does not error removing if board url not added", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		args := []string{"remove", "board-url", "https://someboardurl.com"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})
}

func TestAddListRemoveBuild(t *testing.T) {
	testutil.RunIntegrationTest("errors if not a valid sketch path", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"add", "build", "--name", "somename", "--fqbn", "somefqbn", "--sketch", "noop"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("errors if project not initialized", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"add", "build", "--name", "somename", "--fqbn", "somefqbn", "--sketch", testutil.PixieProjectDir()}
		err := env.Execute(args)
		assert.Error(env.T, err)

		args = []string{"list", "builds"}

		err = env.Execute(args)
		assert.Error(env.T, err)

		args = []string{"remove", "builds", "somename"}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("adding, listing, and removing builds", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		name := "pixie"
		fqbn := "somefqbn"
		sketch := testutil.PixieProjectDir()

		args := []string{"add", "build", "--name", name, "--fqbn", fqbn, "--sketch", sketch}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		env.ClearStdout()
		args = []string{"list", "builds"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), name)
		assert.Contains(env.T, env.Stdout.String(), fqbn)
		assert.Contains(env.T, env.Stdout.String(), sketch)

		args = []string{"remove", "builds", name}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		env.ClearStdout()
		args = []string{"list", "builds"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.NotContains(env.T, env.Stdout.String(), name)
		assert.NotContains(env.T, env.Stdout.String(), fqbn)
		assert.NotContains(env.T, env.Stdout.String(), sketch)
	})

	testutil.RunIntegrationTest("errors adding if name missing", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		args := []string{"add", "build", "--fqbn", "somefqbn", "--sketch", "."}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("errors adding if fqbn missing", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		args := []string{"add", "build", "--name", "pixie", "--sketch", "."}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("errors adding if sketch missing", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		args := []string{"add", "build", "--name", "pixie", "--fqbn", "somefqbn"}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("does not error listing if no builds added", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		args := []string{"list", "builds"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("does not error removing if build not added", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		args := []string{"remove", "build", "noop"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})
}
