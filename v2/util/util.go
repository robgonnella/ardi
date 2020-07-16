package util

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/arduino/arduino-cli/inventory"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/types"
)

// DefaultDaemonPort default port to run arduino-cli daemon
const DefaultDaemonPort = "50051"

// DefaultDaemonLogLevel default arduino-cli daemon log level
const DefaultDaemonLogLevel = "fatal"

// GetAllSettingsOpts options for retrieving all settings
type GetAllSettingsOpts struct {
	Global   bool
	LogLevel string
	Port     string
}

// WriteSettingsOpts options for writing all settings to file
type WriteSettingsOpts struct {
	Global             bool
	ArdiConfig         *types.ArdiConfig
	ArduinoCliSettings *types.ArduinoCliSettings
}

// ArrayContains checks if a string array contains a value
func ArrayContains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// CreateDataDir creates a data dir with proper permissions for ardi / arduino-cli
func CreateDataDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return err
		}
	}
	return nil
}

// ReadArduinoCliSettings reads data config file and returns unmarshalled data and stringified version
func ReadArduinoCliSettings(confPath string) (*types.ArduinoCliSettings, error) {
	var config types.ArduinoCliSettings
	dataFile, err := os.Open(confPath)
	if err != nil {
		return nil, err
	}
	defer dataFile.Close()

	byteData, err := ioutil.ReadAll(dataFile)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(byteData, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// GenArduinoCliSettings generated data config file with default values
func GenArduinoCliSettings(logLevel, port, dataDir string) *types.ArduinoCliSettings {
	return &types.ArduinoCliSettings{
		BoardManager: types.BoardManager{AdditionalUrls: []string{}},
		Daemon: types.Daemon{
			Port: port,
		},
		Directories: types.Directories{
			Data:      dataDir,
			Downloads: path.Join(dataDir, "staging"),
			User:      path.Join(dataDir, "Arduino"),
		},
		Logging: types.Logging{
			Level:  logLevel,
			Format: "text",
			File:   "",
		},
		Telemetry: types.Telemetry{
			Addr:    ":9090",
			Enabled: false,
		},
	}
}

// GenArdiConfig returns default ardi.json in current directory
func GenArdiConfig(logLevel, port string) *types.ArdiConfig {
	return &types.ArdiConfig{
		Daemon: types.ArdiDaemonConfig{
			Port:     port,
			LogLevel: logLevel,
		},
		Platforms: make(map[string]string),
		BoardURLS: []string{},
		Libraries: make(map[string]string),
		Builds:    make(map[string]types.ArdiBuildJSON),
	}
}

// ReadArdiConfig reads ardi.json and returns config
func ReadArdiConfig(confPath string) (*types.ArdiConfig, error) {
	var config types.ArdiConfig
	configFile, err := os.Open(confPath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	byteData, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(byteData, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// GetAllSettings returns settings for both ardi and arduino-cli
func GetAllSettings(opts GetAllSettingsOpts) (*types.ArdiConfig, *types.ArduinoCliSettings) {
	var ardiConfig *types.ArdiConfig
	var cliSettings *types.ArduinoCliSettings
	logLevel := opts.LogLevel
	port := opts.Port

	dataDir := paths.ArdiProjectDataDir
	ardiConf := paths.ArdiProjectConfig
	cliConf := paths.ArduinoCliProjectConfig

	if opts.Global {
		dataDir = paths.ArdiGlobalDataDir
		ardiConf = paths.ArdiGlobalConfig
		cliConf = paths.ArduinoCliGlobalConfig
	}

	if _, err := os.Stat(ardiConf); os.IsNotExist(err) {
		ardiConfig = GenArdiConfig(logLevel, port)
	} else if ardiConfig, err = ReadArdiConfig(ardiConf); err != nil {
		ardiConfig = GenArdiConfig(logLevel, port)
	}
	if port != "" {
		ardiConfig.Daemon.Port = port
	}
	if ardiConfig.Daemon.Port == "" {
		ardiConfig.Daemon.Port = DefaultDaemonPort
	}
	if ardiConfig.Daemon.LogLevel == "" {
		ardiConfig.Daemon.LogLevel = DefaultDaemonLogLevel
	}

	if _, err := os.Stat(cliConf); os.IsNotExist(err) {
		cliSettings = GenArduinoCliSettings(ardiConfig.Daemon.LogLevel, ardiConfig.Daemon.Port, dataDir)
	} else if cliSettings, err = ReadArduinoCliSettings(cliConf); err != nil {
		cliSettings = GenArduinoCliSettings(ardiConfig.Daemon.LogLevel, ardiConfig.Daemon.Port, dataDir)
	}
	cliSettings.Daemon.Port = ardiConfig.Daemon.Port
	cliSettings.Logging.Level = ardiConfig.Daemon.LogLevel

	return ardiConfig, cliSettings
}

// WriteAllSettings writes all settings files
func WriteAllSettings(opts WriteSettingsOpts) error {
	dataDir := paths.ArdiProjectDataDir
	ardiConf := paths.ArdiProjectConfig
	cliConf := paths.ArduinoCliProjectConfig

	if opts.Global {
		dataDir = paths.ArdiGlobalDataDir
		ardiConf = paths.ArdiGlobalConfig
		cliConf = paths.ArduinoCliGlobalConfig
	}

	if err := CreateDataDir(dataDir); err != nil {
		return err
	}

	byteData, _ := json.MarshalIndent(opts.ArdiConfig, "\n", " ")
	if err := ioutil.WriteFile(ardiConf, byteData, 0644); err != nil {
		return err
	}

	byteData, _ = yaml.Marshal(opts.ArduinoCliSettings)
	if err := ioutil.WriteFile(cliConf, byteData, 0644); err != nil {
		return err
	}

	if _, fileErr := os.Stat(path.Join(dataDir, "inventory.yaml")); os.IsNotExist(fileErr) {
		inventory.Init(dataDir)
	}

	return nil
}

// InitProjectDirectory initializes a directory as an ardi project
func InitProjectDirectory(port string) error {
	getOpts := GetAllSettingsOpts{
		Global:   false,
		LogLevel: "fatal",
		Port:     port,
	}
	ardiConfig, cliSettings := GetAllSettings(getOpts)

	writeOpts := WriteSettingsOpts{
		Global:             false,
		ArdiConfig:         ardiConfig,
		ArduinoCliSettings: cliSettings,
	}
	if err := WriteAllSettings(writeOpts); err != nil {
		return err
	}

	return nil
}

// IsProjectDirectory returns whether or not currect directory has been initialized as an ardi project
func IsProjectDirectory() bool {
	_, buildErr := os.Stat(paths.ArdiProjectConfig)
	if os.IsNotExist(buildErr) {
		return false
	}
	return true
}

// GetDaemonLogLevel returns daemon log level string based on logger settings
func GetDaemonLogLevel(logger *log.Logger) string {
	if logger.GetLevel() == log.DebugLevel {
		return "debug"
	}
	return "fatal"
}

// CleanDataDirectory removes directory
func CleanDataDirectory(dir string) error {
	return os.RemoveAll(dir)
}

// ProcessSketch looks for .ino file in specified directory and parses
func ProcessSketch(sketchDir string) (*types.Project, error) {
	if sketchDir == "" {
		err := errors.New("Missing directory argument")
		return nil, err
	}

	stat, err := os.Stat(sketchDir)
	if err != nil {
		return nil, err
	}

	mode := stat.Mode()
	if mode.IsRegular() {
		sketchDir = path.Dir(sketchDir)
	}

	sketchFile, err := findSketch(sketchDir)
	if err != nil {
		return nil, err
	}

	sketchBaud := parseSketchBaud(sketchFile)

	return &types.Project{
		Directory: sketchDir,
		Sketch:    sketchFile,
		Baud:      sketchBaud,
	}, nil
}

// private helpers
// helpers
func findSketch(directory string) (string, error) {
	sketchFile := ""

	d, err := os.Open(directory)
	if err != nil {
		return "", err
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
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
		return "", errors.New(msg)
	}

	if sketchFile, err = filepath.Abs(sketchFile); err != nil {
		msg := "Could not resolve sketch file path"
		return "", errors.New(msg)
	}

	return sketchFile, nil
}

func parseSketchBaud(sketch string) int {
	var baud = 9600
	rgx := regexp.MustCompile(`Serial\.begin\((\d+)\);`)
	file, err := os.Open(sketch)
	if err != nil {
		return baud
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := scanner.Text()
		if match := rgx.MatchString(text); match {
			stringBaud := strings.TrimSpace(rgx.ReplaceAllString(text, "$1"))
			if baud, err = strconv.Atoi(stringBaud); err != nil {
				baud = 9600
			}
			break
		}
	}

	return baud
}
