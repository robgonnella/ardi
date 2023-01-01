package core_test

import (
	"testing"

	"github.com/robgonnella/ardi/v3/paths"
	"github.com/robgonnella/ardi/v3/testutil"
	"github.com/robgonnella/ardi/v3/util"
	"github.com/stretchr/testify/assert"
)

func TestArduinoCliConfig(t *testing.T) {
	testutil.RunUnitTest("adds and removes board urls", t, func(env *testutil.UnitTestEnv) {
		util.InitProjectDirectory()

		boardURL1 := "https://somefakeboardurl.com"
		boardURL2 := "https://anotherfakeboardurl.com"

		err := env.ArdiCore.CliConfig.AddBoardURL(boardURL1)
		assert.NoError(env.T, err)

		err = env.ArdiCore.CliConfig.AddBoardURL(boardURL2)
		assert.NoError(env.T, err)

		settings, err := util.ReadArduinoCliSettings(paths.ArduinoCliProjectConfig)
		assert.NoError(env.T, err)

		assert.Contains(env.T, settings.BoardManager.AdditionalUrls, boardURL1)
		assert.Contains(env.T, settings.BoardManager.AdditionalUrls, boardURL2)

		err = env.ArdiCore.CliConfig.RemoveBoardURL(boardURL1)
		assert.NoError(env.T, err)

		settings, err = util.ReadArduinoCliSettings(paths.ArduinoCliProjectConfig)
		assert.NoError(env.T, err)
		assert.NotContains(env.T, settings.BoardManager.AdditionalUrls, boardURL1)
		assert.Contains(env.T, settings.BoardManager.AdditionalUrls, boardURL2)
	})

	testutil.RunUnitTest("doesnt error adding same url twice", t, func(env *testutil.UnitTestEnv) {
		util.InitProjectDirectory()

		boardURL1 := "https://somefakeboardurl.com"

		err := env.ArdiCore.CliConfig.AddBoardURL(boardURL1)
		assert.NoError(env.T, err)

		err = env.ArdiCore.CliConfig.AddBoardURL(boardURL1)
		assert.NoError(env.T, err)

		settings, err := util.ReadArduinoCliSettings(paths.ArduinoCliProjectConfig)
		assert.NoError(env.T, err)

		assert.Contains(env.T, settings.BoardManager.AdditionalUrls, boardURL1)
		assert.Equal(env.T, len(settings.BoardManager.AdditionalUrls), 1)
	})

	testutil.RunUnitTest("doesnt error removing non-exiting url", t, func(env *testutil.UnitTestEnv) {
		util.InitProjectDirectory()

		boardURL1 := "https://somefakeboardurl.com"

		err := env.ArdiCore.CliConfig.RemoveBoardURL(boardURL1)
		assert.NoError(env.T, err)
	})
}
