package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/types"
	"github.com/robgonnella/ardi/v2/util"
)

// ProjectCore represents an arduino project
type ProjectCore struct {
	Sketch    string
	Directory string
	Baud      int
	client    rpc.Client
	ardiJSON  *ArdiJSON
	ardiYAML  *ArdiYAML
	logger    *log.Logger
}

// NewProjectCore returns new Project instance
func NewProjectCore(client rpc.Client, logger *log.Logger) *ProjectCore {
	return &ProjectCore{
		client: client,
		logger: logger,
	}
}

// Init initialize ardi project core
func (p *ProjectCore) Init(port string) error {
	dataDir := paths.ArdiProjectDataDir
	confPath := paths.ArdiProjectDataConfig

	if err := util.InitDataDirectory(port, dataDir, confPath); err != nil {
		p.logger.WithError(err).Error()
		return err
	}

	p.logger.Info("data directory initialized")

	if err := initializeArdiJSON(); err != nil {
		p.logger.WithError(err).Error()
		return err
	}

	p.logger.Info("ardi.json initialized")

	return nil
}

// SetConfigHelpers sets ardi.json and data-dir helpers
func (p *ProjectCore) SetConfigHelpers() error {
	if p.ardiJSON == nil {
		ardiJSON, err := NewArdiJSON(p.logger)
		if err != nil {
			p.logger.WithError(err).Error()
			return err
		}
		p.ardiJSON = ardiJSON
	}

	if p.ardiYAML == nil {
		ardiYAML, err := NewArdiYAML(p.logger)
		if err != nil {
			p.logger.WithError(err).Error()
			return err
		}
		p.ardiYAML = ardiYAML
	}

	return nil
}

// ProcessSketch to find directory, filepath, and baud
func (p *ProjectCore) ProcessSketch(sketchDir string) error {
	sketchDir, sketchFile, sketchBaud, err := util.ProcessSketch(sketchDir, p.logger)
	if err != nil {
		return err
	}

	p.Sketch = sketchFile
	p.Directory = sketchDir
	p.Baud = sketchBaud
	return nil
}

// ListBuilds specified in ardi.json
func (p *ProjectCore) ListBuilds(builds []string) {
	p.ardiJSON.ListBuilds(builds)
}

// AddLibrary adds a library to ardi.json
func (p *ProjectCore) AddLibrary(name, version string) error {
	return p.ardiJSON.AddLibrary(name, version)
}

// RemoveLibrary removes a library from ardi.json
func (p *ProjectCore) RemoveLibrary(name string) error {
	return p.ardiJSON.RemoveLibrary(name)
}

// ListLibraries specified in ardi.json
func (p *ProjectCore) ListLibraries() {
	p.ardiJSON.ListLibraries()
}

// AddBuild to ardi.json build specifications
func (p *ProjectCore) AddBuild(name, platform, boardURL, path, fqbn string, buildProps []string) {
	if platform != "" {
		p.client.InstallPlatform(platform)
	}
	if boardURL != "" {
		p.ardiYAML.AddBoardURL(boardURL)
	}
	p.ardiJSON.AddBuild(name, platform, boardURL, path, fqbn, buildProps)
}

// RemoveBuild removes specified build(s) from project
func (p *ProjectCore) RemoveBuild(builds []string) {
	for _, build := range builds {
		p.ardiJSON.RemoveBuild(build)
	}
}

// BuildList builds only the build-names specified by the user
func (p *ProjectCore) BuildList(builds []string) error {
	if len(builds) == 0 {
		err := errors.New("Empty build list")
		p.logger.WithError(err).Error("Cannot build")
		return err
	}
	for _, name := range builds {
		build, ok := p.ardiJSON.Config.Builds[name]
		if !ok {
			p.logger.Warnf("No build specification for %s", name)
			continue
		}
		if err := p.ProcessSketch(build.Path); err != nil {
			p.logger.WithError(err).Error()
			return err
		}
		if build.Platform != "" {
			p.client.InstallPlatform(build.Platform)
		}
		if build.BoardURL != "" {
			p.ardiYAML.AddBoardURL(build.BoardURL)
		}
		buildProps := []string{}
		for prop, instruction := range build.Props {
			buildProps = append(buildProps, fmt.Sprintf("%s=%s", prop, instruction))
		}

		p.logger.Infof("Building %s", build)
		opts := rpc.CompileOpts{
			FQBN:       build.FQBN,
			SketchDir:  p.Directory,
			SketchPath: p.Sketch,
			ExportName: name,
			BuildProps: buildProps,
			ShowProps:  false,
		}
		if err := p.client.Compile(opts); err != nil {
			p.logger.WithError(err).Errorf("Build failed for %s", build)
			return err
		}
	}
	return nil
}

// BuildAll builds all builds specified in config
func (p *ProjectCore) BuildAll() error {
	if len(p.ardiJSON.Config.Builds) == 0 {
		p.logger.Warn("No builds defined. Use \"ardi project add build\" to define a build.")
		return nil
	}
	for name, build := range p.ardiJSON.Config.Builds {
		if err := p.ProcessSketch(build.Path); err != nil {
			p.logger.WithError(err).Error()
			return err
		}
		buildProps := []string{}
		for prop, instruction := range build.Props {
			buildProps = append(buildProps, fmt.Sprintf("%s=%s", prop, instruction))
		}

		p.logger.Infof("Building %s", build.Path)
		opts := rpc.CompileOpts{
			FQBN:       build.FQBN,
			SketchDir:  p.Directory,
			SketchPath: p.Sketch,
			ExportName: name,
			BuildProps: buildProps,
			ShowProps:  false,
		}
		if err := p.client.Compile(opts); err != nil {
			p.logger.WithError(err).Errorf("Build faild for %s", name)
			return err
		}
	}
	return nil
}

// GetLibraries returns list of libraries specified in ardi.json
func (p *ProjectCore) GetLibraries() []string {
	libs := []string{}
	for name, vers := range p.ardiJSON.Config.Libraries {
		libs = append(libs, fmt.Sprintf("%s@%s", name, vers))
	}
	return libs
}

// private methods
func (p *ProjectCore) isQuiet() bool {
	return p.logger.Level == log.InfoLevel
}

// helpers
func initializeArdiJSON() error {
	if _, err := os.Stat(paths.ArdiProjectBuildConfig); os.IsNotExist(err) {
		buildConfig := types.ArdiConfig{
			Libraries: make(map[string]string),
			Builds:    make(map[string]types.ArdiBuildJSON),
		}
		jsonConfig, _ := json.MarshalIndent(&buildConfig, "\n", " ")
		if err := ioutil.WriteFile(paths.ArdiProjectBuildConfig, jsonConfig, 0644); err != nil {
			return err
		}
	}
	return nil
}
