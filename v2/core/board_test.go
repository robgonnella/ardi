package core_test

import (
	"errors"
	"testing"

	"github.com/arduino/arduino-cli/rpc/commands"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestBoardCore(t *testing.T) {
	testutil.RunTest("prints fqbns", t, func(st *testing.T, env testutil.TestEnv) {
		defer env.Ctrl.Finish()
		boards := testutil.GenerateCmdBoards(10)
		platform := testutil.GenerateCmdPlatform("test-platform", boards)
		platforms := []*commands.Platform{platform}

		env.Client.EXPECT().GetPlatforms().Return(platforms, nil)
		env.ArdiCore.Board.FQBNS("")

		for _, b := range boards {
			assert.Contains(st, env.Stdout.String(), b.GetName())
			assert.Contains(st, env.Stdout.String(), b.GetFqbn())
		}
	})

	testutil.RunTest("returns fqbn error", t, func(st *testing.T, env testutil.TestEnv) {
		defer env.Ctrl.Finish()
		errString := "dummy error"
		dummyErr := errors.New(errString)

		env.Client.EXPECT().GetPlatforms().Return(nil, dummyErr)
		err := env.ArdiCore.Board.FQBNS("")
		assert.Error(st, err)
		assert.EqualError(st, err, errString)
	})

	testutil.RunTest("prints platforms", t, func(st *testing.T, env testutil.TestEnv) {
		defer env.Ctrl.Finish()
		boards := testutil.GenerateCmdBoards(10)
		platform := testutil.GenerateCmdPlatform("test-platform", boards)
		platforms := []*commands.Platform{platform}

		env.Client.EXPECT().GetPlatforms().Return(platforms, nil)
		env.ArdiCore.Board.Platforms("")

		for _, b := range boards {
			assert.Contains(st, env.Stdout.String(), b.GetName())
			assert.Contains(st, env.Stdout.String(), platform.GetID())
		}
	})

	testutil.RunTest("returns platform error", t, func(st *testing.T, env testutil.TestEnv) {
		defer env.Ctrl.Finish()
		errString := "dummy error"
		dummyErr := errors.New(errString)

		env.Client.EXPECT().GetPlatforms().Return(nil, dummyErr)
		err := env.ArdiCore.Board.Platforms("")
		assert.Error(st, err)
		assert.EqualError(st, err, errString)
	})
}
