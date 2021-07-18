package core_test

import (
	"errors"
	"testing"

	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/golang/mock/gomock"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

// @todo: check that list is actually sorted
func TestPlatformCore(t *testing.T) {
	testutil.RunUnitTest("prints sorted list of all installed platforms to stdout", t, func(env *testutil.UnitTestEnv) {
		platform1 := rpc.Platform{
			Name: "test-platform-1",
		}

		platform2 := rpc.Platform{
			Name: "test-platform-2",
		}
		platforms := []*rpc.Platform{&platform2, &platform1}

		instance := &rpc.Instance{Id: int32(1)}
		req := &rpc.PlatformListRequest{
			Instance:      instance,
			UpdatableOnly: false,
			All:           false,
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().GetPlatforms(req).Return(platforms, nil)

		err := env.ArdiCore.Platform.ListInstalled()
		assert.NoError(env.T, err)

		assert.Contains(env.T, env.Stdout.String(), platform1.Name)
		assert.Contains(env.T, env.Stdout.String(), platform2.Name)
	})

	// @todo: check that list is actually sorted
	testutil.RunUnitTest("prints sorted list of all available platforms to stdout", t, func(env *testutil.UnitTestEnv) {
		platform1 := rpc.Platform{
			Name: "test-platform-1",
		}

		platform2 := rpc.Platform{
			Name: "test-platform-2",
		}
		platforms := []*rpc.Platform{&platform2, &platform1}

		instance := &rpc.Instance{Id: int32(1)}
		req := &rpc.PlatformSearchRequest{
			Instance:    instance,
			AllVersions: true,
		}
		resp := &rpc.PlatformSearchResponse{
			SearchOutput: platforms,
		}
		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().UpdateIndex(gomock.Any(), gomock.Any(), gomock.Any()).Times(2)
		env.ArduinoCli.EXPECT().UpdateLibrariesIndex(gomock.Any(), gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().PlatformSearch(req).Return(resp, nil)

		err := env.ArdiCore.Platform.ListAll()
		assert.NoError(env.T, err)

		assert.Contains(env.T, env.Stdout.String(), platform1.Name)
		assert.Contains(env.T, env.Stdout.String(), platform2.Name)
	})

	testutil.RunUnitTest("adds platforms", t, func(env *testutil.UnitTestEnv) {
		platform1 := &rpc.Platform{
			Id:        "test:platform1",
			Name:      "Platform1",
			Installed: "1.3.8",
		}

		platform2 := &rpc.Platform{
			Id:        "test:platform2",
			Name:      "Platform2",
			Installed: "3.1.9",
		}

		instance := &rpc.Instance{Id: int32(1)}
		req1 := &rpc.PlatformInstallRequest{
			Instance:        instance,
			PlatformPackage: "test",
			Architecture:    "platform1",
		}
		req2 := &rpc.PlatformInstallRequest{
			Instance:        instance,
			PlatformPackage: "test",
			Architecture:    "platform2",
		}

		listReq := &rpc.PlatformListRequest{
			Instance:      instance,
			UpdatableOnly: false,
			All:           false,
		}
		platforms := []*rpc.Platform{platform1, platform2}

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().UpdateIndex(gomock.Any(), gomock.Any(), gomock.Any())

		env.ArduinoCli.EXPECT().PlatformInstall(gomock.Any(), req1, gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().GetPlatforms(listReq).Return([]*rpc.Platform{platform1}, nil)

		env.ArduinoCli.EXPECT().PlatformInstall(gomock.Any(), req2, gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().GetPlatforms(listReq).Return(platforms, nil)

		for _, p := range platforms {
			installed, _, err := env.ArdiCore.Platform.Add(p.GetId())
			assert.NoError(env.T, err)
			assert.Equal(env.T, p.GetId(), installed)
		}
	})

	testutil.RunUnitTest("returns 'platform add' error", t, func(env *testutil.UnitTestEnv) {
		errString := "dummy error"
		dummyErr := errors.New(errString)

		platform1 := &rpc.Platform{
			Id:        "test:platform1",
			Name:      "Platform1",
			Installed: "1.3.8",
		}

		platform2 := &rpc.Platform{
			Id:        "test:platform2",
			Name:      "Platform2",
			Installed: "3.1.9",
		}

		instance := &rpc.Instance{Id: int32(1)}

		req1 := &rpc.PlatformInstallRequest{
			Instance:        instance,
			PlatformPackage: "test",
			Architecture:    "platform1",
		}

		req2 := &rpc.PlatformInstallRequest{
			Instance:        instance,
			PlatformPackage: "test",
			Architecture:    "platform2",
		}

		platforms := []*rpc.Platform{platform1, platform2}

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().UpdateIndex(gomock.Any(), gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().PlatformInstall(gomock.Any(), req1, gomock.Any(), gomock.Any()).Return(nil, dummyErr)
		env.ArduinoCli.EXPECT().PlatformInstall(gomock.Any(), req2, gomock.Any(), gomock.Any()).Return(nil, dummyErr)

		for _, p := range platforms {
			_, _, err := env.ArdiCore.Platform.Add(p.GetId())
			assert.Error(env.T, err)
			assert.EqualError(env.T, err, errString)
		}
	})

	testutil.RunUnitTest("removes a platforms", t, func(env *testutil.UnitTestEnv) {
		platform1 := &rpc.Platform{
			Id:        "test:platform1",
			Name:      "Platform1",
			Installed: "1.3.8",
		}

		platform2 := &rpc.Platform{
			Id:        "test:platform2",
			Name:      "Platform2",
			Installed: "3.1.9",
		}

		instance := &rpc.Instance{Id: int32(1)}

		req1 := &rpc.PlatformUninstallRequest{
			Instance:        instance,
			PlatformPackage: "test",
			Architecture:    "platform1",
		}

		req2 := &rpc.PlatformUninstallRequest{
			Instance:        instance,
			PlatformPackage: "test",
			Architecture:    "platform2",
		}

		platforms := []*rpc.Platform{platform1, platform2}

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().PlatformUninstall(gomock.Any(), req1, gomock.Any())
		env.ArduinoCli.EXPECT().PlatformUninstall(gomock.Any(), req2, gomock.Any())

		for _, p := range platforms {
			removed, err := env.ArdiCore.Platform.Remove(p.GetId())
			assert.NoError(env.T, err)
			assert.Equal(env.T, p.GetId(), removed)
		}
	})

	testutil.RunUnitTest("returns platform remove error", t, func(env *testutil.UnitTestEnv) {
		errString := "dummy error"
		dummyErr := errors.New(errString)

		platform1 := &rpc.Platform{
			Id:        "test:platform1",
			Name:      "Platform1",
			Installed: "1.3.8",
		}

		platform2 := &rpc.Platform{
			Id:        "test:platform2",
			Name:      "Platform2",
			Installed: "3.1.9",
		}

		instance := &rpc.Instance{Id: int32(1)}

		req1 := &rpc.PlatformUninstallRequest{
			Instance:        instance,
			PlatformPackage: "test",
			Architecture:    "platform1",
		}

		req2 := &rpc.PlatformUninstallRequest{
			Instance:        instance,
			PlatformPackage: "test",
			Architecture:    "platform2",
		}

		platforms := []*rpc.Platform{platform1, platform2}

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().PlatformUninstall(gomock.Any(), req1, gomock.Any()).Return(nil, dummyErr)
		env.ArduinoCli.EXPECT().PlatformUninstall(gomock.Any(), req2, gomock.Any()).Return(nil, dummyErr)

		for _, p := range platforms {
			_, err := env.ArdiCore.Platform.Remove(p.GetId())
			assert.Error(env.T, err)
			assert.EqualError(env.T, err, errString)
		}
	})
}
