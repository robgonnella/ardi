package core

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/types"
	"github.com/robgonnella/ardi/v2/util"
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
func (a *ArdiJSON) AddBuild(name, path, fqbn string, buildProps []string) error {
	newBuild := types.ArdiBuildJSON{
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
	a.logger.Println("")
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
	delete(a.Config.Libraries, name)
	return a.write()
}

// ListLibraries lists installed libraries
func (a *ArdiJSON) ListLibraries() {
	a.logger.Println("")
	for library, version := range a.Config.Libraries {
		a.logger.Printf("%s: %s\n", library, version)
	}
	a.logger.Println("")
}

// AddPlatform to ardi.json
func (a *ArdiJSON) AddPlatform(platform, version string) error {
	a.Config.Platforms[platform] = version
	return a.write()
}

// RemovePlatform from ardi.json
func (a *ArdiJSON) RemovePlatform(platform string) error {
	delete(a.Config.Platforms, platform)
	return a.write()
}

// ListPlatforms lists all board urls in config
func (a *ArdiJSON) ListPlatforms() {
	a.logger.Println("")
	for platform, vers := range a.Config.Platforms {
		a.logger.Infof("%s: %s", platform, vers)
	}
	a.logger.Println("")
}

// AddBoardURL to ardi.json
func (a *ArdiJSON) AddBoardURL(url string) error {
	if !util.ArrayContains(a.Config.BoardURLS, url) {
		a.logger.Infof("Adding board url: %s", url)
		a.Config.BoardURLS = append(a.Config.BoardURLS, url)
		return a.write()
	}
	a.logger.Infof("board url already added: %s", url)
	return nil
}

// RemoveBoardURL from ardi.json
func (a *ArdiJSON) RemoveBoardURL(url string) error {
	if util.ArrayContains(a.Config.BoardURLS, url) {
		newList := []string{}
		for _, u := range a.Config.BoardURLS {
			if u != url {
				newList = append(newList, u)
			}
		}
		a.Config.BoardURLS = newList
		return a.write()
	}
	a.logger.Infof("board url not in config: %s", url)
	return nil
}

// ListBoardURLS lists all board urls in config
func (a *ArdiJSON) ListBoardURLS() {
	a.logger.Info("Board URLS")
	for _, url := range a.Config.BoardURLS {
		a.logger.Info(url)
	}
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
	a.logger.Printf("  Path: %s\n", b.Path)
	a.logger.Printf("  FQBN: %s\n", b.FQBN)
	a.logger.Printf("  Props:\n")
	for prop, instruction := range b.Props {
		a.logger.Printf("    %s: %s\n", prop, instruction)
	}
	a.logger.Println("")
}
