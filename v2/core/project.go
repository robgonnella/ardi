package core

import (
	"errors"
	"fmt"

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

var errInitialization = errors.New("project not initialized")

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
		return err
	}

	p.logger.Info("data directory initialized")

	if err := util.InitArdiJSON(); err != nil {
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
			return err
		}
		p.ardiJSON = ardiJSON
	}

	if p.ardiYAML == nil {
		ardiYAML, err := NewArdiYAML()
		if err != nil {
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

// AddLibrary adds a library to ardi.json
func (p *ProjectCore) AddLibrary(name, version string) error {
	if p.ardiJSON == nil {
		return errInitialization
	}
	return p.ardiJSON.AddLibrary(name, version)
}

// RemoveLibrary removes a library from ardi.json
func (p *ProjectCore) RemoveLibrary(name string) error {
	if p.ardiJSON == nil {
		return errInitialization
	}
	return p.ardiJSON.RemoveLibrary(name)
}

// ListLibraries specified in ardi.json
func (p *ProjectCore) ListLibraries() {
	if p.ardiJSON == nil {
		return
	}
	p.ardiJSON.ListLibraries()
}

// GetLibraries returns map of libraries specified in ardi.json
func (p *ProjectCore) GetLibraries() map[string]string {
	if p.ardiJSON == nil {
		return make(map[string]string)
	}
	return p.ardiJSON.Config.Libraries
}

// AddPlatform adds a platform to ardi.json
func (p *ProjectCore) AddPlatform(platform, vers string) error {
	if p.ardiJSON == nil {
		return errInitialization
	}
	return p.ardiJSON.AddPlatform(platform, vers)
}

// RemovePlatform removes a platform from ardi.json
func (p *ProjectCore) RemovePlatform(platform string) error {
	if p.ardiJSON == nil {
		return errInitialization
	}
	return p.ardiJSON.RemovePlatform(platform)
}

// ListPlatforms lists all project platforms in config file
func (p *ProjectCore) ListPlatforms() {
	if p.ardiJSON == nil {
		return
	}
	p.ardiJSON.ListPlatforms()
}

// GetPlatforms returns map of platforms specified in ardi.json
func (p *ProjectCore) GetPlatforms() map[string]string {
	if p.ardiJSON == nil {
		return make(map[string]string)
	}
	return p.ardiJSON.Config.Platforms
}

// AddBuild to ardi.json build specifications
func (p *ProjectCore) AddBuild(name, platform, boardURL, path, fqbn string, buildProps []string) error {
	if p.ardiJSON == nil || p.ardiYAML == nil {
		return errInitialization
	}
	if platform != "" {
		installed, vers, err := p.client.InstallPlatform(platform)
		if err != nil {
			return err
		}
		if err := p.ardiJSON.AddPlatform(installed, vers); err != nil {
			return err
		}
	}
	if boardURL != "" {
		if err := p.ardiYAML.AddBoardURL(boardURL); err != nil {
			return err
		}
		if err := p.ardiJSON.AddBoardURL(boardURL); err != nil {
			return err
		}
	}
	return p.ardiJSON.AddBuild(name, path, fqbn, buildProps)
}

// RemoveBuild removes specified build from project
func (p *ProjectCore) RemoveBuild(build string) error {
	if p.ardiJSON == nil {
		return errInitialization
	}
	return p.ardiJSON.RemoveBuild(build)
}

// GetBuilds returns map of builds stored in ardi.json
func (p *ProjectCore) GetBuilds() map[string]types.ArdiBuildJSON {
	if p.ardiJSON == nil {
		return make(map[string]types.ArdiBuildJSON)
	}
	return p.ardiJSON.Config.Builds
}

// ListBuilds specified in ardi.json
func (p *ProjectCore) ListBuilds(builds []string) {
	if p.ardiJSON == nil {
		return
	}
	p.ardiJSON.ListBuilds(builds)
}

// AddBoardURL add board url to project config
func (p *ProjectCore) AddBoardURL(url string) error {
	if p.ardiJSON == nil || p.ardiYAML == nil {
		return errInitialization
	}
	if err := p.ardiYAML.AddBoardURL(url); err != nil {
		return err
	}
	return p.ardiJSON.AddBoardURL(url)
}

// RemoveBoardURL removes board url from project config
func (p *ProjectCore) RemoveBoardURL(url string) error {
	if p.ardiJSON == nil || p.ardiYAML == nil {
		return errInitialization
	}
	if err := p.ardiYAML.RemoveBoardURL(url); err != nil {
		return err
	}
	return p.ardiJSON.RemoveBoardURL(url)
}

// ListBoardURLS lists board urls specified in ardi.json
func (p *ProjectCore) ListBoardURLS() {
	if p.ardiJSON == nil || p.ardiYAML == nil {
		return
	}
	p.ardiJSON.ListBoardURLS()
}

// GetBoardURLS returns the list of board urls in config file
func (p *ProjectCore) GetBoardURLS() []string {
	if p.ardiJSON == nil {
		return []string{}
	}
	return p.ardiJSON.Config.BoardURLS
}

// Build builds only the build name specified by the user
func (p *ProjectCore) Build(buildName string) error {
	if p.ardiJSON == nil || p.ardiYAML == nil {
		return errInitialization
	}

	if buildName == "" {
		return errors.New("Empty build list")
	}

	build, ok := p.ardiJSON.Config.Builds[buildName]

	if !ok {
		return fmt.Errorf("No build specification for %s", buildName)
	}
	if err := p.ProcessSketch(build.Path); err != nil {
		return err
	}

	buildProps := []string{}
	for prop, instruction := range build.Props {
		buildProps = append(buildProps, fmt.Sprintf("%s=%s", prop, instruction))
	}

	p.logger.Infof("Building %s", buildName)
	opts := rpc.CompileOpts{
		FQBN:       build.FQBN,
		SketchDir:  p.Directory,
		SketchPath: p.Sketch,
		ExportName: buildName,
		BuildProps: buildProps,
		ShowProps:  false,
	}
	if err := p.client.Compile(opts); err != nil {
		return err
	}

	return nil
}

// BuildAll builds all builds specified in config
func (p *ProjectCore) BuildAll() error {
	if p.ardiJSON == nil || p.ardiYAML == nil {
		return errInitialization
	}

	if len(p.ardiJSON.Config.Builds) == 0 {
		p.logger.Warn("No builds defined. Use 'ardi project add build' to define a build.")
		return nil
	}
	for buildName, build := range p.ardiJSON.Config.Builds {
		if err := p.ProcessSketch(build.Path); err != nil {
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
			ExportName: buildName,
			BuildProps: buildProps,
			ShowProps:  false,
		}
		if err := p.client.Compile(opts); err != nil {
			return err
		}
	}
	return nil
}

// GetDataConfig returns config file contents in .ardi/arduino-cli.yml
func (p *ProjectCore) GetDataConfig() types.DataConfig {
	if p.ardiYAML == nil {
		return types.DataConfig{}
	}
	return p.ardiYAML.Config
}

// private methods
func (p *ProjectCore) isQuiet() bool {
	return p.logger.Level == log.InfoLevel
}
