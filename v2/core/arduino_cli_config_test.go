package core_test

import (
	"testing"

	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/robgonnella/ardi/v2/util"
	"github.com/stretchr/testify/assert"
)

func TestArduinoCliConfig(t *testing.T) {
	testutil.RunUnitTest("adds and removes board urls", t, func(env *testutil.UnitTestEnv) {
		util.InitProjectDirectory()

		boardURL := "https://somefakeboardurl.com"
		err := env.ArdiCore.CliConfig.AddBoardURL(boardURL)
		assert.NoError(env.T, err)

		settings, err := util.ReadArduinoCliSettings(paths.ArduinoCliProjectConfig)
		assert.NoError(env.T, err)

		assert.Contains(env.T, settings.BoardManager.AdditionalUrls, boardURL)

		err = env.ArdiCore.CliConfig.RemoveBoardURL(boardURL)
		assert.NoError(env.T, err)

		settings, err = util.ReadArduinoCliSettings(paths.ArduinoCliProjectConfig)
		assert.NoError(env.T, err)
		assert.NotContains(env.T, settings.BoardManager.AdditionalUrls, boardURL)
	})
}
