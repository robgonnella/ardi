package core_test

import (
	"fmt"
	"testing"

	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestLibCore(t *testing.T) {
	testutil.RunTest("installs versioned library", t, func(st *testing.T, env testutil.TestEnv) {
		defer env.Ctrl.Finish()

		lib := "Adafruit_Pixie"
		version := "1.0.0"
		library := fmt.Sprintf("%s@%s", lib, version)
		installedVersion := "1.0.0-alpha.2"

		env.Client.EXPECT().InstallLibrary(lib, version).Times(1).Return(installedVersion, nil)

		returnedLib, returnedVers, err := env.ArdiCore.Lib.Add(library)
		assert.NoError(t, err)
		assert.Equal(t, returnedLib, lib)
		assert.Equal(t, returnedVers, installedVersion)
	})

	testutil.RunTest("uninstalls library", t, func(st *testing.T, env testutil.TestEnv) {
		defer env.Ctrl.Finish()

		libName := "Adafruit_Pixie"
		version := "1.0.0"
		libWithVers := fmt.Sprintf("%s@%s", libName, version)

		installedVersion := "1.0.0-alpha.2"

		env.Client.EXPECT().InstallLibrary(libName, version).Times(1).Return(installedVersion, nil)

		returnedLib, returnedVers, err := env.ArdiCore.Lib.Add(libWithVers)
		assert.NoError(t, err)
		assert.Equal(t, returnedLib, libName)
		assert.Equal(t, returnedVers, installedVersion)

		env.Client.EXPECT().UninstallLibrary(libName).Times(1).Return(nil)

		err = env.ArdiCore.Lib.Remove(libName)
		assert.NoError(t, err)
	})
}
