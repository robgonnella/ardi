package core_test

import (
	"errors"
	"testing"

	rpc "github.com/arduino/arduino-cli/rpc/commands"
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
		req := &rpc.PlatformListReq{
			Instance:      instance,
			UpdatableOnly: false,
			All:           false,
		}

		env.Cli.EXPECT().CreateInstance().Return(instance, nil)
		env.Cli.EXPECT().GetPlatforms(req).Return(platforms, nil)

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
		req := &rpc.PlatformSearchReq{
			Instance:    instance,
			AllVersions: true,
		}
		resp := &rpc.PlatformSearchResp{
			SearchOutput: platforms,
		}
		env.Cli.EXPECT().CreateInstanceIgnorePlatformIndexErrors().Return(instance).Times(3)
		env.Cli.EXPECT().UpdateIndex(gomock.Any(), gomock.Any(), gomock.Any()).Times(2)
		env.Cli.EXPECT().UpdateLibrariesIndex(gomock.Any(), gomock.Any(), gomock.Any())
		env.Cli.EXPECT().CreateInstance().Return(instance, nil)
		env.Cli.EXPECT().PlatformSearch(req).Return(resp, nil)

		err := env.ArdiCore.Platform.ListAll()
		assert.NoError(env.T, err)

		assert.Contains(env.T, env.Stdout.String(), platform1.Name)
		assert.Contains(env.T, env.Stdout.String(), platform2.Name)
	})

	testutil.RunUnitTest("adds platforms", t, func(env *testutil.UnitTestEnv) {
		platform1 := &rpc.Platform{
			ID:        "test:platform1",
			Name:      "Platform1",
			Installed: "1.3.8",
		}

		platform2 := &rpc.Platform{
			ID:        "test:platform2",
			Name:      "Platform2",
			Installed: "3.1.9",
		}

		instance := &rpc.Instance{Id: int32(1)}
		req1 := &rpc.PlatformInstallReq{
			Instance:        instance,
			PlatformPackage: "test",
			Architecture:    "platform1",
		}
		req2 := &rpc.PlatformInstallReq{
			Instance:        instance,
			PlatformPackage: "test",
			Architecture:    "platform2",
		}

		listReq := &rpc.PlatformListReq{
			Instance:      instance,
			UpdatableOnly: false,
			All:           false,
		}
		platforms := []*rpc.Platform{platform1, platform2}

		env.Cli.EXPECT().CreateInstanceIgnorePlatformIndexErrors().Return(instance)
		env.Cli.EXPECT().UpdateIndex(gomock.Any(), gomock.Any(), gomock.Any())

		env.Cli.EXPECT().CreateInstance().Return(instance, nil)
		env.Cli.EXPECT().PlatformInstall(gomock.Any(), req1, gomock.Any(), gomock.Any())
		env.Cli.EXPECT().CreateInstance().Return(instance, nil)
		env.Cli.EXPECT().GetPlatforms(listReq).Return([]*rpc.Platform{platform1}, nil)

		env.Cli.EXPECT().CreateInstance().Return(instance, nil)
		env.Cli.EXPECT().PlatformInstall(gomock.Any(), req2, gomock.Any(), gomock.Any())
		env.Cli.EXPECT().CreateInstance().Return(instance, nil)
		env.Cli.EXPECT().GetPlatforms(listReq).Return(platforms, nil)

		for _, p := range platforms {
			installed, _, err := env.ArdiCore.Platform.Add(p.GetID())
			assert.NoError(env.T, err)
			assert.Equal(env.T, p.GetID(), installed)
		}
	})

	testutil.RunUnitTest("returns 'platform add' error", t, func(env *testutil.UnitTestEnv) {
		errString := "dummy error"
		dummyErr := errors.New(errString)

		platform1 := &rpc.Platform{
			ID:        "test:platform1",
			Name:      "Platform1",
			Installed: "1.3.8",
		}

		platform2 := &rpc.Platform{
			ID:        "test:platform2",
			Name:      "Platform2",
			Installed: "3.1.9",
		}

		instance := &rpc.Instance{Id: int32(1)}

		req1 := &rpc.PlatformInstallReq{
			Instance:        instance,
			PlatformPackage: "test",
			Architecture:    "platform1",
		}

		req2 := &rpc.PlatformInstallReq{
			Instance:        instance,
			PlatformPackage: "test",
			Architecture:    "platform2",
		}

		platforms := []*rpc.Platform{platform1, platform2}

		env.Cli.EXPECT().CreateInstanceIgnorePlatformIndexErrors().Return(instance)
		env.Cli.EXPECT().UpdateIndex(gomock.Any(), gomock.Any(), gomock.Any())

		env.Cli.EXPECT().CreateInstance().Return(instance, nil)
		env.Cli.EXPECT().PlatformInstall(gomock.Any(), req1, gomock.Any(), gomock.Any()).Return(nil, dummyErr)

		env.Cli.EXPECT().CreateInstance().Return(instance, nil)
		env.Cli.EXPECT().PlatformInstall(gomock.Any(), req2, gomock.Any(), gomock.Any()).Return(nil, dummyErr)

		for _, p := range platforms {
			_, _, err := env.ArdiCore.Platform.Add(p.GetID())
			assert.Error(env.T, err)
			assert.EqualError(env.T, err, errString)
		}
	})

	testutil.RunUnitTest("removes a platforms", t, func(env *testutil.UnitTestEnv) {
		platform1 := &rpc.Platform{
			ID:        "test:platform1",
			Name:      "Platform1",
			Installed: "1.3.8",
		}

		platform2 := &rpc.Platform{
			ID:        "test:platform2",
			Name:      "Platform2",
			Installed: "3.1.9",
		}

		instance := &rpc.Instance{Id: int32(1)}

		req1 := &rpc.PlatformUninstallReq{
			Instance:        instance,
			PlatformPackage: "test",
			Architecture:    "platform1",
		}

		req2 := &rpc.PlatformUninstallReq{
			Instance:        instance,
			PlatformPackage: "test",
			Architecture:    "platform2",
		}

		platforms := []*rpc.Platform{platform1, platform2}

		env.Cli.EXPECT().CreateInstance().Return(instance, nil)
		env.Cli.EXPECT().PlatformUninstall(gomock.Any(), req1, gomock.Any())

		env.Cli.EXPECT().CreateInstance().Return(instance, nil)
		env.Cli.EXPECT().PlatformUninstall(gomock.Any(), req2, gomock.Any())

		for _, p := range platforms {
			removed, err := env.ArdiCore.Platform.Remove(p.GetID())
			assert.NoError(env.T, err)
			assert.Equal(env.T, p.GetID(), removed)
		}
	})

	testutil.RunUnitTest("returns platform remove error", t, func(env *testutil.UnitTestEnv) {
		errString := "dummy error"
		dummyErr := errors.New(errString)

		platform1 := &rpc.Platform{
			ID:        "test:platform1",
			Name:      "Platform1",
			Installed: "1.3.8",
		}

		platform2 := &rpc.Platform{
			ID:        "test:platform2",
			Name:      "Platform2",
			Installed: "3.1.9",
		}

		instance := &rpc.Instance{Id: int32(1)}

		req1 := &rpc.PlatformUninstallReq{
			Instance:        instance,
			PlatformPackage: "test",
			Architecture:    "platform1",
		}

		req2 := &rpc.PlatformUninstallReq{
			Instance:        instance,
			PlatformPackage: "test",
			Architecture:    "platform2",
		}

		platforms := []*rpc.Platform{platform1, platform2}

		env.Cli.EXPECT().CreateInstance().Return(instance, nil)
		env.Cli.EXPECT().PlatformUninstall(gomock.Any(), req1, gomock.Any()).Return(nil, dummyErr)

		env.Cli.EXPECT().CreateInstance().Return(instance, nil)
		env.Cli.EXPECT().PlatformUninstall(gomock.Any(), req2, gomock.Any()).Return(nil, dummyErr)

		for _, p := range platforms {
			_, err := env.ArdiCore.Platform.Remove(p.GetID())
			assert.Error(env.T, err)
			assert.EqualError(env.T, err, errString)
		}
	})
}
