package commands_test

import (
	"errors"
	"testing"

	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/golang/mock/gomock"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestSearchLibCommand(t *testing.T) {
	instance := &rpc.Instance{Id: 1}

	indexReq := &rpc.UpdateLibrariesIndexRequest{
		Instance: instance,
	}

	testutil.RunMockIntegrationTest("errors if project not initialized", t, func(env *testutil.MockIntegrationTestEnv) {
		searchLib := "Adafruit Pixie"
		args := []string{"search", "libs", searchLib}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockIntegrationTest("searches for library", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		searchLib := "Adafruit Pixie"

		searchReq := &rpc.LibrarySearchRequest{
			Instance: instance,
			Query:    searchLib,
		}

		expectedLib := &rpc.SearchedLibrary{
			Name: searchLib,
		}

		searchResp := &rpc.LibrarySearchResponse{
			Libraries: []*rpc.SearchedLibrary{expectedLib},
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
		env.ArduinoCli.EXPECT().UpdateLibrariesIndex(gomock.Any(), indexReq, gomock.Any())
		env.ArduinoCli.EXPECT().LibrarySearch(gomock.Any(), searchReq).Return(searchResp, nil)

		args := []string{"search", "libs", searchLib}
		err = env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), expectedLib.Name)
	})

	testutil.RunMockIntegrationTest("returns search error", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		searchLib := "Adafruit Pixie"

		searchReq := &rpc.LibrarySearchRequest{
			Instance: instance,
			Query:    searchLib,
		}

		dummyErr := errors.New("dummy err")

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
		env.ArduinoCli.EXPECT().UpdateLibrariesIndex(gomock.Any(), indexReq, gomock.Any())
		env.ArduinoCli.EXPECT().LibrarySearch(gomock.Any(), searchReq).Return(nil, dummyErr)

		args := []string{"search", "libs", searchLib}
		err = env.Execute(args)
		assert.Error(env.T, err)
		assert.ErrorIs(env.T, err, dummyErr)
	})

	testutil.RunMockIntegrationTest("does not error if search arg not provided", t, func(env *testutil.MockIntegrationTestEnv) {
		err := env.RunProjectInit()
		assert.NoError(env.T, err)

		searchLib := ""

		searchReq := &rpc.LibrarySearchRequest{
			Instance: instance,
			Query:    searchLib,
		}

		expectedLib := &rpc.SearchedLibrary{
			Name: searchLib,
		}

		searchResp := &rpc.LibrarySearchResponse{
			Libraries: []*rpc.SearchedLibrary{expectedLib},
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
		env.ArduinoCli.EXPECT().UpdateLibrariesIndex(gomock.Any(), indexReq, gomock.Any())
		env.ArduinoCli.EXPECT().LibrarySearch(gomock.Any(), searchReq).Return(searchResp, nil)

		args := []string{"search", "libs", searchLib}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})
}

func TestSearchPlatformCommand(t *testing.T) {
	instance := &rpc.Instance{Id: 1}

	indexReq := &rpc.UpdateIndexRequest{
		Instance: instance,
	}

	testutil.RunMockIntegrationTest("errors if project not initialized", t, func(env *testutil.MockIntegrationTestEnv) {
		args := []string{"search", "platforms"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockIntegrationTest("searches platforms", t, func(env *testutil.MockIntegrationTestEnv) {
		env.RunProjectInit()

		platform := "some:platform"

		searchReq := &rpc.PlatformSearchRequest{
			Instance:    instance,
			AllVersions: true,
		}

		searchResp := &rpc.PlatformSearchResponse{
			SearchOutput: []*rpc.Platform{
				{
					Id: platform,
				},
			},
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
		env.ArduinoCli.EXPECT().UpdateIndex(gomock.Any(), indexReq, gomock.Any()).MaxTimes(2)
		env.ArduinoCli.EXPECT().PlatformSearch(searchReq).Return(searchResp, nil)

		args := []string{"search", "platforms"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), platform)
	})

	testutil.RunMockIntegrationTest("return platform search error", t, func(env *testutil.MockIntegrationTestEnv) {
		env.RunProjectInit()
		searchReq := &rpc.PlatformSearchRequest{
			Instance:    instance,
			AllVersions: true,
		}

		dummyErr := errors.New("dummy error")

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
		env.ArduinoCli.EXPECT().UpdateIndex(gomock.Any(), indexReq, gomock.Any()).MaxTimes(2)
		env.ArduinoCli.EXPECT().PlatformSearch(searchReq).Return(nil, dummyErr)

		args := []string{"search", "platforms"}
		err := env.Execute(args)
		assert.Error(env.T, err)
		assert.ErrorIs(env.T, err, dummyErr)
	})
}
