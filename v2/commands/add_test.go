package commands_test

import (
	"errors"
	"path"
	"testing"

	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/golang/mock/gomock"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAddPlatformCommand(t *testing.T) {
	pkg := "pkg"
	arch := "arch"
	version := "1.3.5"
	platform := pkg + ":" + arch + "@" + version

	instance := &rpc.Instance{Id: 1}

	installReq := &rpc.PlatformInstallRequest{
		Instance:        instance,
		PlatformPackage: pkg,
		Architecture:    arch,
		Version:         version,
	}

	platformReq := &rpc.PlatformListRequest{
		Instance: instance,
	}

	indexReq := &rpc.UpdateIndexRequest{
		Instance: instance,
	}

	testutil.RunMockIntegrationTest("errors if project not initialized", t, func(env *testutil.MockIntegrationTestEnv) {
		args := []string{"add", "platforms", "arduino:avr"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockIntegrationTest("adds platform", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
		env.ArduinoCli.EXPECT().UpdateIndex(gomock.Any(), indexReq, gomock.Any())
		env.ArduinoCli.EXPECT().GetPlatforms(platformReq)
		env.ArduinoCli.EXPECT().PlatformInstall(gomock.Any(), installReq, gomock.Any(), gomock.Any())

		args := []string{"add", "platform", platform}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunMockIntegrationTest("returns add platform error", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		dummyErr := errors.New("dummy error")

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
		env.ArduinoCli.EXPECT().UpdateIndex(gomock.Any(), indexReq, gomock.Any())
		env.ArduinoCli.EXPECT().PlatformInstall(gomock.Any(), installReq, gomock.Any(), gomock.Any()).Return(nil, dummyErr)

		args := []string{"add", "platform", platform}
		err = env.Execute(args)
		assert.Error(env.T, err)
		assert.ErrorIs(env.T, err, dummyErr)
	})
}

func TestAddLibraryCommand(t *testing.T) {
	instance := &rpc.Instance{Id: 1}
	library := "Some_Fancy_Library"
	version := "1.2.3"
	libStr := library + "@" + version

	installReq := &rpc.LibraryInstallRequest{
		Instance: instance,
		Name:     library,
		Version:  version,
	}

	listReq := &rpc.LibraryListRequest{
		Instance: instance,
	}

	indexReq := &rpc.UpdateLibrariesIndexRequest{
		Instance: instance,
	}

	testutil.RunMockIntegrationTest("errors if project not initialized", t, func(env *testutil.MockIntegrationTestEnv) {
		args := []string{"add", "lib", "Some Lib"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockIntegrationTest("adds library", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
		env.ArduinoCli.EXPECT().UpdateLibrariesIndex(gomock.Any(), indexReq, gomock.Any())
		env.ArduinoCli.EXPECT().LibraryList(gomock.Any(), listReq)
		env.ArduinoCli.EXPECT().LibraryInstall(gomock.Any(), installReq, gomock.Any(), gomock.Any())

		args := []string{"add", "lib", libStr}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunMockIntegrationTest("returns add library error", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		dummyErr := errors.New("dummy error")

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
		env.ArduinoCli.EXPECT().UpdateLibrariesIndex(gomock.Any(), indexReq, gomock.Any())
		env.ArduinoCli.EXPECT().LibraryInstall(gomock.Any(), installReq, gomock.Any(), gomock.Any()).Return(dummyErr)

		args := []string{"add", "lib", libStr}
		err = env.Execute(args)
		assert.Error(env.T, err)
		assert.ErrorIs(env.T, err, dummyErr)
	})
}

func TestAddBuildCommand(t *testing.T) {
	build := "cool"
	fqbn := "super:cool:fqbn"
	sketchDir := testutil.BlinkProjectDir()
	sketch := path.Join(sketchDir, "blink.ino")

	testutil.RunMockIntegrationTest("errors if project not initialized", t, func(env *testutil.MockIntegrationTestEnv) {
		args := []string{"add", "build", "-n", "default", "-f", "some:f:qbn", "-s", "sketch.ino"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockIntegrationTest("adds build", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		args := []string{"add", "build", "-n", build, "-f", fqbn, "-s", sketch}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		builds := env.ArdiCore.Config.GetBuilds()
		b, ok := builds[build]

		assert.True(env.T, ok)
		assert.Equal(env.T, b.FQBN, fqbn)
		assert.Equal(env.T, b.Sketch, sketch)
		assert.Equal(env.T, b.Directory, sketchDir)
	})

	testutil.RunMockIntegrationTest("return error if sketch not found", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		args := []string{"add", "build", "-n", build, "-f", fqbn, "-s", "noop.ino"}
		err = env.Execute(args)
		assert.Error(env.T, err)

		builds := env.ArdiCore.Config.GetBuilds()
		b, ok := builds[build]

		assert.False(env.T, ok)
		assert.Empty(env.T, b.FQBN)
		assert.Empty(env.T, b.Sketch)
		assert.Empty(env.T, b.Directory)
	})
}

func TestAddBoardURLCommand(t *testing.T) {
	testutil.RunMockIntegrationTest("errors if project not initialized", t, func(env *testutil.MockIntegrationTestEnv) {
		args := []string{"add", "board-url", "https://someboardurl.com"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockIntegrationTest("adds board url", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		url := "https://somefancyboardurl.com"
		args := []string{"add", "board-url", url}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		urls := env.ArdiCore.CliConfig.Config.BoardManager.AdditionalUrls

		assert.NotEmpty(env.T, urls)
		assert.Contains(env.T, urls, url)
	})

	testutil.RunMockIntegrationTest("adds multiple board url", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		url1 := "https://somefancyboardurl.com"
		url2 := "https://anotherfancyboardurl.com"

		args := []string{"add", "board-url", url1}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"add", "board-url", url2}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		urls := env.ArdiCore.CliConfig.Config.BoardManager.AdditionalUrls

		assert.NotEmpty(env.T, urls)
		assert.Contains(env.T, urls, url1)
		assert.Contains(env.T, urls, url2)
	})

	testutil.RunMockIntegrationTest("does not error if board url already added", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		url := "https://somefancyboardurl.com"

		args := []string{"add", "board-url", url}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"add", "board-url", url}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		urls := env.ArdiCore.CliConfig.Config.BoardManager.AdditionalUrls

		assert.NotEmpty(env.T, urls)
		assert.Contains(env.T, urls, url)
		assert.Equal(env.T, len(urls), 1)
	})
}
