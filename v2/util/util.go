package util

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/types"
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
func GenDefaultDataConfig(port, dataDirPath string) types.DataConfig {
	return types.DataConfig{
		BoardManager: types.BoardManager{AdditionalUrls: []string{}},
		Daemon: types.Daemon{
			Port: port,
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
func InitDataDirectory(port, dataDirPath, dataConfigPath string) error {
	if _, err := os.Stat(dataDirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDirPath, 0777); err != nil {
			return err
		}
	}

	if _, err := os.Stat(dataConfigPath); os.IsNotExist(err) {
		dataConfig := GenDefaultDataConfig(port, dataDirPath)
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

// ProcessSketch looks for .ino file in specified directory and parses
func ProcessSketch(sketchDir string, logger *log.Logger) (string, string, int, error) {
	if sketchDir == "" {
		msg := "Must provide a sketch directory as an argument"
		err := errors.New("Missing directory argument")
		logger.WithError(err).Error(msg)
		return "", "", 0, err
	}

	stat, err := os.Stat(sketchDir)
	if err != nil {
		logger.WithError(err).Error()
		return "", "", 0, err
	}

	mode := stat.Mode()
	if mode.IsRegular() {
		sketchDir = path.Dir(sketchDir)
	}

	sketchFile, err := findSketch(sketchDir, logger)
	if err != nil {
		return "", "", 0, err
	}

	sketchBaud := parseSketchBaud(sketchFile, logger)
	if sketchBaud != 0 {
		fmt.Println("")
		logger.WithField("detected baud", sketchBaud).Info("Detected baud rate from sketch file.")
		fmt.Println("")
	}

	return sketchDir, sketchFile, sketchBaud, nil
}

// private helpers
// helpers
func findSketch(directory string, logger *log.Logger) (string, error) {
	sketchFile := ""

	d, err := os.Open(directory)
	if err != nil {
		logger.WithError(err).Error("Failed to open sketch directory")
		return "", err
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		logger.WithError(err).Error("Cannot process .ino file")
		return "", err
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ".ino" {
				sketchFile = path.Join(directory, file.Name())
			}
		}
	}
	if sketchFile == "" {
		msg := fmt.Sprintf("Failed to find .ino file in %s", directory)
		logger.Error(msg)
		return "", errors.New(msg)
	}

	if sketchFile, err = filepath.Abs(sketchFile); err != nil {
		msg := "Could not resolve sketch file path"
		logger.WithError(err).Error(msg)
		return "", errors.New(msg)
	}

	return sketchFile, nil
}

func parseSketchBaud(sketch string, logger *log.Logger) int {
	var baud int
	rgx := regexp.MustCompile(`Serial\.begin\((\d+)\);`)
	file, err := os.Open(sketch)
	if err != nil {
		// Log the error and return 0 for baud to let script continue
		// with either default value or value specified from command-line.
		logger.WithError(err).
			WithField("sketch", sketch).
			Info("Failed to read sketch")
		return baud
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := scanner.Text()
		if match := rgx.MatchString(text); match {
			stringBaud := strings.TrimSpace(rgx.ReplaceAllString(text, "$1"))
			if baud, err = strconv.Atoi(stringBaud); err != nil {
				// set baud to 0 and let script continue with either default
				// value or value specified from command-line.
				logger.WithError(err).Info("Failed to parse baud rate from sketch")
				baud = 0
			}
			break
		}
	}

	return baud
}
