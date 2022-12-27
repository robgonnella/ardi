package commands_test

import (
	"os"
	"testing"

	"github.com/robgonnella/ardi/v3/paths"
	"github.com/robgonnella/ardi/v3/testutil"
	"github.com/stretchr/testify/assert"
)

func TestProjectInitCommand(t *testing.T) {
	testutil.RunIntegrationTest("initializes a project directory", t, func(env *testutil.IntegrationTestEnv) {
		_, dataConfigErr := os.Stat(paths.ArduinoCliProjectConfig)
		_, buildConfigErr := os.Stat(paths.ArdiProjectConfig)
		assert.True(env.T, os.IsNotExist(dataConfigErr))
		assert.True(env.T, os.IsNotExist(buildConfigErr))

		args := []string{"project-init"}
		err := env.Execute(args)
		assert.NoError(env.T, err)

		_, dataConfigErr = os.Stat(paths.ArduinoCliProjectConfig)
		_, buildConfigErr = os.Stat(paths.ArdiProjectConfig)
		assert.NoError(env.T, dataConfigErr)
		assert.NoError(env.T, buildConfigErr)

	})
}
