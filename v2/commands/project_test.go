package commands_test

import (
	"os"
	"testing"

	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestProjectInitCommand(t *testing.T) {
	testutil.RunIntegrationTest("initializes a project directory", t, func(env *testutil.IntegrationTestEnv) {
		_, dataConfigErr := os.Stat(paths.ArdiProjectDataConfig)
		_, buildConfigErr := os.Stat(paths.ArdiProjectBuildConfig)
		assert.True(env.T, os.IsNotExist(dataConfigErr))
		assert.True(env.T, os.IsNotExist(buildConfigErr))

		args := []string{"project", "init"}
		err := env.Execute(args)
		assert.NoError(env.T, err)

		_, dataConfigErr = os.Stat(paths.ArdiProjectDataConfig)
		_, buildConfigErr = os.Stat(paths.ArdiProjectBuildConfig)
		assert.NoError(env.T, dataConfigErr)
		assert.NoError(env.T, buildConfigErr)

	})
}

func TestProjectRequiringPlatform(t *testing.T) {
	testutil.RunIntegrationTest("project platform", t, func(groupEnv *testutil.IntegrationTestEnv) {
		err := groupEnv.RunProjectInit()
		assert.NoError(groupEnv.T, err)

		platform := "arduino:avr"
		groupEnv.T.Logf("Adding platform: %s", platform)
		err = groupEnv.AddPlatform(platform, testutil.GlobalOpt{false})
		assert.NoError(groupEnv.T, err)

		groupEnv.T.Run("lists platforms", func(st *testing.T) {
			args := []string{"project", "list", "platforms"}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
			assert.Contains(st, groupEnv.Stdout.String(), platform)
		})

		groupEnv.T.Run("builds all projects", func(st *testing.T) {
			testutil.CleanBuilds()
			args := []string{"project", "add", "build", "--name", "blink", "--fqbn", testutil.ArduinoMegaFQBN(), "--sketch", testutil.BlinkProjectDir()}
			err := groupEnv.Execute(args)
			assert.NoError(st, err)

			args = []string{"project", "build"}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
		})

		groupEnv.T.Run("builds a single project", func(st *testing.T) {
			testutil.CleanBuilds()
			buildName := "blink"
			args := []string{"project", "add", "build", "--name", buildName, "--fqbn", testutil.ArduinoMegaFQBN(), "--sketch", testutil.BlinkProjectDir()}
			err := groupEnv.Execute(args)
			assert.NoError(st, err)

			args = []string{"project", "build", buildName}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
		})

		groupEnv.T.Run("removes platform", func(st *testing.T) {
			args := []string{"project", "remove", "platform", platform}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
		})
	})
}

func TestProjectListCommand(t *testing.T) {
	testutil.RunIntegrationTest("does not error if no platforms installed", t, func(env *testutil.IntegrationTestEnv) {
		env.RunProjectInit()
		args := []string{"project", "list", "platforms"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("lists libraries", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)
		lib := "Adafruit Pixie"
		env.AddLib(lib, testutil.GlobalOpt{false})
		args := []string{"project", "list", "libs"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), lib)
	})

	testutil.RunIntegrationTest("does not error if no libraries installed", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)
		args := []string{"project", "list", "libs"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("lists builds", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)
		name := "some_build_name"
		platform := "some-platform"
		fqbn := "some-fqbn"
		args := []string{"project", "add", "build", "--name", name, "--platform", platform, "--fqbn", fqbn, "--sketch", "."}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		args = []string{"project", "list", "builds"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		out := env.Stdout.String()
		assert.Contains(env.T, out, name)
		assert.Contains(env.T, out, platform)
		assert.Contains(env.T, out, fqbn)
	})

	testutil.RunIntegrationTest("does not error if no builds added", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)
		args := []string{"project", "list", "builds"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("errors if project not initialized", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"project", "list", "platforms"}
		err := env.Execute(args)
		assert.Error(env.T, err)

		args = []string{"project", "list", "libs"}

		err = env.Execute(args)
		assert.Error(env.T, err)

		args = []string{"project", "list", "builds"}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})
}

func TestProjectAddCommand(t *testing.T) {
	testutil.RunIntegrationTest("errors on invalid platform", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)
		args := []string{"project", "add", "platform", "noop"}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("adds library", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)
		args := []string{"project", "add", "library", "Adafruit Pixie"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("errors on invalid library", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)
		args := []string{"project", "add", "lib", "noop"}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("adds build", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)
		args := []string{"project", "add", "build", "--name", "pixie", "--fqbn", "somefqbn", "--sketch", "."}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("errors if name missing", t, func(env *testutil.IntegrationTestEnv) {
		env.RunProjectInit()
		args := []string{"project", "add", "build", "--fqbn", "somefqbn", "--sketch", "."}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("errors if fqbn missing", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)
		args := []string{"project", "add", "build", "--name", "pixie", "--sketch", "."}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("errors if sketch missing", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)
		args := []string{"project", "add", "build", "--name", "pixie", "--fqbn", "somefqbn"}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("errors if project not initialized", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"project", "add", "platform", "arduino:avr"}
		err := env.Execute(args)
		assert.Error(env.T, err)

		args = []string{"project", "add", "build", "--name", "pixie", "--sketch", "."}
		err = env.Execute(args)
		assert.Error(env.T, err)

		args = []string{"project", "add", "build", "--name", "pixie", "--sketch", "."}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})
}

func TestProjectRemoveCommand(t *testing.T) {
	testutil.RunIntegrationTest("errors if platform not installed", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)
		args := []string{"project", "remove", "platform", "arduino:avr"}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("removes library", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)
		lib := "Adafruit Pixie"
		env.AddLib(lib, testutil.GlobalOpt{false})
		args := []string{"project", "remove", "lib", lib}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("does not error if library not added", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)
		args := []string{"project", "remove", "lib", "noop"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("removes build", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)
		buildName := "pixie"
		args := []string{"project", "add", "build", "--name", buildName, "--fqbn", "somefqbn", "--sketch", "."}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"project", "remove", "build", buildName}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("does not error if build not added", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)
		args := []string{"project", "remove", "build", "noop"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("errors if project not initialized", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"project", "remove", "platform", "arduino:sam"}
		err := env.Execute(args)
		assert.Error(env.T, err)

		args = []string{"project", "remove", "lib", "Adafruit Pixie"}
		err = env.Execute(args)
		assert.Error(env.T, err)

		args = []string{"project", "remove", "build", "pixie"}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})
}

func TestProjectInstallCommand(t *testing.T) {
	testutil.RunIntegrationTest("installs dependencies", t, func(env *testutil.IntegrationTestEnv) {
		env.RunProjectInit()
		lib := "Adafruit Pixie"
		env.AddLib(lib, testutil.GlobalOpt{false})

		args := []string{"lib", "remove", lib}
		err := env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"project", "install"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})
}
