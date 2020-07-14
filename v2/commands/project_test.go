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
	testutil.RunIntegrationTest("with platform and library added", t, func(groupEnv *testutil.IntegrationTestEnv) {
		err := groupEnv.RunProjectInit()
		assert.NoError(groupEnv.T, err)

		platform := "arduino:avr"
		groupEnv.T.Logf("Adding platform: %s", platform)
		platArgs := []string{"project", "add", "platform", platform}
		err = groupEnv.Execute(platArgs)
		assert.NoError(groupEnv.T, err)

		lib := "Adafruit Pixie"
		groupEnv.T.Logf("Adding library: %s", lib)
		installedLib := "Adafruit_Pixie"
		libArgs := []string{"project", "add", "lib", lib}
		err = groupEnv.Execute(libArgs)
		assert.NoError(groupEnv.T, err)

		groupEnv.T.Run("lists platforms", func(st *testing.T) {
			groupEnv.ClearStdout()
			args := []string{"project", "list", "platforms"}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
			assert.Contains(st, groupEnv.Stdout.String(), platform)
		})

		groupEnv.T.Run("builds all projects", func(st *testing.T) {
			testutil.CleanBuilds()
			args := []string{"project", "add", "build", "--name", "pixie", "--fqbn", testutil.ArduinoMegaFQBN(), "--sketch", testutil.PixieProjectDir()}
			err := groupEnv.Execute(args)
			assert.NoError(st, err)

			args = []string{"project", "build"}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
		})

		groupEnv.T.Run("builds a single project", func(st *testing.T) {
			testutil.CleanBuilds()
			buildName := "pixie"
			args := []string{"project", "add", "build", "--name", buildName, "--fqbn", testutil.ArduinoMegaFQBN(), "--sketch", testutil.PixieProjectDir()}
			err := groupEnv.Execute(args)
			assert.NoError(st, err)

			args = []string{"project", "build", buildName}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
		})

		groupEnv.T.Run("removes lib and platform then reinstalls dependencies", func(st *testing.T) {
			groupEnv.ClearStdout()
			args := []string{"lib", "list"}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
			assert.Contains(st, groupEnv.Stdout.String(), installedLib)

			groupEnv.ClearStdout()
			args = []string{"platform", "list", "--installed"}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
			assert.Contains(st, groupEnv.Stdout.String(), platform)

			args = []string{"lib", "remove", installedLib}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)

			args = []string{"platform", "remove", platform}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)

			groupEnv.ClearStdout()
			args = []string{"lib", "list"}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
			assert.NotContains(st, groupEnv.Stdout.String(), lib)
			assert.NotContains(st, groupEnv.Stdout.String(), installedLib)

			groupEnv.ClearStdout()
			args = []string{"platform", "list", "--installed"}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
			assert.NotContains(st, groupEnv.Stdout.String(), platform)

			args = []string{"project", "install"}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)

			groupEnv.ClearStdout()
			args = []string{"lib", "list"}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
			assert.Contains(st, groupEnv.Stdout.String(), installedLib)

			groupEnv.ClearStdout()
			args = []string{"platform", "list", "--installed"}
			err = groupEnv.Execute(args)
			assert.NoError(st, err)
			assert.Contains(st, groupEnv.Stdout.String(), platform)
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
		name := "pixie"
		fqbn := "arduino:avr:mega"
		sketchPath := testutil.PixieProjectDir()
		args := []string{"project", "add", "build", "--name", name, "--fqbn", fqbn, "--sketch", sketchPath}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		args = []string{"project", "list", "builds"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		out := env.Stdout.String()
		assert.Contains(env.T, out, name)
		assert.Contains(env.T, out, sketchPath)
		assert.Contains(env.T, out, fqbn)
	})

	testutil.RunIntegrationTest("does not error if no builds added", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)
		args := []string{"project", "list", "builds"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunIntegrationTest("lists board urls", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		boardURL := "https://someboardurl.com"
		args := []string{"project", "add", "board-url", boardURL}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		env.ClearStdout()
		args = []string{"project", "list", "board-urls"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), boardURL)
	})

	testutil.RunIntegrationTest("does not error if no board urls added", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)
		args := []string{"project", "list", "board-urls"}
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

		args = []string{"project", "list", "board-urls"}
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

	testutil.RunIntegrationTest("adds board url", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		boardURL := "https://someboardurl.com"
		args := []string{"project", "add", "board-url", boardURL}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		env.ClearStdout()
		args = []string{"project", "list", "board-urls"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), boardURL)
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

		args = []string{"project", "add", "board-url", "https://someboardurl.com"}
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

	testutil.RunIntegrationTest("removes board url", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)
		boardURL := "https://someboardurl.com"
		args := []string{"project", "add", "board-url", boardURL}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		env.ClearStdout()
		args = []string{"project", "list", "board-urls"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), boardURL)

		args = []string{"project", "remove", "board-url", boardURL}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		env.ClearStdout()
		args = []string{"project", "list", "board-urls"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.NotContains(env.T, env.Stdout.String(), boardURL)
	})

	testutil.RunIntegrationTest("does not error if board url not added", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)
		args := []string{"project", "remove", "board-url", "noop"}
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

		args = []string{"project", "remove", "board-url", "https://someboardurl.com"}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})
}
