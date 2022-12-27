package cli_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/golang/mock/gomock"
	"github.com/robgonnella/ardi/v3/cli-wrapper"
	"github.com/robgonnella/ardi/v3/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type cliTestEnv struct {
	Ctx        context.Context
	Ctrl       *gomock.Controller
	ArduinoCli *mocks.MockCli
	CliWrapper *cli.Wrapper
	Logger     *logrus.Logger
	Stdout     *bytes.Buffer
}

// ClearStdout clears stdout for unit test
func (e *cliTestEnv) ClearStdout() {
	var b bytes.Buffer
	e.Logger.SetOutput(&b)
	e.Stdout = &b
}

func runCliTest(name string, t *testing.T, fn func(env cliTestEnv, t2 *testing.T)) {
	t.Run(name, func(st *testing.T) {
		logger := logrus.New()
		var b bytes.Buffer
		logger.SetOutput(&b)

		ctrl := gomock.NewController(st)
		mockArduinoCli := mocks.NewMockCli(ctrl)
		settingsPath := "."
		ctx := context.Background()

		mockArduinoCli.EXPECT().InitSettings(settingsPath)

		withMockCli := cli.WithArduinoCli(mockArduinoCli)
		cliWrapper := cli.NewCli(ctx, settingsPath, logger, withMockCli)

		env := cliTestEnv{
			Ctx:        ctx,
			Ctrl:       ctrl,
			ArduinoCli: mockArduinoCli,
			CliWrapper: cliWrapper,
			Logger:     logger,
			Stdout:     &b,
		}

		fn(env, st)
	})
}

func TestCliWrapperTest(t *testing.T) {
	runCliTest("updates both indexes", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}
		platIndexReq := &rpc.UpdateIndexRequest{Instance: inst}
		libIndexReq := &rpc.UpdateLibrariesIndexRequest{Instance: inst}

		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().UpdateIndex(gomock.Any(), platIndexReq, gomock.Any())
		env.ArduinoCli.EXPECT().UpdateLibrariesIndex(gomock.Any(), libIndexReq, gomock.Any())

		err := env.CliWrapper.UpdateIndexFiles()
		assert.NoError(st, err)
	})

	runCliTest("updates library index", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}
		libIndexReq := &rpc.UpdateLibrariesIndexRequest{Instance: inst}

		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().UpdateLibrariesIndex(gomock.Any(), libIndexReq, gomock.Any())

		err := env.CliWrapper.UpdateLibraryIndex()
		assert.NoError(st, err)
	})

	runCliTest("updates platform index", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}
		platIndexReq := &rpc.UpdateIndexRequest{Instance: inst}

		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().UpdateIndex(gomock.Any(), platIndexReq, gomock.Any())

		err := env.CliWrapper.UpdatePlatformIndex()
		assert.NoError(st, err)
	})

	runCliTest("installs platforms", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}
		pkg := "something"
		arch := "validish"
		version := "2.3.5"
		platform := fmt.Sprintf("%s:%s", pkg, arch)
		platformWithVersion := fmt.Sprintf("%s@%s", platform, version)

		env.Logger.SetLevel(logrus.DebugLevel)

		req := &rpc.PlatformInstallRequest{
			Instance:        inst,
			PlatformPackage: pkg,
			Architecture:    arch,
			Version:         version,
		}

		listReq := &rpc.PlatformListRequest{
			Instance:      inst,
			UpdatableOnly: false,
			All:           false,
		}
		installed := []*rpc.Platform{
			{
				Id:        platform,
				Installed: version,
			},
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().PlatformInstall(gomock.Any(), req, gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().GetPlatforms(listReq).Return(installed, nil)

		installedPlat, installedVers, err := env.CliWrapper.InstallPlatform(platformWithVersion)
		assert.NoError(st, err)
		assert.Equal(st, platform, installedPlat)
		assert.Equal(st, version, installedVers)
	})

	runCliTest("uninstalls platforms", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}
		pkg := "something"
		arch := "validish"
		platform := fmt.Sprintf("%s:%s", pkg, arch)
		req := &rpc.PlatformUninstallRequest{
			Instance:        inst,
			PlatformPackage: pkg,
			Architecture:    arch,
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().PlatformUninstall(gomock.Any(), req, gomock.Any())

		removedPlatform, err := env.CliWrapper.UninstallPlatform(platform)
		assert.NoError(st, err)
		assert.Equal(st, platform, removedPlatform)
	})

	runCliTest("returns installed platforms", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}
		listReq := &rpc.PlatformListRequest{
			Instance:      inst,
			UpdatableOnly: false,
			All:           false,
		}
		installed := []*rpc.Platform{
			{
				Id:        "some:platform",
				Installed: "2.3.6",
			},
		}
		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().GetPlatforms(listReq).Return(installed, nil)

		list, err := env.CliWrapper.GetInstalledPlatforms()
		assert.NoError(st, err)
		assert.Equal(st, installed, list)
	})

	runCliTest("returns all platforms", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}

		platIndexReq := &rpc.UpdateIndexRequest{Instance: inst}

		searchReq := &rpc.PlatformSearchRequest{
			Instance:    inst,
			AllVersions: false,
		}

		expectedResp := &rpc.PlatformSearchResponse{
			SearchOutput: []*rpc.Platform{
				{
					Id:        "some:platform",
					Installed: "2.6.4",
				},
				{
					Id: "some:otherplatform",
				},
			},
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().UpdateIndex(gomock.Any(), platIndexReq, gomock.Any())
		env.ArduinoCli.EXPECT().PlatformSearch(searchReq).Return(expectedResp, nil)

		resp, err := env.CliWrapper.SearchPlatforms()
		assert.NoError(st, err)
		assert.Equal(st, expectedResp.SearchOutput, resp)
	})

	runCliTest("returns all supported boards", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}

		req := &rpc.PlatformListRequest{
			Instance:      inst,
			UpdatableOnly: false,
			All:           true,
		}

		resp := []*rpc.Platform{
			{
				Boards: []*rpc.Board{
					{
						Name: "Some board name",
						Fqbn: "some:fqbn",
					},
				},
			},
		}

		expectBoards := []*cli.BoardWithPort{
			{
				FQBN: resp[0].Boards[0].Fqbn,
				Name: resp[0].Boards[0].Name,
			},
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().GetPlatforms(req).Return(resp, nil)

		boards := env.CliWrapper.AllBoards()
		assert.Equal(st, expectBoards, boards)
	})

	runCliTest("searches libraries", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}
		query := "some query"
		req := &rpc.LibrarySearchRequest{
			Instance: inst,
			Query:    query,
		}
		resp := &rpc.LibrarySearchResponse{
			Libraries: []*rpc.SearchedLibrary{
				{
					Name:   "Some library",
					Latest: &rpc.LibraryRelease{Version: "3.3.3"},
				},
			},
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().LibrarySearch(gomock.Any(), req).Return(resp, nil)

		libs, err := env.CliWrapper.SearchLibraries(query)
		assert.NoError(st, err)
		assert.Equal(st, resp.Libraries, libs)
	})

	runCliTest("installs libraries", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}
		lib := "somelib"
		version := "6.6.6"

		req := &rpc.LibraryInstallRequest{
			Instance: inst,
			Name:     lib,
			Version:  version,
		}

		listReq := &rpc.LibraryListRequest{
			Instance: inst,
		}

		listResp := &rpc.LibraryListResponse{
			InstalledLibraries: []*rpc.InstalledLibrary{
				{
					Library: &rpc.Library{
						Name:    lib,
						Version: version,
					},
				},
			},
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().LibraryInstall(gomock.Any(), req, gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().LibraryList(gomock.Any(), listReq).Return(listResp, nil)

		installedVers, err := env.CliWrapper.InstallLibrary(lib, version)
		assert.NoError(st, err)
		assert.Equal(st, version, installedVers)
	})

	runCliTest("uninstalls libraries", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}
		lib := "somelib"
		req := &rpc.LibraryUninstallRequest{
			Instance: inst,
			Name:     lib,
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().LibraryUninstall(gomock.Any(), req, gomock.Any())

		err := env.CliWrapper.UninstallLibrary(lib)
		assert.NoError(st, err)
	})

	runCliTest("returns installed libraries", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}
		req := &rpc.LibraryListRequest{
			Instance: inst,
		}
		resp := &rpc.LibraryListResponse{
			InstalledLibraries: []*rpc.InstalledLibrary{
				{
					Library: &rpc.Library{
						Name:    "somelib",
						Version: "1.2.3",
					},
				},
			},
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().LibraryList(gomock.Any(), req).Return(resp, nil)

		libs, err := env.CliWrapper.GetInstalledLibs()
		assert.NoError(st, err)
		assert.Equal(st, resp.InstalledLibraries, libs)
	})

	runCliTest("returns client version", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}
		version := "1.8.7"

		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().Version().Return(version)

		v := env.CliWrapper.ClientVersion()
		assert.Equal(st, v, version)
	})
}
