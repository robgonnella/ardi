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
		_, cliConfErr := os.Stat(paths.ArduinoCliProjectConfig)

		assert.True(env.T, os.IsNotExist(dirErr))
		assert.True(env.T, os.IsNotExist(cliConfErr))
		assert.FileExists(env.T, paths.ArdiProjectConfig)
	})
}
