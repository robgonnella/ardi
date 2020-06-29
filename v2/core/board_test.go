package core_test

import (
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
}
