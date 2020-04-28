package util

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/types"
	"gopkg.in/yaml.v2"
)

// ArrayContains checks if a string array contains a value
func ArrayContains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// GenDefaultDataConfig generated data config file with default values
func GenDefaultDataConfig(dataDirPath string) types.DataConfig {
	return types.DataConfig{
		BoardManager: types.BoardManager{AdditionalUrls: []string{}},
		Daemon: types.Daemon{
			Port: "50051",
		},
		Directories: types.Directories{
			Data:      dataDirPath,
			Downloads: path.Join(dataDirPath, "staging"),
			User:      path.Join(dataDirPath, "Arduino"),
		},
		Logging: types.Logging{
			Level:  "info",
			Format: "text",
			File:   "",
		},
		Telemetry: types.Telemetry{
			Addr:    ":9090",
			Enabled: false,
		},
	}
}

// IsProjectDirectory returns whether or not currect directory has been initialized as an ardi project
func IsProjectDirectory() bool {
	_, dirErr := os.Stat(paths.ArdiProjectDataDir)
	_, buildErr := os.Stat(paths.ArdiProjectBuildConfig)
	if os.IsNotExist(dirErr) && os.IsNotExist(buildErr) {
		return false
	}
	return true
}

// InitDataDirectory creates and initializes project data directory if necessary
func InitDataDirectory(dataDirPath, dataConfigPath string) error {
	if _, err := os.Stat(dataDirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDirPath, 0777); err != nil {
			return err
		}
	}

	if _, err := os.Stat(dataConfigPath); os.IsNotExist(err) {
		dataConfig := GenDefaultDataConfig(dataDirPath)
		yamlConfig, _ := yaml.Marshal(&dataConfig)
		if err := ioutil.WriteFile(dataConfigPath, yamlConfig, 0644); err != nil {
			return err
		}
	}

	return nil
}

// CleanDataDirectory removes directory
func CleanDataDirectory(dir string) error {
	return os.RemoveAll(dir)
}
