package commands_test

import (
	"errors"
	"path"
	"testing"

	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/golang/mock/gomock"
	"github.com/robgonnella/ardi/v3/testutil"
	"github.com/stretchr/testify/assert"
)

func TestListPlatformCommand(t *testing.T) {
	instance := &rpc.Instance{Id: 1}

	platformReq := &rpc.PlatformListRequest{
		Instance:      instance,
		UpdatableOnly: false,
		All:           false,
	}

	testutil.RunMockIntegrationTest("errors if project not initialized", t, func(env *testutil.MockIntegrationTestEnv) {
		args := []string{"list", "platforms"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockIntegrationTest("lists platforms", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		expectedPlatform := &rpc.Platform{
			Id:        "cool:platform",
			Installed: "1.2.3",
			Latest:    "1.2.3",
			Name:      "Super Cool Platform",
		}

		expectedPlatforms := []*rpc.Platform{expectedPlatform}

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
		env.ArduinoCli.EXPECT().GetPlatforms(platformReq).Return(expectedPlatforms, nil)

		args := []string{"list", "platforms"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), expectedPlatform.Name)
	})

	testutil.RunMockIntegrationTest("returns list platforms error", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		dummyErr := errors.New("dummy error")

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
		env.ArduinoCli.EXPECT().GetPlatforms(platformReq).Return(nil, dummyErr)

		args := []string{"list", "platforms"}
		err = env.Execute(args)
		assert.Error(env.T, err)
		assert.ErrorIs(env.T, err, dummyErr)
	})
}

func TestListLibraryCommand(t *testing.T) {
	instance := &rpc.Instance{Id: 1}

	listReq := &rpc.LibraryListRequest{
		Instance: instance,
	}

	testutil.RunMockIntegrationTest("errors if project not initialized", t, func(env *testutil.MockIntegrationTestEnv) {
		args := []string{"list", "libs"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockIntegrationTest("lists installed libraries", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		installedLib := &rpc.InstalledLibrary{
			Library: &rpc.Library{
				Name:    "Some_Cool_Library",
				Version: "1.3.5",
			},
		}
		listResp := &rpc.LibraryListResponse{
			InstalledLibraries: []*rpc.InstalledLibrary{installedLib},
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
		env.ArduinoCli.EXPECT().LibraryList(gomock.Any(), listReq).Return(listResp, nil)

		args := []string{"list", "libs"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), installedLib.Library.Name)
		assert.Contains(env.T, env.Stdout.String(), installedLib.Library.Version)
	})

	testutil.RunMockIntegrationTest("returns list libs error", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		dummyErr := errors.New("dummy error")

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
		env.ArduinoCli.EXPECT().LibraryList(gomock.Any(), listReq).Return(nil, dummyErr)

		args := []string{"list", "libs"}
		err = env.Execute(args)
		assert.Error(env.T, err)
		assert.ErrorIs(env.T, err, dummyErr)
	})
}

func TestListBuildCommand(t *testing.T) {
	build := "cool"
	fqbn := "super:cool:fqbn"
	sketchDir := testutil.BlinkProjectDir()
	sketch := path.Join(sketchDir, "blink.ino")

	testutil.RunMockIntegrationTest("errors if project not initialized", t, func(env *testutil.MockIntegrationTestEnv) {
		args := []string{"list", "builds"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockIntegrationTest("lists builds", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		args := []string{"add", "build", "-n", build, "-f", fqbn, "-s", sketch}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"list", "builds"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), build)
		assert.Contains(env.T, env.Stdout.String(), fqbn)
		assert.Contains(env.T, env.Stdout.String(), sketchDir)
		assert.Contains(env.T, env.Stdout.String(), sketch)
	})

	testutil.RunMockIntegrationTest("doesnt error if no builds to list", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		args := []string{"list", "builds"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})
}

func TestListBoardURLCommand(t *testing.T) {
	testutil.RunMockIntegrationTest("errors if project not initialized", t, func(env *testutil.MockIntegrationTestEnv) {
		args := []string{"list", "board-urls"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockIntegrationTest("lists board urls", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		boardUrl := "https://somecoolboardurl.com"
		args := []string{"add", "board-url", boardUrl}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"list", "board-urls"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), boardUrl)
	})

	testutil.RunMockIntegrationTest("doesnt error if no board urls to list", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		args := []string{"list", "board-urls"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})
}
