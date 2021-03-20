package util_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/robgonnella/ardi/v2/util"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func writeSettings(conf string, data []byte) error {
	os.RemoveAll(conf)
	return ioutil.WriteFile(conf, data, 0644)
}

func TestUtilArrayContains(t *testing.T) {
	t.Run("returns true if array contains item", func(st *testing.T) {
		item := "someitem"
		array := []string{item}
		assert.True(st, util.ArrayContains(array, item))
	})

	t.Run("returns false if array does not contain item", func(st *testing.T) {
		item := "someitem"
		array := []string{"someotheritem"}
		assert.False(st, util.ArrayContains(array, item))
	})
}

func TestUtilCreateDataDirector(t *testing.T) {
	t.Run("creates data directory", func(st *testing.T) {
		dir := "./test-data-dir"
		os.RemoveAll(dir)
		err := util.CreateDataDir(dir)
		assert.NoError(st, err)
		assert.DirExists(st, dir)
		os.RemoveAll(dir)
	})

	t.Run("does not error if directory exists", func(st *testing.T) {
		dir := "./another-data-dir"
		os.RemoveAll(dir)
		err := util.CreateDataDir(dir)
		assert.NoError(st, err)
		assert.DirExists(st, dir)

		err = util.CreateDataDir(dir)
		assert.NoError(st, err)
		assert.DirExists(st, dir)
		os.RemoveAll(dir)
	})

	t.Run("deletes data directory", func(st *testing.T) {
		dir := "somefancydatadirectory"
		os.RemoveAll(dir)
		err := util.CreateDataDir(dir)
		assert.NoError(st, err)
		assert.DirExists(st, dir)

		err = util.CleanDataDirectory(dir)
		assert.NoError(st, err)
		_, fileErr := os.Stat(dir)
		assert.True(st, os.IsNotExist(fileErr))
	})
}

func TestUtilArduinoCliSettings(t *testing.T) {
	t.Run("errors if file does not exist", func(st *testing.T) {
		data, err := util.ReadArduinoCliSettings("./noop")
		assert.Error(st, err)
		assert.Nil(st, data)
	})

	t.Run("errors if file malformed", func(st *testing.T) {
		conf := "testconf"
		data := []byte("noop\ndoublenoop")
		err := writeSettings(conf, data)
		assert.NoError(st, err)

		settings, err := util.ReadArduinoCliSettings(conf)
		assert.Error(st, err)
		assert.Nil(st, settings)
		os.RemoveAll(conf)
	})

	t.Run("returns settings from file", func(st *testing.T) {
		conf := "success-conf"
		level := "debug"
		dataDir := "."
		expected := util.GenArduinoCliSettings(level, dataDir)
		assert.Equal(st, expected.Logging.Level, level)
		assert.Equal(st, expected.Directories.Data, dataDir)
		assert.Equal(st, expected.Directories.Downloads, path.Join(dataDir, "staging"))
		assert.Equal(st, expected.Directories.User, path.Join(dataDir, "Arduino"))

		byteData, err := yaml.Marshal(expected)
		assert.NoError(st, err)
		err = writeSettings(conf, byteData)
		assert.NoError(st, err)

		data, err := util.ReadArduinoCliSettings(conf)
		assert.NoError(st, err)
		assert.Equal(st, expected, data)
		os.RemoveAll(conf)
	})

	t.Run("returns project path", func(st *testing.T) {
		opts := util.GetAllSettingsOpts{
			Global: false,
		}
		settingsPath := util.GetCliSettingsPath(opts)
		assert.Equal(st, paths.ArduinoCliProjectConfig, settingsPath)
	})

	t.Run("returns global path", func(st *testing.T) {
		opts := util.GetAllSettingsOpts{
			Global: true,
		}
		settingsPath := util.GetCliSettingsPath(opts)
		assert.Equal(st, paths.ArduinoCliGlobalConfig, settingsPath)
	})
}

func TestUtilArdiConfig(t *testing.T) {
	t.Run("errors if file does not exist", func(st *testing.T) {
		data, err := util.ReadArdiConfig("./noop")
		assert.Error(st, err)
		assert.Nil(st, data)
	})

	t.Run("errors if file malformed", func(st *testing.T) {
		conf := "testconf"
		data := []byte("noop\ndoublenoop")
		err := writeSettings(conf, data)
		assert.NoError(st, err)

		settings, err := util.ReadArdiConfig(conf)
		assert.Error(st, err)
		assert.Nil(st, settings)
		os.RemoveAll(conf)
	})

	t.Run("returns settings from file", func(st *testing.T) {
		conf := "ardi-success-conf"
		level := "debug"
		expected := util.GenArdiConfig(level)
		assert.Equal(st, expected.Daemon.LogLevel, level)

		byteData, err := json.Marshal(expected)
		assert.NoError(st, err)
		err = writeSettings(conf, byteData)
		assert.NoError(st, err)

		data, err := util.ReadArdiConfig(conf)
		assert.NoError(st, err)
		assert.Equal(st, expected, data)
		os.RemoveAll(conf)
	})
}

func TestUtilGetAllSettings(t *testing.T) {
	t.Run("returns default settings if project files not found", func(st *testing.T) {
		dataDir := paths.ArdiProjectDataDir
		level := "fancy-log-level"
		os.RemoveAll(dataDir)

		expectedConfig := util.GenArdiConfig(level)
		expectedSettings := util.GenArduinoCliSettings(level, dataDir)

		opts := util.GetAllSettingsOpts{
			Global:   false,
			LogLevel: level,
		}
		config, settings := util.GetAllSettings(opts)
		assert.Equal(st, expectedConfig, config)
		assert.Equal(st, expectedSettings, settings)
	})

	t.Run("returns settings from project files", func(st *testing.T) {
		dataDir := paths.ArdiProjectDataDir
		level := "fancy-log-level"
		expectedConfig := util.GenArdiConfig(level)
		expectedSettings := util.GenArduinoCliSettings(level, dataDir)

		os.RemoveAll(dataDir)

		writeOpts := util.WriteSettingsOpts{
			Global:             false,
			ArdiConfig:         expectedConfig,
			ArduinoCliSettings: expectedSettings,
		}
		util.WriteAllSettings(writeOpts)

		assert.DirExists(st, dataDir)
		assert.FileExists(st, paths.ArdiProjectConfig)
		assert.FileExists(st, paths.ArduinoCliProjectConfig)

		opts := util.GetAllSettingsOpts{
			Global:   false,
			LogLevel: level,
		}

		config, settings := util.GetAllSettings(opts)
		assert.Equal(st, expectedConfig, config)
		assert.Equal(st, expectedSettings, settings)

		os.RemoveAll(dataDir)
		os.RemoveAll(paths.ArdiProjectConfig)
	})

	t.Run("returns default settings if global files not found", func(st *testing.T) {
		dataDir := paths.ArdiGlobalDataDir
		level := "fancy-log-level"
		os.RemoveAll(dataDir)

		expectedConfig := util.GenArdiConfig(level)
		expectedSettings := util.GenArduinoCliSettings(level, dataDir)

		opts := util.GetAllSettingsOpts{
			Global:   true,
			LogLevel: level,
		}
		config, settings := util.GetAllSettings(opts)
		assert.Equal(st, expectedConfig, config)
		assert.Equal(st, expectedSettings, settings)
	})

	t.Run("returns settings from global files", func(st *testing.T) {
		dataDir := paths.ArdiGlobalDataDir
		level := "fancy-log-level"
		expectedConfig := util.GenArdiConfig(level)
		expectedSettings := util.GenArduinoCliSettings(level, dataDir)

		os.RemoveAll(dataDir)

		writeOpts := util.WriteSettingsOpts{
			Global:             true,
			ArdiConfig:         expectedConfig,
			ArduinoCliSettings: expectedSettings,
		}
		err := util.WriteAllSettings(writeOpts)
		assert.NoError(st, err)
		assert.DirExists(st, dataDir)
		assert.FileExists(st, paths.ArdiGlobalConfig)
		assert.FileExists(st, paths.ArduinoCliGlobalConfig)

		opts := util.GetAllSettingsOpts{
			Global:   true,
			LogLevel: level,
		}
		config, settings := util.GetAllSettings(opts)
		assert.Equal(st, expectedConfig, config)
		assert.Equal(st, expectedSettings, settings)

		os.RemoveAll(dataDir)
	})
}

func TestUtilInitProjectDirectory(t *testing.T) {
	t.Run("initialized directory with project config files", func(st *testing.T) {
		os.RemoveAll(paths.ArdiProjectDataDir)
		os.RemoveAll(paths.ArdiProjectConfig)

		err := util.InitProjectDirectory()
		assert.NoError(st, err)
		assert.DirExists(st, paths.ArdiProjectDataDir)
		assert.FileExists(st, paths.ArdiProjectConfig)

		os.RemoveAll(paths.ArdiProjectDataDir)
		os.RemoveAll(paths.ArdiProjectConfig)
	})
}

func TestUtilIsProjectDirectory(t *testing.T) {
	t.Run("returns false if project ardi.json not found", func(st *testing.T) {
		os.RemoveAll(paths.ArdiProjectConfig)
		assert.False(st, util.IsProjectDirectory())
	})

	t.Run("returns true if project ardi.json found", func(st *testing.T) {
		os.RemoveAll(paths.ArdiProjectConfig)
		file, _ := os.Create(paths.ArdiProjectConfig)
		defer func() {
			file.Close()
			os.RemoveAll(paths.ArdiProjectConfig)
		}()
		assert.True(st, util.IsProjectDirectory())
	})
}

func TestUtilGetDaemonLogLevel(t *testing.T) {
	t.Run("returns debug for logger level", func(st *testing.T) {
		logger := logrus.New()
		logger.SetLevel(logrus.DebugLevel)
		level := util.GetDaemonLogLevel(logger)
		assert.Equal(st, "debug", level)
	})

	t.Run("returns fata for other logger levels", func(st *testing.T) {
		logger := logrus.New()
		logger.SetLevel(logrus.InfoLevel)
		level := util.GetDaemonLogLevel(logger)
		assert.Equal(st, "fatal", level)
	})
}

func TestUtilGeneratePropsArray(t *testing.T) {
	t.Run("generates props array from props object", func(st *testing.T) {
		propsObject := make(map[string]string)
		prop := "someprop"
		propValue := "somevalue"
		propsObject[prop] = propValue
		propsArray := util.GeneratePropsArray(propsObject)
		assert.Equal(st, len(propsArray), 1)
		assert.Contains(st, propsArray, fmt.Sprintf("%s=%s", prop, propValue))
	})
}

func TestUtilGeneratePropsMap(t *testing.T) {
	t.Run("generates props map from props array", func(st *testing.T) {
		prop := "someprop"
		propValue := "somevalue"
		propsArray := []string{fmt.Sprintf("%s=%s", prop, propValue)}
		propsMap := util.GeneratePropsMap(propsArray)
		assert.Equal(st, propsMap[prop], propValue)
	})
}

func TestUtilProcessSketch(t *testing.T) {
	t.Run("errors if sketch param empty", func(st *testing.T) {
		project, err := util.ProcessSketch("")
		assert.Error(st, err)
		assert.Nil(st, project)
	})

	t.Run("errors if path does not contain sketch", func(st *testing.T) {
		project, err := util.ProcessSketch(".")
		assert.Error(st, err)
		assert.Nil(st, project)
	})

	t.Run("returns project for valid sketch directory", func(st *testing.T) {
		dir := testutil.BlinkProjectDir()
		sketch := path.Join(dir, "blink.ino")

		project, err := util.ProcessSketch(dir)
		assert.NoError(st, err)
		assert.Equal(st, project.Directory, dir)
		assert.Equal(st, project.Sketch, sketch)
		assert.Equal(st, project.Baud, 9600)
	})

	t.Run("returns project for valid sketch file", func(st *testing.T) {
		dir := testutil.BlinkProjectDir()
		sketch := path.Join(dir, "blink.ino")
		project, err := util.ProcessSketch(sketch)
		assert.NoError(st, err)
		assert.Equal(st, project.Directory, dir)
		assert.Equal(st, project.Sketch, sketch)
		assert.Equal(st, project.Baud, 9600)
	})

	t.Run("reads baud from file", func(st *testing.T) {
		dir := testutil.Blink14400ProjectDir()
		sketch := path.Join(dir, "blink14400.ino")
		project, err := util.ProcessSketch(sketch)
		assert.NoError(st, err)
		assert.Equal(st, project.Directory, dir)
		assert.Equal(st, project.Sketch, sketch)
		assert.Equal(st, project.Baud, 14400)
	})
}
