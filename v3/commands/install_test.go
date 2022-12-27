package commands_test

import (
	"errors"
	"testing"

	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/golang/mock/gomock"
	"github.com/robgonnella/ardi/v3/testutil"
	"github.com/stretchr/testify/assert"
)

func TestInstallCommand(t *testing.T) {
	instance := &rpc.Instance{Id: 1}

	lib := "Some_Library"
	libVers := "3.4.1"

	platform := "some:platform"
	platformVers := "3.1.0"

	installPlatReq := &rpc.PlatformInstallRequest{
		Instance:        instance,
		PlatformPackage: "some",
		Architecture:    "platform",
		Version:         platformVers,
	}

	installLibReq := &rpc.LibraryInstallRequest{
		Instance: instance,
		Name:     lib,
		Version:  libVers,
	}

	platformListReq := &rpc.PlatformListRequest{
		Instance:      instance,
		UpdatableOnly: false,
		All:           false,
	}

	libraryListReq := &rpc.LibraryListRequest{
		Instance: instance,
	}

	indexReq := &rpc.UpdateIndexRequest{
		Instance: instance,
	}

	libIndexReq := &rpc.UpdateLibrariesIndexRequest{
		Instance: instance,
	}

	expectUsual := func(env *testutil.MockIntegrationTestEnv) {
		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
		env.ArduinoCli.EXPECT().UpdateIndex(gomock.Any(), indexReq, gomock.Any())
	}

	testutil.RunIntegrationTest("errors if project not initialized", t, func(env *testutil.IntegrationTestEnv) {
		args := []string{"install"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockIntegrationTest("installs all dependencies listed in ardi.json", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		err = env.ArdiCore.Config.AddLibrary(lib, libVers)
		assert.NoError(env.T, err)

		err = env.ArdiCore.Config.AddPlatform(platform, platformVers)
		assert.NoError(env.T, err)

		err = env.ArdiCore.Config.AddBoardURL(testutil.Esp8266BoardURL())
		assert.NoError(env.T, err)

		expectUsual(env)
		env.ArduinoCli.EXPECT().PlatformInstall(gomock.Any(), installPlatReq, gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().GetPlatforms(platformListReq)
		env.ArduinoCli.EXPECT().UpdateLibrariesIndex(gomock.Any(), libIndexReq, gomock.Any())
		env.ArduinoCli.EXPECT().LibraryInstall(gomock.Any(), installLibReq, gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().LibraryList(gomock.Any(), libraryListReq)

		args := []string{"install"}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		boardURLs := env.ArdiCore.CliConfig.Config.BoardManager.AdditionalUrls
		assert.Contains(env.T, boardURLs, testutil.Esp8266BoardURL())
	})

	testutil.RunMockIntegrationTest("returns platform install error", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		dummyErr := errors.New("dummy error")

		expectUsual(env)
		env.ArduinoCli.EXPECT().PlatformInstall(gomock.Any(), installPlatReq, gomock.Any(), gomock.Any()).Return(nil, dummyErr)

		err = env.ArdiCore.Config.AddLibrary(lib, libVers)
		assert.NoError(env.T, err)

		err = env.ArdiCore.Config.AddPlatform(platform, platformVers)
		assert.NoError(env.T, err)

		args := []string{"install"}
		err = env.Execute(args)
		assert.Error(env.T, err)
		assert.ErrorIs(env.T, err, dummyErr)
	})

	testutil.RunMockIntegrationTest("returns library install error", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		dummyErr := errors.New("dummy error")

		expectUsual(env)
		env.ArduinoCli.EXPECT().PlatformInstall(gomock.Any(), installPlatReq, gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().GetPlatforms(platformListReq)
		env.ArduinoCli.EXPECT().UpdateLibrariesIndex(gomock.Any(), libIndexReq, gomock.Any())
		env.ArduinoCli.EXPECT().LibraryInstall(gomock.Any(), installLibReq, gomock.Any(), gomock.Any()).Return(dummyErr)

		err = env.ArdiCore.Config.AddLibrary(lib, libVers)
		assert.NoError(env.T, err)

		err = env.ArdiCore.Config.AddPlatform(platform, platformVers)
		assert.NoError(env.T, err)

		args := []string{"install"}
		err = env.Execute(args)
		assert.Error(env.T, err)
		assert.ErrorIs(env.T, err, dummyErr)
	})
}
