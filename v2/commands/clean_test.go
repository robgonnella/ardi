package commands_test

import (
	"os"
	"testing"

	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCleanCommand(t *testing.T) {
	testutil.RunIntegrationTest("deletes project level .ardi directory and ardi.json file", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		assert.DirExists(env.T, paths.ArdiProjectDataDir)
		assert.FileExists(env.T, paths.ArdiProjectConfig)
		assert.FileExists(env.T, paths.ArduinoCliProjectConfig)

		args := []string{"clean"}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		_, dirErr := os.Stat(paths.ArdiProjectDataDir)
		_, file1Err := os.Stat(paths.ArdiProjectConfig)
		_, file2Err := os.Stat(paths.ArduinoCliProjectConfig)

		assert.True(env.T, os.IsNotExist(dirErr))
		assert.True(env.T, os.IsNotExist(file1Err))
		assert.True(env.T, os.IsNotExist(file2Err))
	})

	testutil.RunIntegrationTest("deletes global level .ardi directory and ardi.json file", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"version"}
		err := env.Execute(args)
		assert.NoError(env.T, err)

		assert.DirExists(env.T, paths.ArdiGlobalDataDir)
		assert.FileExists(env.T, paths.ArdiGlobalConfig)
		assert.FileExists(env.T, paths.ArduinoCliGlobalConfig)

		args = []string{"clean", "--global"}
		err = env.Execute(args)

		assert.NoError(env.T, err)

		_, dirErr := os.Stat(paths.ArdiGlobalDataDir)
		_, file1Err := os.Stat(paths.ArdiGlobalConfig)
		_, file2Err := os.Stat(paths.ArduinoCliGlobalConfig)

		assert.True(env.T, os.IsNotExist(dirErr))
		assert.True(env.T, os.IsNotExist(file1Err))
		assert.True(env.T, os.IsNotExist(file2Err))
	})
}
