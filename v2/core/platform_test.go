package core_test

import (
	"errors"
	"testing"

	"github.com/arduino/arduino-cli/rpc/commands"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

// @todo: check that list is actually sorted
func TestPlatformCore(t *testing.T) {
	testutil.RunUnitTest("prints sorted list of all installed platforms to stdout", t, func(env *testutil.UnitTestEnv) {
		platform1 := commands.Platform{
			Name: "test-platform-1",
		}

		platform2 := commands.Platform{
			Name: "test-platform-2",
		}
		platforms := []*commands.Platform{&platform2, &platform1}

		env.Client.EXPECT().GetInstalledPlatforms().Times(1).Return(platforms, nil)

		err := env.ArdiCore.Platform.ListInstalled()
		assert.NoError(env.T, err)

		assert.Contains(env.T, env.Stdout.String(), platform1.Name)
		assert.Contains(env.T, env.Stdout.String(), platform2.Name)
	})

	// @todo: check that list is actually sorted
	testutil.RunUnitTest("prints sorted list of all available platforms to stdout", t, func(env *testutil.UnitTestEnv) {
		platform1 := commands.Platform{
			Name: "test-platform-1",
		}

		platform2 := commands.Platform{
			Name: "test-platform-2",
		}
		platforms := []*commands.Platform{&platform2, &platform1}

		env.Client.EXPECT().UpdatePlatformIndex().Times(1).Return(nil)
		env.Client.EXPECT().GetPlatforms().Times(1).Return(platforms, nil)

		err := env.ArdiCore.Platform.ListAll()
		assert.NoError(env.T, err)

		assert.Contains(env.T, env.Stdout.String(), platform1.Name)
		assert.Contains(env.T, env.Stdout.String(), platform2.Name)
	})

	testutil.RunUnitTest("adds platforms", t, func(env *testutil.UnitTestEnv) {
		testPlatform1 := "test-platform1"
		testPlatform2 := "test-platform2"

		env.Client.EXPECT().UpdatePlatformIndex().Times(1).Return(nil)
		env.Client.EXPECT().InstallPlatform(testPlatform1).Times(1).Return(testPlatform1, "1.0.0", nil)
		env.Client.EXPECT().InstallPlatform(testPlatform2).Times(1).Return(testPlatform2, "2.1.2", nil)

		platforms := []string{testPlatform1, testPlatform2}

		for _, p := range platforms {
			installed, _, err := env.ArdiCore.Platform.Add(p)
			assert.NoError(env.T, err)
			assert.Equal(env.T, p, installed)
		}
	})

	testutil.RunUnitTest("returns 'platform add' error", t, func(env *testutil.UnitTestEnv) {
		errString := "dummy error"
		dummyErr := errors.New(errString)

		testPlatform1 := "test-platform1"
		testPlatform2 := "test-platform2"

		env.Client.EXPECT().UpdatePlatformIndex().Times(1).Return(nil)
		env.Client.EXPECT().InstallPlatform(testPlatform1).Times(1).Return("", "", dummyErr)
		env.Client.EXPECT().InstallPlatform(testPlatform2).Times(1).Return("", "", dummyErr)

		platforms := []string{testPlatform1, testPlatform2}

		for _, p := range platforms {
			_, _, err := env.ArdiCore.Platform.Add(p)
			assert.Error(env.T, err)
			assert.EqualError(env.T, err, errString)
		}
	})

	testutil.RunUnitTest("removes a platforms", t, func(env *testutil.UnitTestEnv) {
		testPlatform1 := "test-platform1"
		testPlatform2 := "test-platform2"

		env.Client.EXPECT().UninstallPlatform(testPlatform1).Times(1).Return(testPlatform1, nil)
		env.Client.EXPECT().UninstallPlatform(testPlatform2).Times(1).Return(testPlatform2, nil)

		platforms := []string{testPlatform1, testPlatform2}

		for _, p := range platforms {
			removed, err := env.ArdiCore.Platform.Remove(p)
			assert.NoError(env.T, err)
			assert.Equal(env.T, p, removed)
		}
	})

	testutil.RunUnitTest("returns platform remove error", t, func(env *testutil.UnitTestEnv) {
		errString := "dummy error"
		dummyErr := errors.New(errString)

		testPlatform1 := "test-platform1"
		testPlatform2 := "test-platform2"

		env.Client.EXPECT().UninstallPlatform(testPlatform1).Times(1).Return("", dummyErr)
		env.Client.EXPECT().UninstallPlatform(testPlatform2).Times(1).Return("", dummyErr)

		platforms := []string{testPlatform1, testPlatform2}

		for _, p := range platforms {
			_, err := env.ArdiCore.Platform.Remove(p)
			assert.Error(env.T, err)
			assert.EqualError(env.T, err, errString)
		}
	})
}
