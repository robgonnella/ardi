package commands_test

import (
	"os"
	"testing"

	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/robgonnella/ardi/v2/util"
	"github.com/stretchr/testify/assert"
)

func TestCleanCommand(t *testing.T) {
	testutil.RunIntegrationTest("deletes project level .ardi directory and ardi.json file", t, func(env *testutil.IntegrationTestEnv) {
		err := util.InitDataDirectory("2222", paths.ArdiProjectDataDir, paths.ArdiProjectDataConfig)
		assert.NoError(env.T, err)
		err = util.InitArdiJSON()
		assert.NoError(env.T, err)
		assert.DirExists(env.T, ".ardi")
		assert.FileExists(env.T, "ardi.json")

		args := []string{"clean"}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		_, dirErr := os.Stat(".ardi")
		_, fileErr := os.Stat("ardi.json")

		assert.True(env.T, os.IsNotExist(dirErr))
		assert.True(env.T, os.IsNotExist(fileErr))
	})

	testutil.RunIntegrationTest("deletes global level .ardi directory and ardi.json file", t, func(env *testutil.IntegrationTestEnv) {
		err := util.InitDataDirectory("2222", paths.ArdiGlobalDataDir, paths.ArdiGlobalDataConfig)
		assert.NoError(env.T, err)

		assert.DirExists(env.T, paths.ArdiGlobalDataDir)
		assert.FileExists(env.T, paths.ArdiGlobalDataConfig)

		args := []string{"clean", "--global"}
		err = env.Execute(args)

		assert.NoError(env.T, err)

		_, dirErr := os.Stat(paths.ArdiGlobalDataDir)
		_, fileErr := os.Stat(paths.ArdiGlobalDataConfig)
		assert.True(env.T, os.IsNotExist(dirErr))
		assert.True(env.T, os.IsNotExist(fileErr))
	})
}
