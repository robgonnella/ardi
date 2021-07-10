package cli_test

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"path/filepath"
	"testing"

	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/golang/mock/gomock"
	"github.com/robgonnella/ardi/v2/cli-wrapper"
	"github.com/robgonnella/ardi/v2/mocks"
	"github.com/robgonnella/ardi/v2/util"
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
		svrSettings := util.GenArduinoCliSettings(".")
		ctx := context.Background()

		mockArduinoCli.EXPECT().InitSettings(settingsPath)

		cliWrapper := cli.NewCli(ctx, settingsPath, svrSettings, logger, mockArduinoCli)

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

	runCliTest("upgrades platforms", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}
		pkg := "something"
		arch := "validish"
		platform := fmt.Sprintf("%s:%s", pkg, arch)
		req := &rpc.PlatformUpgradeRequest{
			Instance:        inst,
			PlatformPackage: pkg,
			Architecture:    arch,
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().PlatformUpgrade(gomock.Any(), req, gomock.Any(), gomock.Any())

		err := env.CliWrapper.UpgradePlatform(platform)
		assert.NoError(st, err)
	})

	runCliTest("installs platforms", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}
		pkg := "something"
		arch := "validish"
		version := "2.3.5"
		platform := fmt.Sprintf("%s:%s", pkg, arch)
		platformWithVersion := fmt.Sprintf("%s@%s", platform, version)
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

		libIndexReq := &rpc.UpdateLibrariesIndexRequest{Instance: inst}

		searchReq := &rpc.PlatformSearchRequest{
			Instance:    inst,
			AllVersions: true,
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
		env.ArduinoCli.EXPECT().UpdateLibrariesIndex(gomock.Any(), libIndexReq, gomock.Any())
		env.ArduinoCli.EXPECT().PlatformSearch(searchReq).Return(expectedResp, nil)

		resp, err := env.CliWrapper.SearchPlatforms()
		assert.NoError(st, err)
		assert.Equal(st, expectedResp.SearchOutput, resp)
	})

	runCliTest("returns connected boards", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}

		resp := []*rpc.DetectedPort{
			{
				Address: "/dev/null",
				Boards: []*rpc.BoardListItem{
					{
						Name: "some-board-name",
						Fqbn: "some:fqbn",
					},
				},
			},
		}

		expectBoards := []*cli.BoardWithPort{
			{
				FQBN: resp[0].Boards[0].Fqbn,
				Name: resp[0].Boards[0].Name,
				Port: resp[0].Address,
			},
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().ConnectedBoards(inst.GetId()).Return(resp, nil)

		boards := env.CliWrapper.ConnectedBoards()
		assert.Equal(st, expectBoards, boards)
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

	runCliTest("uploads builds", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}
		fqbn := "some:fqbn"
		sketchDir := "."
		resolvedSketchDir, _ := filepath.Abs(sketchDir)
		device := "/dev/null"

		req := &rpc.UploadRequest{
			Instance:   inst,
			Fqbn:       fqbn,
			SketchPath: resolvedSketchDir,
			Port:       device,
			Verbose:    false,
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().Upload(gomock.Any(), req, gomock.Any(), gomock.Any())

		err := env.CliWrapper.Upload(fqbn, sketchDir, device)
		assert.NoError(st, err)
	})

	runCliTest("compiles sketches", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}
		sketchDir := "."
		resolvedSketchDir, _ := filepath.Abs(sketchDir)
		resolvedSketchPath := path.Join(resolvedSketchDir, "some_sketch.ino")
		resolvedBuildDir := path.Join(resolvedSketchDir, "build")

		opts := cli.CompileOpts{
			FQBN:       "some:fqbn",
			SketchDir:  sketchDir,
			SketchPath: "./some_sketch.ino",
		}

		req := &rpc.CompileRequest{
			Instance:        inst,
			Fqbn:            opts.FQBN,
			SketchPath:      resolvedSketchPath,
			ExportDir:       resolvedBuildDir,
			BuildProperties: opts.BuildProps,
			ShowProperties:  opts.ShowProps,
			Verbose:         false,
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), req, gomock.Any(), gomock.Any(), gomock.Any())

		err := env.CliWrapper.Compile(opts)
		assert.NoError(st, err)
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

	runCliTest("returns target board when fqbn and port specified", t, func(env cliTestEnv, st *testing.T) {
		fqbn := "some:fqbn"
		port := "/dev/null"

		expectedBoard := &cli.BoardWithPort{
			FQBN: fqbn,
			Port: port,
		}

		board, err := env.CliWrapper.GetTargetBoard(fqbn, port, false)
		assert.NoError(st, err)
		assert.Equal(st, expectedBoard, board)

		board, err = env.CliWrapper.GetTargetBoard(fqbn, port, true)
		assert.NoError(st, err)
		assert.Equal(st, expectedBoard, board)
	})

	runCliTest("returns connected board match", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}

		fqbn := "some:fqbn"

		resp := []*rpc.DetectedPort{
			{
				Address: "/dev/null",
				Boards: []*rpc.BoardListItem{
					{
						Name: "some-board-name",
						Fqbn: fqbn,
					},
				},
			},
		}

		expectedBoard := &cli.BoardWithPort{
			FQBN: resp[0].Boards[0].Fqbn,
			Name: resp[0].Boards[0].Name,
			Port: resp[0].Address,
		}

		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().ConnectedBoards(inst.GetId()).Return(resp, nil)
		env.ArduinoCli.EXPECT().GetPlatforms(gomock.Any())

		board, err := env.CliWrapper.GetTargetBoard(fqbn, "", true)
		assert.NoError(st, err)
		assert.Equal(st, expectedBoard, board)
	})

	runCliTest("returns error if no match for provided fqbn and onlyConnected=true", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}
		resp := []*rpc.DetectedPort{}
		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().ConnectedBoards(inst.GetId()).Return(resp, nil)
		env.ArduinoCli.EXPECT().GetPlatforms(gomock.Any())

		board, err := env.CliWrapper.GetTargetBoard("some:fqbn", "", true)
		assert.Error(st, err)
		assert.Nil(st, board)
	})

	runCliTest("returns board without port if fqbn provided and onlyConnected=false", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}
		resp := []*rpc.DetectedPort{}
		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().ConnectedBoards(inst.GetId()).Return(resp, nil)
		env.ArduinoCli.EXPECT().GetPlatforms(gomock.Any())

		fqbn := "some:fqbn"
		expectedBoard := &cli.BoardWithPort{FQBN: fqbn}

		board, err := env.CliWrapper.GetTargetBoard(fqbn, "", false)
		assert.NoError(st, err)
		assert.Equal(st, expectedBoard, board)
	})

	runCliTest("returns error if fqbn empty and no boards connected and onlyConnected=true", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}

		resp := []*rpc.DetectedPort{}

		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().ConnectedBoards(inst.GetId()).Return(resp, nil)
		env.ArduinoCli.EXPECT().GetPlatforms(gomock.Any())

		board, err := env.CliWrapper.GetTargetBoard("", "", true)
		assert.Error(st, err)
		assert.Nil(st, board)
	})

	runCliTest("returns error and prints list if fqbn empty and no boards connected and onlyConnected=false", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}

		installed := []*rpc.Platform{
			{
				Id:        "some:platform",
				Installed: "2.4.6",
				Boards: []*rpc.Board{
					{
						Name: "Some board name",
						Fqbn: "some:board:fqbn",
					},
				},
			},
		}
		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().ConnectedBoards(inst.GetId())
		env.ArduinoCli.EXPECT().GetPlatforms(gomock.Any()).Return(installed, nil)

		env.ClearStdout()
		board, err := env.CliWrapper.GetTargetBoard("", "", false)
		assert.Error(st, err)
		assert.Nil(st, board)

		output := env.Stdout.String()
		assert.Contains(st, output, installed[0].Boards[0].Name)
		assert.Contains(st, output, installed[0].Boards[0].Fqbn)
	})

	runCliTest("returns first connected board", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}

		resp := []*rpc.DetectedPort{
			{
				Address: "/dev/null",
				Boards: []*rpc.BoardListItem{
					{
						Name: "some-board-name",
						Fqbn: "some:board:fqbn",
					},
				},
			},
		}
		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().ConnectedBoards(inst.GetId()).Return(resp, nil).Times(2)
		env.ArduinoCli.EXPECT().GetPlatforms(gomock.Any()).Times(2)

		expectedBoard := &cli.BoardWithPort{
			FQBN: resp[0].Boards[0].Fqbn,
			Name: resp[0].Boards[0].Name,
			Port: resp[0].Address,
		}

		board, err := env.CliWrapper.GetTargetBoard("", "", false)
		assert.NoError(st, err)
		assert.Equal(st, expectedBoard, board)

		board, err = env.CliWrapper.GetTargetBoard("", "", true)
		assert.NoError(st, err)
		assert.Equal(st, expectedBoard, board)
	})

	runCliTest("prints connected boards if fqbn empty and more than one board connected", t, func(env cliTestEnv, st *testing.T) {
		inst := &rpc.Instance{Id: int32(1)}

		resp := []*rpc.DetectedPort{
			{
				Address: "/dev/null",
				Boards: []*rpc.BoardListItem{
					{
						Name: "some-board-name",
						Fqbn: "some:board:fqbn",
					},
					{
						Name: "another-board-name",
						Fqbn: "another:board:fqbn",
					},
				},
			},
		}
		env.ArduinoCli.EXPECT().CreateInstance().Return(inst).AnyTimes()
		env.ArduinoCli.EXPECT().ConnectedBoards(inst.GetId()).Return(resp, nil).Times(2)
		env.ArduinoCli.EXPECT().GetPlatforms(gomock.Any()).Times(2)

		env.ClearStdout()
		board, err := env.CliWrapper.GetTargetBoard("", "", false)
		assert.Error(st, err)
		assert.Nil(st, board)

		output := env.Stdout.String()
		assert.Contains(st, output, resp[0].Boards[0].Name)
		assert.Contains(st, output, resp[0].Boards[0].Fqbn)
		assert.Contains(st, output, resp[0].Boards[1].Name)
		assert.Contains(st, output, resp[0].Boards[1].Fqbn)

		env.ClearStdout()
		board, err = env.CliWrapper.GetTargetBoard("", "", true)
		assert.Error(st, err)
		assert.Nil(st, board)

		output = env.Stdout.String()
		assert.Contains(st, output, resp[0].Boards[0].Name)
		assert.Contains(st, output, resp[0].Boards[0].Fqbn)
		assert.Contains(st, output, resp[0].Boards[1].Name)
		assert.Contains(st, output, resp[0].Boards[1].Fqbn)
	})
}
