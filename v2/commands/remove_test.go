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

func TestRemovePlatformCommand(t *testing.T) {
	pkg := "pkg"
	arch := "arch"
	version := "1.3.5"
	platform := pkg + ":" + arch + "@" + version

	instance := &rpc.Instance{Id: 1}

	uninstallReq := &rpc.PlatformUninstallRequest{
		Instance:        instance,
		PlatformPackage: pkg,
		Architecture:    arch,
	}

	testutil.RunMockIntegrationTest("errors if project not initialized", t, func(env *testutil.MockIntegrationTestEnv) {
		args := []string{"add", "platforms", "arduino:avr"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockIntegrationTest("removes platform", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
		env.ArduinoCli.EXPECT().PlatformUninstall(gomock.Any(), uninstallReq, gomock.Any())

		args := []string{"remove", "platform", platform}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunMockIntegrationTest("returns remove platform error", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		dummyErr := errors.New("dummy error")

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
		env.ArduinoCli.EXPECT().PlatformUninstall(gomock.Any(), uninstallReq, gomock.Any()).Return(nil, dummyErr)

		args := []string{"remove", "platform", platform}
		err = env.Execute(args)
		assert.Error(env.T, err)
		assert.ErrorIs(env.T, err, dummyErr)
	})
}

func TestRemoveLibraryCommand(t *testing.T) {
	instance := &rpc.Instance{Id: 1}
	library := "Some_Fancy_Library"

	uninstallReq := &rpc.LibraryUninstallRequest{
		Instance: instance,
		Name:     library,
	}

	testutil.RunMockIntegrationTest("errors if project not initialized", t, func(env *testutil.MockIntegrationTestEnv) {
		args := []string{"add", "lib", "Some Lib"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockIntegrationTest("removes library", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
		env.ArduinoCli.EXPECT().LibraryUninstall(gomock.Any(), uninstallReq, gomock.Any())

		args := []string{"remove", "lib", library}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunMockIntegrationTest("returns remove library error", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		dummyErr := errors.New("dummy error")

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
		env.ArduinoCli.EXPECT().LibraryUninstall(gomock.Any(), uninstallReq, gomock.Any()).Return(dummyErr)

		args := []string{"remove", "lib", library}
		err = env.Execute(args)
		assert.Error(env.T, err)
		assert.ErrorIs(env.T, err, dummyErr)
	})
}

func TestRemoveBuildCommand(t *testing.T) {
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

		args = []string{"remove", "build", build}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		builds = env.ArdiCore.Config.GetBuilds()
		b, ok = builds[build]
		assert.False(env.T, ok)
		assert.Empty(env.T, b)
	})
}

func TestRemoveBoardURLCommand(t *testing.T) {
	testutil.RunMockIntegrationTest("errors if project not initialized", t, func(env *testutil.MockIntegrationTestEnv) {
		args := []string{"add", "board-url", "https://someboardurl.com"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockIntegrationTest("removes board url", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		url := "https://somefancyboardurl.com"
		args := []string{"add", "board-url", url}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		urls := env.ArdiCore.CliConfig.Config.BoardManager.AdditionalUrls

		assert.NotEmpty(env.T, urls)
		assert.Contains(env.T, urls, url)

		args = []string{"remove", "board-url", url}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		urls = env.ArdiCore.CliConfig.Config.BoardManager.AdditionalUrls
		assert.NotContains(env.T, urls, url)
		assert.Empty(env.T, urls)
	})

	testutil.RunMockIntegrationTest("removes multiple board urls", t, func(env *testutil.MockIntegrationTestEnv) {
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

		args = []string{"remove", "board-url", url1, url2}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		urls = env.ArdiCore.CliConfig.Config.BoardManager.AdditionalUrls
		assert.NotContains(env.T, urls, url1)
		assert.NotContains(env.T, urls, url2)
		assert.Empty(env.T, urls)
	})

	testutil.RunMockIntegrationTest("does not error if no board url exists", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		url := "https://somefancyboardurl.com"

		args := []string{"remove", "board-url", url}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		urls := env.ArdiCore.CliConfig.Config.BoardManager.AdditionalUrls
		assert.Empty(env.T, urls)
	})
}
