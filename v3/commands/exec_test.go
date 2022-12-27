package commands_test

import (
	"testing"

	"github.com/robgonnella/ardi/v3/testutil"
	"github.com/stretchr/testify/assert"
)

func TestExecCommand(t *testing.T) {
	testutil.RunIntegrationTest("errors if project not initialized", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"exec", "arduino-cli", "--", "compile", "--fqbn", testutil.ArduinoMegaFQBN(), testutil.PixieProjectDir()}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunIntegrationTest("compiles project sketch that requires library", t, func(env *testutil.IntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		err = env.Execute([]string{"add", "lib", "Adafruit Pixie"})
		assert.NoError(env.T, err)

		err = env.Execute([]string{"add", "platform", "arduino:avr"})
		assert.NoError(env.T, err)

		args := []string{"exec", "--", "arduino-cli", "compile", "--fqbn", testutil.ArduinoMegaFQBN(), testutil.PixieProjectDir()}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})
}
