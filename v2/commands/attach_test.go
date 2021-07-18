package commands_test

import (
	"testing"

	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

// TODO: Find a way to test this
func TestAttachCommand(t *testing.T) {
	testutil.RunIntegrationTest("placeholder", t, func(env *testutil.IntegrationTestEnv) {
		// this cannot be tested at the moment as we don't have a way to inject a
		// mock serial port from here.
		assert.Equal(env.T, true, true)
	})
}
