package cli_test

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/arduino/arduino-cli/cli/output"
	"github.com/arduino/arduino-cli/configuration"
	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/robgonnella/ardi/v3/cli-wrapper"
	"github.com/robgonnella/ardi/v3/types"
	"github.com/robgonnella/ardi/v3/util"
	"github.com/stretchr/testify/assert"
)

func TestArduinoCli(t *testing.T) {
	here, _ := filepath.Abs(".")
	dataDir := path.Join(here, ".ardi")

	ctx := context.Background()
	settingsPath := util.GetCliSettingsPath()

	tearDown := func() {
		os.RemoveAll(".ardi")
		os.RemoveAll("ardi.json")
	}

	writeSettings := func() *types.ArduinoCliSettings {
		util.InitProjectDirectory()
		config := util.GenArdiConfig()
		settings := util.GenArduinoCliSettings(dataDir)
		util.WriteAllSettings(config, settings)
		return settings
	}

	setUp := func(tt *testing.T) *cli.ArduinoCli {
		writeSettings()
		tt.Cleanup(tearDown)
		arduinoCli := cli.NewArduinoCli()
		arduinoCli.InitSettings(settingsPath)
		return arduinoCli
	}

	updatePlatformIndex := func(arduinoCli *cli.ArduinoCli) {
		rpcInstance := arduinoCli.CreateInstance()
		indexReq := &rpc.UpdateIndexRequest{Instance: rpcInstance}
		arduinoCli.UpdateIndex(ctx, indexReq, output.ProgressBar())
	}

	updateLibraryIndex := func(arduinoCli *cli.ArduinoCli) {
		rpcInstance := arduinoCli.CreateInstance()
		indexReq := &rpc.UpdateLibrariesIndexRequest{Instance: rpcInstance}
		arduinoCli.UpdateLibrariesIndex(ctx, indexReq, output.ProgressBar())
	}

	tearDown()
	t.Cleanup(tearDown)

	t.Run("return new arduino-cli instance", func(st *testing.T) {
		arudinoCli := cli.NewArduinoCli()
		assert.NotNil(st, arudinoCli)
	})

	t.Run("initializes settings", func(st *testing.T) {
		expected := writeSettings()
		st.Cleanup(tearDown)

		arduinoCli := cli.NewArduinoCli()
		arduinoCli.InitSettings(settingsPath)

		actual := configuration.Settings

		assert.Equal(st, actual.Get("directories.data"), expected.Directories.Data)
		assert.Equal(st, actual.Get("directories.downloads"), expected.Directories.Downloads)
		assert.Equal(st, actual.Get("directories.user"), expected.Directories.User)
	})

	t.Run("creates rpc instance", func(st *testing.T) {
		arduinoCli := setUp(st)

		rpcInstance := arduinoCli.CreateInstance()
		assert.NotNil(st, rpcInstance)
		assert.NotNil(st, rpcInstance.GetId())
	})

	t.Run("updates platform index", func(st *testing.T) {
		arduinoCli := setUp(st)

		rpcInst := arduinoCli.CreateInstance()
		downloadCb := output.ProgressBar()
		req := &rpc.UpdateIndexRequest{Instance: rpcInst}

		err := arduinoCli.UpdateIndex(ctx, req, downloadCb)
		assert.NoError(st, err)
		assert.FileExists(st, ".ardi/package_index.json")

	})

	t.Run("updates library index", func(st *testing.T) {
		arduinoCli := setUp(st)

		rpcInst := arduinoCli.CreateInstance()
		downloadCb := output.ProgressBar()
		req := &rpc.UpdateLibrariesIndexRequest{Instance: rpcInst}

		err := arduinoCli.UpdateLibrariesIndex(ctx, req, downloadCb)
		assert.NoError(st, err)
		assert.FileExists(st, ".ardi/library_index.json")
	})

	t.Run("installs, lists, and uninstalls platforms", func(st *testing.T) {
		arduinoCli := setUp(st)
		updatePlatformIndex(arduinoCli)

		pkg := "arduino"
		arch := "avr"
		platformID := pkg + ":" + arch
		version := "1.8.2"

		rpcInst := arduinoCli.CreateInstance()

		/**
		 * Install
		**/
		installReq := &rpc.PlatformInstallRequest{
			Instance:        rpcInst,
			PlatformPackage: pkg,
			Architecture:    arch,
			Version:         version,
		}
		downloadCb := output.ProgressBar()
		taskCb := output.TaskProgress()

		installResp, err := arduinoCli.PlatformInstall(ctx, installReq, downloadCb, taskCb)
		assert.NoError(st, err)
		assert.NotNil(st, installResp)
		assert.DirExists(st, ".ardi/packages/arduino/hardware/avr/1.8.2")

		/**
		 * List
		**/
		listReq := &rpc.PlatformListRequest{
			Instance: rpcInst,
		}
		platforms, err := arduinoCli.GetPlatforms(listReq)
		assert.NoError(st, err)
		assert.Equal(st, len(platforms), 1)
		assert.Equal(st, platforms[0].GetId(), platformID)
		assert.Equal(st, platforms[0].Installed, "1.8.2")

		/**
		 * Uninstall
		**/
		uninstallReq := &rpc.PlatformUninstallRequest{
			Instance:        rpcInst,
			PlatformPackage: pkg,
			Architecture:    arch,
		}
		taskCb = output.TaskProgress()
		_, err = arduinoCli.PlatformUninstall(ctx, uninstallReq, taskCb)
		assert.NoError(st, err)

		platforms, _ = arduinoCli.GetPlatforms(listReq)
		assert.Equal(st, len(platforms), 0)
	})

	t.Run("searches platforms", func(st *testing.T) {
		arduinoCli := setUp(st)
		updatePlatformIndex(arduinoCli)
		rpcInst := arduinoCli.CreateInstance()

		req := &rpc.PlatformSearchRequest{
			Instance:    rpcInst,
			AllVersions: false,
			SearchArgs:  "arduino:avr",
		}

		resp, err := arduinoCli.PlatformSearch(req)
		assert.NoError(st, err)
		assert.NotEmpty(st, resp.SearchOutput)
	})

	t.Run("searches libraries", func(st *testing.T) {
		arduinoCli := setUp(st)
		updateLibraryIndex(arduinoCli)

		rpcInst := arduinoCli.CreateInstance()

		req := &rpc.LibrarySearchRequest{
			Instance: rpcInst,
			Query:    "WiFi",
		}

		resp, err := arduinoCli.LibrarySearch(ctx, req)
		assert.NoError(st, err)
		assert.NotEmpty(st, resp.Libraries)
	})

	t.Run("installs, lists, and removes library", func(st *testing.T) {
		arduinoCli := setUp(st)
		updateLibraryIndex(arduinoCli)

		rpcInst := arduinoCli.CreateInstance()

		lib := "Adafruit Pixie"

		/**
		 * Install
		**/
		installReq := &rpc.LibraryInstallRequest{
			Instance: rpcInst,
			Name:     lib,
		}
		taskCb := output.TaskProgress()
		downloadCb := output.ProgressBar()
		err := arduinoCli.LibraryInstall(ctx, installReq, downloadCb, taskCb)
		assert.NoError(st, err)
		assert.DirExists(st, ".ardi/Arduino/libraries/Adafruit_Pixie")

		/**
		 * List Installed
		**/
		listReq := &rpc.LibraryListRequest{
			Instance: rpcInst,
		}
		libs, err := arduinoCli.LibraryList(ctx, listReq)
		assert.NoError(st, err)
		assert.Equal(st, len(libs.InstalledLibraries), 1)
		assert.Equal(st, libs.InstalledLibraries[0].Library.Name, lib)

		/**
		 * Remove
		**/
		uninstallReq := &rpc.LibraryUninstallRequest{
			Instance: rpcInst,
			Name:     lib,
		}
		taskCb = output.TaskProgress()
		err = arduinoCli.LibraryUninstall(ctx, uninstallReq, taskCb)
		assert.NoError(st, err)

		libs, _ = arduinoCli.LibraryList(ctx, listReq)
		assert.Equal(st, len(libs.InstalledLibraries), 0)
	})

	t.Run("prints arduino-cli version", func(st *testing.T) {
		arduinoCli := setUp(st)
		vers := arduinoCli.Version()
		assert.NotEmpty(st, vers)
	})
}
