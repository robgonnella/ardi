package commands_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/robgonnella/ardi/v2/types"
	"github.com/robgonnella/ardi/v2/util"
	"github.com/stretchr/testify/assert"
)

func TestInstallCommand(t *testing.T) {
	testutil.RunIntegrationTest("errors if project not initialized", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"install"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("removes lib and platform then reinstalls dependencies", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		platform := "arduino:avr"
		platArgs := []string{"add", "platform", platform}
		err = env.Execute(platArgs)
		assert.NoError(env.T, err)

		lib := "Adafruit Pixie"
		installedLib := "Adafruit Pixie"
		libArgs := []string{"add", "lib", lib}
		err = env.Execute(libArgs)
		assert.NoError(env.T, err)

		env.ClearStdout()
		args := []string{"list", "libs"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), installedLib)

		env.ClearStdout()
		args = []string{"list", "platforms"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), platform)

		// remove data directory
		os.RemoveAll(paths.ArdiProjectDataDir)

		args = []string{"install"}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		env.ClearStdout()
		args = []string{"list", "libraries"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), installedLib)

		env.ClearStdout()
		args = []string{"list", "platforms"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), platform)
	})

	testutil.RunIntegrationTest("returns platform install error", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		platform := "noop"

		var ardiConfig types.ArdiConfig
		byteData, _ := ioutil.ReadFile(paths.ArdiProjectConfig)
		json.Unmarshal(byteData, &ardiConfig)
		ardiConfig.Platforms[platform] = "0.1.0"

		util.WriteAllSettings(&ardiConfig, util.GenArduinoCliSettings(paths.ArdiProjectDataDir))

		// remove data directory
		os.RemoveAll(paths.ArdiProjectDataDir)

		args := []string{"install"}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("returns library install error", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		library := "noop"

		var ardiConfig types.ArdiConfig
		byteData, _ := ioutil.ReadFile(paths.ArdiProjectConfig)
		json.Unmarshal(byteData, &ardiConfig)
		ardiConfig.Libraries[library] = "1.3.5"

		util.WriteAllSettings(&ardiConfig, util.GenArduinoCliSettings(paths.ArdiProjectDataDir))

		// remove data directory
		os.RemoveAll(paths.ArdiProjectDataDir)

		args := []string{"install"}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})
}
