package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/types"
	log "github.com/sirupsen/logrus"
)

// ArdiJSON represents core module for ardi.json manipulation
type ArdiJSON struct {
	Config types.ArdiConfig
	logger *log.Logger
}

// NewArdiJSON returns core json module for handling ardi.json config
func NewArdiJSON(logger *log.Logger) (*ArdiJSON, error) {
	config := types.ArdiConfig{}
	buildConfig, err := ioutil.ReadFile(paths.ArdiProjectBuildConfig)
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
func (a *ArdiJSON) AddBuild(name, platform, boardURL, path, fqbn string, buildProps []string) error {
	newBuild := types.ArdiBuildJSON{
		Platform: platform,
		BoardURL: boardURL,
		Path:     path,
		FQBN:     fqbn,
		Props:    make(map[string]string),
	}

	for _, p := range buildProps {
		parts := strings.SplitN(p, "=", 2)
		label := parts[0]
		instruction := parts[1]
		newBuild.Props[label] = instruction
	}

	a.logger.Infof("Addding build: %s", name)
	a.printBuild(name, newBuild)
	a.Config.Builds[name] = newBuild
	return a.write()
}

// RemoveBuild removes specified build from ardi.json
func (a *ArdiJSON) RemoveBuild(build string) error {
	delete(a.Config.Builds, build)
	return a.write()
}

// ListBuilds lists build specifications in ardi.json
func (a *ArdiJSON) ListBuilds(builds []string) {
	fmt.Println("")
	if len(builds) > 0 {
		for _, name := range builds {
			if b, ok := a.Config.Builds[name]; ok {
				a.printBuild(name, b)
			}
		}
	}
	for name, build := range a.Config.Builds {
		a.printBuild(name, build)
	}
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
	newData, err := json.MarshalIndent(a.Config, "", "  ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(paths.ArdiProjectBuildConfig, newData, 0644); err != nil {
		return err
	}

	return nil
}

// private
func (a *ArdiJSON) printBuild(name string, b types.ArdiBuildJSON) {
	a.logger.Println("")
	a.logger.Printf("%s:\n", name)
	a.logger.Printf("  Platform: %s\n", b.Platform)
	a.logger.Printf("  Board URL: %s\n", b.BoardURL)
	a.logger.Printf("  Path: %s\n", b.Path)
	a.logger.Printf("  FQBN: %s\n", b.FQBN)
	a.logger.Printf("  Props:\n")
	for prop, instruction := range b.Props {
		a.logger.Printf("    %s: %s\n", prop, instruction)
	}
	a.logger.Println("")
}
