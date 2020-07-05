package core_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/arduino/arduino-cli/rpc/commands"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestLibCore(t *testing.T) {
	testutil.RunUnitTest("installs versioned library", t, func(env testutil.UnitTestEnv) {
		defer env.Ctrl.Finish()

		lib := "Adafruit_Pixie"
		version := "1.0.0"
		library := fmt.Sprintf("%s@%s", lib, version)
		installedVersion := "1.0.0-alpha.2"

		env.Client.EXPECT().InstallLibrary(lib, version).Times(1).Return(installedVersion, nil)

		returnedLib, returnedVers, err := env.ArdiCore.Lib.Add(library)
		assert.NoError(env.T, err)
		assert.Equal(env.T, returnedLib, lib)
		assert.Equal(env.T, returnedVers, installedVersion)
	})

	testutil.RunUnitTest("returns install error", t, func(env testutil.UnitTestEnv) {
		defer env.Ctrl.Finish()

		errString := "dummy error"
		dummyErr := errors.New(errString)

		lib := "Adafruit_Pixie"
		version := "1.0.0"
		library := fmt.Sprintf("%s@%s", lib, version)

		env.Client.EXPECT().InstallLibrary(lib, version).Times(1).Return("", dummyErr)

		_, _, err := env.ArdiCore.Lib.Add(library)
		assert.Error(env.T, err)
		assert.EqualError(env.T, err, errString)
	})

	testutil.RunUnitTest("uninstalls library", t, func(env testutil.UnitTestEnv) {
		defer env.Ctrl.Finish()
		libName := "Adafruit_Pixie"
		env.Client.EXPECT().UninstallLibrary(libName).Times(1).Return(nil)
		err := env.ArdiCore.Lib.Remove(libName)
		assert.NoError(env.T, err)
	})

	testutil.RunUnitTest("returns uninstall error", t, func(env testutil.UnitTestEnv) {
		defer env.Ctrl.Finish()
		errString := "dummy error"
		dummyErr := errors.New(errString)
		libName := "Adafruit_Pixie"
		env.Client.EXPECT().UninstallLibrary(libName).Times(1).Return(dummyErr)
		err := env.ArdiCore.Lib.Remove(libName)
		assert.Error(env.T, err)
		assert.EqualError(env.T, err, errString)
	})

	testutil.RunUnitTest("prints library searches to stdout", t, func(env testutil.UnitTestEnv) {
		defer env.Ctrl.Finish()

		searchQuery := "wifi101"

		latest := commands.LibraryRelease{Version: "1.2.1"}

		libReleaseMap := map[string]*commands.LibraryRelease{
			"1.2.1": &latest,
		}

		lib := commands.SearchedLibrary{
			Name:     "WIFI101",
			Latest:   &latest,
			Releases: libReleaseMap,
		}

		searchedLibs := []*commands.SearchedLibrary{&lib}

		env.Client.EXPECT().SearchLibraries(searchQuery).Times(1).Return(searchedLibs, nil)

		err := env.ArdiCore.Lib.Search(searchQuery)
		assert.NoError(env.T, err)

		assert.Contains(env.T, env.Stdout.String(), lib.Name)
	})

	testutil.RunUnitTest("prints installed libraries to stdout", t, func(env testutil.UnitTestEnv) {
		defer env.Ctrl.Finish()

		installedLib := commands.InstalledLibrary{
			Library: &commands.Library{
				Name:     "My favorite library",
				Version:  "1.2.2",
				Sentence: "This is my favoritest library",
			},
		}

		env.Client.EXPECT().GetInstalledLibs().Times(1).Return([]*commands.InstalledLibrary{&installedLib}, nil)
		env.ArdiCore.Lib.ListInstalled()

		stdout := env.Stdout.String()
		assert.Contains(env.T, stdout, installedLib.Library.Name)
		assert.Contains(env.T, stdout, installedLib.Library.Version)
		assert.Contains(env.T, stdout, installedLib.Library.Sentence)
	})
}
