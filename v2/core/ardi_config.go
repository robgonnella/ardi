package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
	"github.com/robgonnella/ardi/v2/types"
	"github.com/robgonnella/ardi/v2/util"
	log "github.com/sirupsen/logrus"
)

// ArdiConfig represents core module for ardi.json manipulation
type ArdiConfig struct {
	config   types.ArdiConfig
	confPath string
	logger   *log.Logger
	mux      sync.Mutex
}

// NewArdiConfig returns core json module for handling ardi.json config
func NewArdiConfig(confPath string, initialConfig types.ArdiConfig, logger *log.Logger) *ArdiConfig {
	return &ArdiConfig{
		config:   initialConfig,
		confPath: confPath,
		logger:   logger,
		mux:      sync.Mutex{},
	}
}

// AddBuild to ardi.json
func (a *ArdiConfig) AddBuild(name, sketch, fqbn string, buildProps []string) error {
	project, err := util.ProcessSketch(sketch)
	if err != nil {
		return err
	}

	newBuild := types.ArdiBuild{
		Directory: project.Directory,
		Sketch:    project.Sketch,
		Baud:      project.Baud,
		FQBN:      fqbn,
	}

	props := util.GeneratePropsMap(buildProps)
	newBuild.Props = props

	a.logger.Infof("Addding build: %s", name)
	a.printBuild(name, newBuild)
	a.config.Builds[name] = newBuild
	return a.write()
}

// GetCompileOpts returns appropriate compile options for an ardi build
func (a *ArdiConfig) GetCompileOpts(buildName string) (*cli.CompileOpts, error) {
	build, ok := a.config.Builds[buildName]
	if !ok {
		return nil, fmt.Errorf("no builds found for %s", buildName)
	}

	buildProps := util.GeneratePropsArray(build.Props)

	compileOpts := &cli.CompileOpts{
		FQBN:       build.FQBN,
		SketchDir:  build.Directory,
		SketchPath: build.Sketch,
		BuildProps: buildProps,
	}

	return compileOpts, nil
}

// RemoveBuild removes specified build from ardi.json
func (a *ArdiConfig) RemoveBuild(build string) error {
	delete(a.config.Builds, build)
	return a.write()
}

// ListBuilds lists build specifications in ardi.json
func (a *ArdiConfig) ListBuilds(builds []string) {
	a.logger.Println("")
	if len(builds) > 0 {
		for _, name := range builds {
			if b, ok := a.config.Builds[name]; ok {
				a.printBuild(name, b)
			}
		}
	}
	for name, build := range a.config.Builds {
		a.printBuild(name, build)
	}
}

// GetBuilds returns builds specified in config
func (a *ArdiConfig) GetBuilds() map[string]types.ArdiBuild {
	return a.config.Builds
}

// AddLibrary to ardi.json
func (a *ArdiConfig) AddLibrary(name, version string) error {
	a.config.Libraries[name] = version
	return a.write()
}

// RemoveLibrary from ardi.json
func (a *ArdiConfig) RemoveLibrary(name string) error {
	delete(a.config.Libraries, name)
	return a.write()
}

// ListLibraries lists installed libraries
func (a *ArdiConfig) ListLibraries() {
	a.logger.Println("")
	for library, version := range a.config.Libraries {
		a.logger.Printf("%s: %s\n", library, version)
	}
	a.logger.Println("")
}

// GetLibraries returns libraries specired in config
func (a *ArdiConfig) GetLibraries() map[string]string {
	return a.config.Libraries
}

// AddPlatform to ardi.json
func (a *ArdiConfig) AddPlatform(platform, version string) error {
	a.config.Platforms[platform] = version
	return a.write()
}

// RemovePlatform from ardi.json
func (a *ArdiConfig) RemovePlatform(platform string) error {
	delete(a.config.Platforms, platform)
	return a.write()
}

// ListPlatforms lists all board urls in config
func (a *ArdiConfig) ListPlatforms() {
	a.logger.Println("")
	for platform, vers := range a.config.Platforms {
		a.logger.Infof("%s: %s", platform, vers)
	}
	a.logger.Println("")
}

// GetPlatforms returns platforms specified in config
func (a *ArdiConfig) GetPlatforms() map[string]string {
	return a.config.Platforms
}

// AddBoardURL to ardi.json
func (a *ArdiConfig) AddBoardURL(url string) error {
	if !util.ArrayContains(a.config.BoardURLS, url) {
		a.logger.Infof("Adding board url: %s", url)
		a.config.BoardURLS = append(a.config.BoardURLS, url)
		return a.write()
	}
	a.logger.Infof("board url already added: %s", url)
	return nil
}

// RemoveBoardURL from ardi.json
func (a *ArdiConfig) RemoveBoardURL(url string) error {
	if util.ArrayContains(a.config.BoardURLS, url) {
		newList := []string{}
		for _, u := range a.config.BoardURLS {
			if u != url {
				newList = append(newList, u)
			}
		}
		a.config.BoardURLS = newList
		return a.write()
	}
	a.logger.Infof("board url not in config: %s", url)
	return nil
}

// ListBoardURLS lists all board urls in config
func (a *ArdiConfig) ListBoardURLS() {
	a.logger.Info("Board URLS")
	for _, url := range a.config.BoardURLS {
		a.logger.Info(url)
	}
}

// GetBoardURLS returns board urls specified in config
func (a *ArdiConfig) GetBoardURLS() []string {
	return a.config.BoardURLS
}

func (a *ArdiConfig) write() error {
	a.mux.Lock()
	defer a.mux.Unlock()

	newData, err := json.MarshalIndent(a.config, "", "\t")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(a.confPath, newData, 0644); err != nil {
		return err
	}

	return nil
}

// private
func (a *ArdiConfig) printBuild(name string, b types.ArdiBuild) {
	a.logger.Println("")
	a.logger.Printf("%s:\n", name)
	a.logger.Printf("  Directory: %s\n", b.Directory)
	a.logger.Printf("  Sketch: %s\n", b.Sketch)
	a.logger.Printf("  Baud: %d\n", b.Baud)
	a.logger.Printf("  FQBN: %s\n", b.FQBN)
	a.logger.Printf("  Props:\n")
	for prop, instruction := range b.Props {
		a.logger.Printf("    %s: %s\n", prop, instruction)
	}
	a.logger.Println("")
}
