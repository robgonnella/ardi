package commands_test

import (
	"testing"

	"github.com/robgonnella/ardi/v3/testutil"
	"github.com/robgonnella/ardi/v3/version"
	"github.com/stretchr/testify/assert"
)

func TestVersionCommand(t *testing.T) {
	testutil.RunIntegrationTest("prints ardi vesion", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"version"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), version.VERSION)
	})
}
