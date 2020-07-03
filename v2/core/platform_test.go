package core_test

import (
	"testing"

	"github.com/arduino/arduino-cli/rpc/commands"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

// @todo: check that list is actually sorted
func TestPlatformCore(t *testing.T) {
	testutil.RunTest("prints sorted list of all installed platforms to stdout", t, func(st *testing.T, env testutil.TestEnv) {
		defer env.Ctrl.Finish()
		platform1 := commands.Platform{
			Name: "test-platform-1",
		}

		platform2 := commands.Platform{
			Name: "test-platform-2",
		}
		platforms := []*commands.Platform{&platform2, &platform1}

		env.Client.EXPECT().GetInstalledPlatforms().Times(1).Return(platforms, nil)

		err := env.ArdiCore.Platform.ListInstalled()
		assert.NoError(st, err)

		assert.Contains(st, env.Stdout.String(), platform1.Name)
		assert.Contains(st, env.Stdout.String(), platform2.Name)
	})

	// @todo: check that list is actually sorted
	testutil.RunTest("prints sorted list of all available platforms to stdout", t, func(st *testing.T, env testutil.TestEnv) {
		defer env.Ctrl.Finish()
		platform1 := commands.Platform{
			Name: "test-platform-1",
		}

		platform2 := commands.Platform{
			Name: "test-platform-2",
		}
		platforms := []*commands.Platform{&platform2, &platform1}

		env.Client.EXPECT().GetPlatforms().Times(1).Return(platforms, nil)

		err := env.ArdiCore.Platform.ListAll()
		assert.NoError(st, err)

		assert.Contains(st, env.Stdout.String(), platform1.Name)
		assert.Contains(st, env.Stdout.String(), platform2.Name)
	})

	testutil.RunTest("adds platforms", t, func(st *testing.T, env testutil.TestEnv) {
		defer env.Ctrl.Finish()
		testPlatform1 := "test-platform1"
		testPlatform2 := "test-platform2"

		env.Client.EXPECT().InstallPlatform(testPlatform1).Times(1).Return(nil)
		env.Client.EXPECT().InstallPlatform(testPlatform2).Times(1).Return(nil)

		platforms := []string{testPlatform1, testPlatform2}

		for _, p := range platforms {
			err := env.ArdiCore.Platform.Add(p)
			assert.NoError(st, err)
		}
	})

	testutil.RunTest("removes a platforms", t, func(st *testing.T, env testutil.TestEnv) {
		defer env.Ctrl.Finish()
		testPlatform1 := "test-platform1"
		testPlatform2 := "test-platform2"

		env.Client.EXPECT().UninstallPlatform(testPlatform1).Times(1).Return(nil)
		env.Client.EXPECT().UninstallPlatform(testPlatform2).Times(1).Return(nil)

		platforms := []string{testPlatform1, testPlatform2}

		for _, p := range platforms {
			err := env.ArdiCore.Platform.Remove(p)
			assert.NoError(st, err)
		}
	})

	testutil.RunTest("adds all available platforms", t, func(st *testing.T, env testutil.TestEnv) {
		env.Client.EXPECT().InstallAllPlatforms().Times(1).Return(nil)
		err := env.ArdiCore.Platform.AddAll()
		assert.NoError(st, err)
	})
}
