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
	"github.com/google/uuid"
	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/types"
	"gopkg.in/yaml.v2"
)

var appID = uuid.New().String()
var appSecret = uuid.New().String()

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
func GenArduinoCliSettings(dataDir string) *types.ArduinoCliSettings {
	return &types.ArduinoCliSettings{
		BoardManager: types.BoardManager{AdditionalUrls: []string{}},
		Daemon: types.Daemon{
			Port: "",
		},
		Directories: types.Directories{
			Data:      dataDir,
			Downloads: path.Join(dataDir, "staging"),
			User:      path.Join(dataDir, "Arduino"),
		},
		Installation: types.Installation{
			ID:     appID,
			Secret: appSecret,
		},
		Library: types.Library{
			EnableUnsafeInstall: false,
		},
		Logging: types.Logging{
			Level:  "fatal",
			Format: "text",
			File:   "",
		},
		Metrics: types.Metrics{
			Addr:    ":9090",
			Enabled: false,
		},
	}
}

// GenArdiConfig returns default ardi.json in current directory
func GenArdiConfig() *types.ArdiConfig {
	return &types.ArdiConfig{
		Platforms: make(map[string]string),
		BoardURLS: []string{},
		Libraries: make(map[string]string),
		Builds:    make(map[string]types.ArdiBuild),
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
func GetAllSettings() (*types.ArdiConfig, *types.ArduinoCliSettings) {
	var ardiConfig *types.ArdiConfig
	var cliSettings *types.ArduinoCliSettings

	dataDir := paths.ArdiProjectDataDir
	ardiConf := paths.ArdiProjectConfig
	cliConf := paths.ArduinoCliProjectConfig

	if _, err := os.Stat(ardiConf); os.IsNotExist(err) {
		ardiConfig = GenArdiConfig()
	} else if ardiConfig, err = ReadArdiConfig(ardiConf); err != nil {
		ardiConfig = GenArdiConfig()
	}

	if _, err := os.Stat(cliConf); os.IsNotExist(err) {
		cliSettings = GenArduinoCliSettings(dataDir)
	} else if cliSettings, err = ReadArduinoCliSettings(cliConf); err != nil {
		cliSettings = GenArduinoCliSettings(dataDir)
	}

	return ardiConfig, cliSettings
}

// GetCliSettingsPath returns path to arduino-cli.yaml based on scope
func GetCliSettingsPath() string {
	cliConf := paths.ArduinoCliProjectConfig
	return cliConf
}

// WriteAllSettings writes all settings files
func WriteAllSettings(ardiConfig *types.ArdiConfig, arduinoSettings *types.ArduinoCliSettings) error {
	dataDir := paths.ArdiProjectDataDir
	ardiConf := paths.ArdiProjectConfig
	cliConf := paths.ArduinoCliProjectConfig

	if err := CreateDataDir(dataDir); err != nil {
		return err
	}

	byteData, _ := json.MarshalIndent(ardiConfig, "", "\t")
	if err := ioutil.WriteFile(ardiConf, byteData, 0644); err != nil {
		return err
	}

	byteData, _ = yaml.Marshal(arduinoSettings)
	if err := ioutil.WriteFile(cliConf, byteData, 0644); err != nil {
		return err
	}

	if _, fileErr := os.Stat(path.Join(dataDir, "inventory.yaml")); os.IsNotExist(fileErr) {
		inventory.Init(dataDir)
	}

	return nil
}

// InitProjectDirectory initializes a directory as an ardi project
func InitProjectDirectory() error {
	ardiConfig, cliSettings := GetAllSettings()
	return WriteAllSettings(ardiConfig, cliSettings)
}

// IsProjectDirectory returns whether or not currect directory has been initialized as an ardi project
func IsProjectDirectory() bool {
	_, buildErr := os.Stat(paths.ArdiProjectConfig)
	return !os.IsNotExist(buildErr)
}

// CleanDataDirectory removes directory
func CleanDataDirectory(dir string) error {
	return os.RemoveAll(dir)
}

// GeneratePropsMap returns map of build props from string array
func GeneratePropsMap(buildProps []string) map[string]string {
	props := make(map[string]string)

	for _, p := range buildProps {
		parts := strings.SplitN(p, "=", 2)
		label := parts[0]
		instruction := parts[1]
		props[label] = instruction
	}

	return props
}

// GeneratePropsArray returns an arrary of props from props map
func GeneratePropsArray(props map[string]string) []string {
	buildProps := []string{}
	for prop, instruction := range props {
		buildProps = append(buildProps, fmt.Sprintf("%s=%s", prop, instruction))
	}
	return buildProps
}

// ProcessSketch looks for .ino file in specified directory and parses
func ProcessSketch(filePath string) (*types.Project, error) {
	if filePath == "" {
		err := errors.New("missing sketch argument")
		return nil, err
	}

	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	sketchDir := ""
	sketchFile := ""

	mode := stat.Mode()
	if mode.IsRegular() {
		sketchFile = filePath
		sketchDir = filepath.Dir(sketchFile)
	} else {
		sketchDir = filePath
		if sketchFile, err = findSketch(sketchDir); err != nil {
			return nil, err
		}
	}

	sketchBaud := ParseSketchBaud(sketchFile)

	return &types.Project{
		Directory: sketchDir,
		Sketch:    sketchFile,
		Baud:      sketchBaud,
	}, nil
}

// ParseSketchBaud reads a sketch file and tries to parse baud rate
func ParseSketchBaud(sketch string) int {
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

// private helpers
// helpers
func findSketch(directory string) (string, error) {
	stat, err := os.Stat(directory)
	if err != nil {
		return "", err
	}

	if !stat.IsDir() {
		return "", errors.New("not a directory")
	}

	sketchFile := ""
	absPath, _ := filepath.Abs(directory)
	searchName := fmt.Sprintf("%s.ino", filepath.Base(absPath))

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
		if file.Mode().IsRegular() && file.Name() == searchName {
			sketchFile = path.Join(directory, file.Name())
		}
	}
	if sketchFile == "" {
		msg := fmt.Sprintf("Failed to find %s in %s", searchName, directory)
		return "", errors.New(msg)
	}

	return sketchFile, nil
}
