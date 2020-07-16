package commands_test

import (
	"os"
	"path"
	"testing"

	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestBuildCommandGlobal(t *testing.T) {
	testutil.RunIntegrationTest("builds projects", t, func(env *testutil.IntegrationTestEnv) {
		testutil.CleanBuilds()

		buildName := "pixie"
		projectDir := testutil.PixieProjectDir()
		buildDir := path.Join(projectDir, "build")
		boardURL := testutil.Esp8266BoardURL()
		platform := testutil.Esp8266Platform()
		fqbn := testutil.Esp8266WifiduinoFQBN()
		lib := "Adafruit Pixie"

		args := []string{"add", "board-urls", boardURL, "--global"}
		err := env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"add", "platform", platform, "--global"}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"add", "lib", lib, "--global"}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"add", "build", "--name", buildName, "--fqbn", fqbn, "--sketch", projectDir, "--global"}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"build", "--global"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.DirExists(env.T, buildDir)

		// build a single project (use same test to avoid reinstalling deps)
		testutil.CleanBuilds()
		_, err = os.Stat(buildDir)
		assert.True(env.T, os.IsNotExist(err))

		args = []string{"build", buildName, "--global"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.DirExists(env.T, buildDir)
	})
}

func TestBuildCommand(t *testing.T) {
	testutil.RunIntegrationTest("errors if project not initialized", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"build"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("builds projects", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		buildName := "pixie"
		projectDir := testutil.PixieProjectDir()
		buildDir := path.Join(projectDir, "build")
		boardURL := testutil.Esp8266BoardURL()
		platform := testutil.Esp8266Platform()
		fqbn := testutil.Esp8266WifiduinoFQBN()
		lib := "Adafruit Pixie"

		args := []string{"add", "board-urls", boardURL}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"add", "platform", platform}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"add", "lib", lib}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"add", "build", "--name", buildName, "--fqbn", fqbn, "--sketch", projectDir}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"build"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.DirExists(env.T, buildDir)

		// build a single project (use same test to avoid reinstalling deps)
		testutil.CleanBuilds()
		_, err = os.Stat(buildDir)
		assert.True(env.T, os.IsNotExist(err))

		args = []string{"build", buildName}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.DirExists(env.T, buildDir)
	})
}
