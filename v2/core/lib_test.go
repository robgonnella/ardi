package core_test

import (
	"errors"
	"fmt"
	"testing"

	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/golang/mock/gomock"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestLibCore(t *testing.T) {
	testutil.RunUnitTest("installs versioned library", t, func(env *testutil.UnitTestEnv) {
		lib := "Adafruit_Pixie"
		version := "1.0.0"
		library := fmt.Sprintf("%s@%s", lib, version)
		installedVersion := "1.0.0-alpha.2"

		instance := &rpc.Instance{Id: int32(1)}
		req := &rpc.LibraryInstallRequest{
			Instance: instance,
			Name:     lib,
			Version:  version,
		}
		listReq := &rpc.LibraryListRequest{
			Instance: instance,
		}
		listResp := &rpc.LibraryListResponse{
			InstalledLibraries: []*rpc.InstalledLibrary{
				{
					Library: &rpc.Library{Name: lib, Version: installedVersion},
				},
			},
		}
		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().UpdateLibrariesIndex(gomock.Any(), gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().LibraryInstall(gomock.Any(), req, gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().LibraryList(gomock.Any(), listReq).Return(listResp, nil)

		returnedLib, returnedVers, err := env.ArdiCore.Lib.Add(library)
		assert.NoError(env.T, err)
		assert.Equal(env.T, returnedLib, lib)
		assert.Equal(env.T, returnedVers, installedVersion)
	})

	testutil.RunUnitTest("returns install error", t, func(env *testutil.UnitTestEnv) {
		errString := "dummy error"
		dummyErr := errors.New(errString)

		lib := "Adafruit_Pixie"
		version := "1.0.0"
		library := fmt.Sprintf("%s@%s", lib, version)

		instance := &rpc.Instance{Id: int32(1)}
		req := &rpc.LibraryInstallRequest{
			Instance: instance,
			Name:     lib,
			Version:  version,
		}
		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().UpdateLibrariesIndex(gomock.Any(), gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().LibraryInstall(gomock.Any(), req, gomock.Any(), gomock.Any()).Return(dummyErr)

		_, _, err := env.ArdiCore.Lib.Add(library)
		assert.Error(env.T, err)
		assert.EqualError(env.T, err, errString)
	})

	testutil.RunUnitTest("uninstalls library", t, func(env *testutil.UnitTestEnv) {
		libName := "Adafruit_Pixie"
		instance := &rpc.Instance{Id: int32(1)}
		req := &rpc.LibraryUninstallRequest{
			Instance: instance,
			Name:     libName,
		}
		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().LibraryUninstall(gomock.Any(), req, gomock.Any()).Return(nil)
		err := env.ArdiCore.Lib.Remove(libName)
		assert.NoError(env.T, err)
	})

	testutil.RunUnitTest("returns uninstall error", t, func(env *testutil.UnitTestEnv) {
		errString := "dummy error"
		dummyErr := errors.New(errString)
		libName := "Adafruit_Pixie"
		instance := &rpc.Instance{Id: int32(1)}
		req := &rpc.LibraryUninstallRequest{
			Instance: instance,
			Name:     libName,
		}
		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().LibraryUninstall(gomock.Any(), req, gomock.Any()).Return(dummyErr)
		err := env.ArdiCore.Lib.Remove(libName)
		assert.Error(env.T, err)
		assert.EqualError(env.T, err, errString)
	})

	testutil.RunUnitTest("prints library searches to stdout", t, func(env *testutil.UnitTestEnv) {
		searchQuery := "wifi101"

		latest := rpc.LibraryRelease{Version: "1.2.1"}

		libReleaseMap := map[string]*rpc.LibraryRelease{
			"1.2.1": &latest,
		}

		lib := rpc.SearchedLibrary{
			Name:     "WIFI101",
			Latest:   &latest,
			Releases: libReleaseMap,
		}

		searchedLibs := []*rpc.SearchedLibrary{&lib}
		instance := &rpc.Instance{Id: int32(1)}
		req := &rpc.LibrarySearchRequest{
			Instance: instance,
			Query:    searchQuery,
		}
		resp := &rpc.LibrarySearchResponse{
			Libraries: searchedLibs,
		}
		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().UpdateLibrariesIndex(gomock.Any(), gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().LibrarySearch(gomock.Any(), req).Return(resp, nil)

		err := env.ArdiCore.Lib.Search(searchQuery)
		assert.NoError(env.T, err)

		assert.Contains(env.T, env.Stdout.String(), lib.Name)
	})

	testutil.RunUnitTest("prints installed libraries to stdout", t, func(env *testutil.UnitTestEnv) {
		installedLib := rpc.InstalledLibrary{
			Library: &rpc.Library{
				Name:     "My favorite library",
				Version:  "1.2.2",
				Sentence: "This is my favoritest library",
			},
		}

		instance := &rpc.Instance{Id: int32(1)}
		req := &rpc.LibraryListRequest{
			Instance: instance,
		}
		resp := &rpc.LibraryListResponse{
			InstalledLibraries: []*rpc.InstalledLibrary{
				{
					Library: &rpc.Library{
						Name:     installedLib.Library.Name,
						Version:  installedLib.Library.Version,
						Sentence: installedLib.Library.Sentence,
					},
				},
			},
		}
		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.ArduinoCli.EXPECT().LibraryList(gomock.Any(), req).Return(resp, nil)

		env.ArdiCore.Lib.ListInstalled()
		stdout := env.Stdout.String()
		assert.Contains(env.T, stdout, installedLib.Library.Name)
		assert.Contains(env.T, stdout, installedLib.Library.Version)
		assert.Contains(env.T, stdout, installedLib.Library.Sentence)
	})
}
