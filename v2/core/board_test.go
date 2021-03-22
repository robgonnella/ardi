package core_test

import (
	"errors"
	"testing"

	"github.com/arduino/arduino-cli/rpc/commands"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestBoardCore(t *testing.T) {
	testutil.RunUnitTest("prints fqbns", t, func(env *testutil.UnitTestEnv) {
		boards := testutil.GenerateCmdBoards(10)
		platform := testutil.GenerateCmdPlatform("test-platform", boards)
		platforms := []*commands.Platform{platform}

		env.Cli.EXPECT().GetPlatforms().Return(platforms, nil)
		env.ArdiCore.Board.FQBNS("")

		for _, b := range boards {
			assert.Contains(env.T, env.Stdout.String(), b.GetName())
			assert.Contains(env.T, env.Stdout.String(), b.GetFqbn())
		}
	})

	testutil.RunUnitTest("returns fqbn error", t, func(env *testutil.UnitTestEnv) {
		errString := "dummy error"
		dummyErr := errors.New(errString)

		env.Cli.EXPECT().GetPlatforms().Return(nil, dummyErr)
		err := env.ArdiCore.Board.FQBNS("")
		assert.Error(env.T, err)
		assert.EqualError(env.T, err, errString)
	})

	testutil.RunUnitTest("prints platforms", t, func(env *testutil.UnitTestEnv) {
		boards := testutil.GenerateCmdBoards(10)
		platform := testutil.GenerateCmdPlatform("test-platform", boards)
		platforms := []*commands.Platform{platform}

		env.Cli.EXPECT().GetPlatforms().Return(platforms, nil)
		env.ArdiCore.Board.Platforms("")

		for _, b := range boards {
			assert.Contains(env.T, env.Stdout.String(), b.GetName())
			assert.Contains(env.T, env.Stdout.String(), platform.GetID())
		}
	})

	testutil.RunUnitTest("returns platform error", t, func(env *testutil.UnitTestEnv) {
		errString := "dummy error"
		dummyErr := errors.New(errString)

		env.Cli.EXPECT().GetPlatforms().Return(nil, dummyErr)
		err := env.ArdiCore.Board.Platforms("")
		assert.Error(env.T, err)
		assert.EqualError(env.T, err, errString)
	})
}
