package ardijson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/types"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// BuildJSON represents the build properties in ardi.json
type BuildJSON struct {
	Path  string            `json:"path"`
	FQBN  string            `json:"fqbn"`
	Props map[string]string `json:"props"`
}

// Config represents the ardi.json file
type Config struct {
	Libraries map[string]string    `json:"libraries"`
	Builds    map[string]BuildJSON `json:"builds"`
}

// ArdiJSON represents core module for ardi.json manipulation
type ArdiJSON struct {
	Config Config
	logger *log.Logger
}

// New returns core json module for handling ardi.json config
func New(logger *log.Logger) (*ArdiJSON, error) {
	initConfigFiles(logger)
	config := Config{}
	buildConfig, err := ioutil.ReadFile(paths.ArdiBuildConfig)
	if err != nil {
		logger.WithError(err).Error("Failed to read ardi.json")
		return nil, err
	}
	if err := json.Unmarshal(buildConfig, &config); err != nil {
		logger.WithError(err).Error("Failed to parse ardi.json")
		return nil, err
	}
	return &ArdiJSON{
		Config: config,
		logger: logger,
	}, nil
}

// AddBuild to ardi.json
func (a *ArdiJSON) AddBuild(name, path, fqbn string, buildProps []string) error {
	newBuild := BuildJSON{
		Path:  path,
		FQBN:  fqbn,
		Props: make(map[string]string),
	}
	for _, p := range buildProps {
		parts := strings.SplitN(p, "=", 2)
		label := parts[0]
		instruction := parts[1]
		newBuild.Props[label] = instruction
	}

	a.Config.Builds[name] = newBuild
	return a.write()
}

// RemoveBuild removes specified build from ardi.json
func (a *ArdiJSON) RemoveBuild(build string) error {
	delete(a.Config.Builds, build)
	return a.write()
}

// ListBuilds lists build specifications in ardi.json
func (a *ArdiJSON) ListBuilds() {
	fmt.Println("")
	for name, build := range a.Config.Builds {
		fmt.Printf("%s:\n", name)
		fmt.Printf("\tPath: %s\n", build.Path)
		fmt.Printf("\tFQBN: %s\n", build.FQBN)
		fmt.Printf("\tProps:\n")
		for prop, instruction := range build.Props {
			fmt.Printf("\t\t%s: %s\n", prop, instruction)
		}
	}
	fmt.Println("")
}

// AddLibrary to ardi.json
func (a *ArdiJSON) AddLibrary(name, version string) error {
	a.Config.Libraries[name] = version
	return a.write()
}

// RemoveLibrary from ardi.json
func (a *ArdiJSON) RemoveLibrary(name string) error {
	for lib := range a.Config.Libraries {
		if lib == name {
			delete(a.Config.Libraries, lib)
		}
	}
	return a.write()
}

// ListLibraries lists installed libraries
func (a *ArdiJSON) ListLibraries() {
	fmt.Println("")
	for library, version := range a.Config.Libraries {
		fmt.Printf("%s: %s\n", library, version)
	}
	fmt.Println("")
}

func (a *ArdiJSON) write() error {
	newData, err := json.MarshalIndent(a.Config, "", " ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(paths.ArdiBuildConfig, newData, 0644); err != nil {
		return err
	}

	return nil
}

// private helper
func initConfigFiles(logger *log.Logger) {
	if _, err := os.Stat(paths.ArdiDataConfig); os.IsNotExist(err) {
		dataConfig := types.DataConfig{
			ProxyType:      "auto",
			SketchbookPath: ".",
			ArduinoData:    ".",
			BoardManager:   make(map[string]interface{}),
		}
		yamlConfig, _ := yaml.Marshal(&dataConfig)
		ioutil.WriteFile(paths.ArdiDataConfig, yamlConfig, 0644)
		logger.Info("ardi.yaml initialized")
	}
	if _, err := os.Stat(paths.ArdiBuildConfig); os.IsNotExist(err) {
		buildConfig := Config{
			Libraries: make(map[string]string),
			Builds:    make(map[string]BuildJSON),
		}
		jsonConfig, _ := json.MarshalIndent(&buildConfig, "\n", " ")
		ioutil.WriteFile(paths.ArdiBuildConfig, jsonConfig, 0644)
		logger.Info("ardi.json initialized")
	}
}
