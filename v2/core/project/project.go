package project

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

	ardijson "github.com/robgonnella/ardi/v2/core/ardi-json"
	ardiyaml "github.com/robgonnella/ardi/v2/core/ardi-yaml"
	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/types"
	"github.com/robgonnella/ardi/v2/util"
	log "github.com/sirupsen/logrus"
)

// Project represents an arduino project
type Project struct {
	Sketch    string
	Directory string
	Baud      int
	Client    *rpc.Client
	ardiJSON  *ardijson.ArdiJSON
	ardiYAML  *ardiyaml.ArdiYAML
	logger    *log.Logger
}

// New returns new Project instance
func New(logger *log.Logger) (*Project, error) {
	client, err := rpc.NewClient(logger)
	if err != nil {
		return nil, err
	}

	ardiJSON, err := ardijson.New(logger)
	if err != nil {
		logger.WithError(err).Error()
		return nil, err
	}

	ardiYAML, err := ardiyaml.New(logger)
	if err != nil {
		logger.WithError(err).Error()
		return nil, err
	}

	return &Project{
		Client:   client,
		ardiJSON: ardiJSON,
		ardiYAML: ardiYAML,
		logger:   logger,
	}, nil
}

// Init initializes directory as an ardi project
func Init(logger *log.Logger) error {
	dataDir := paths.ArdiProjectDataDir
	confPath := paths.ArdiProjectDataConfig
	if err := util.InitDataDirectory(dataDir, confPath); err != nil {
		logger.WithError(err).Error()
		return err
	}
	logger.Info("data directory initialized")
	if err := initializeArdiJSON(); err != nil {
		logger.WithError(err).Error()
		return err
	}
	logger.Info("ardi.json initialized")
	return nil
}

// ProcessSketch to find directory, filepath, and baud
func (p *Project) ProcessSketch(sketchDir string) error {
	if sketchDir == "" {
		msg := "Must provide a sketch directory as an argument"
		err := errors.New("Missing directory argument")
		p.logger.WithError(err).Error(msg)
		return err
	}

	stat, err := os.Stat(sketchDir)
	if err != nil {
		p.logger.WithError(err).Error()
		return err
	}

	mode := stat.Mode()
	if mode.IsRegular() {
		sketchDir = path.Dir(sketchDir)
	}

	fmt.Printf("SKETCH DIRECTORY: %s\n", sketchDir)

	sketchFile, err := findSketch(sketchDir, p.logger)
	if err != nil {
		return err
	}

	sketchBaud := parseSketchBaud(sketchFile, p.logger)
	if sketchBaud != 0 {
		fmt.Println("")
		p.logger.WithField("detected baud", sketchBaud).Info("Detected baud rate from sketch file.")
		fmt.Println("")
	}

	p.Sketch = sketchFile
	p.Directory = sketchDir
	p.Baud = sketchBaud
	return nil
}

// ListBuilds specified in ardi.json
func (p *Project) ListBuilds(builds []string) {
	p.ardiJSON.ListBuilds(builds)
}

// ListLibraries specified in ardi.json
func (p *Project) ListLibraries() {
	p.ardiJSON.ListLibraries()
}

// AddBuild to ardi.json build specifications
func (p *Project) AddBuild(name, platform, boardURL, path, fqbn string, buildProps []string) {
	if platform != "" {
		p.Client.InstallPlatform(platform)
	}
	if boardURL != "" {
		p.ardiYAML.AddBoardURL(boardURL)
	}
	p.ardiJSON.AddBuild(name, platform, boardURL, path, fqbn, buildProps)
}

// RemoveBuild removes specified build(s) from project
func (p *Project) RemoveBuild(builds []string) {
	for _, build := range builds {
		p.ardiJSON.RemoveBuild(build)
	}
}

// Build specified project from ardi.json, or build all projects if left blank
func (p *Project) Build(builds []string) error {
	if len(builds) > 0 {
		for _, name := range builds {
			if build, ok := p.ardiJSON.Config.Builds[name]; ok {
				fmt.Printf("Build PATH: %s\n", build.Path)
				if err := p.ProcessSketch(build.Path); err != nil {
					p.logger.WithError(err).Error()
					return err
				}
				if build.Platform != "" {
					p.Client.InstallPlatform(build.Platform)
				}
				if build.BoardURL != "" {
					p.ardiYAML.AddBoardURL(build.BoardURL)
				}
				buildProps := []string{}
				for prop, instruction := range build.Props {
					buildProps = append(buildProps, fmt.Sprintf("%s=%s", prop, instruction))
				}
				p.logger.Infof("Building %s", build)
				if err := p.Client.Compile(build.FQBN, p.Directory, p.Sketch, name, buildProps, false); err != nil {
					p.logger.WithError(err).Errorf("Build failed for %s", build)
					return err
				}
			} else {
				p.logger.Warnf("No build specification for %s", build)
			}
		}
		return nil
	}
	// Build all
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
		if err := p.Client.Compile(build.FQBN, p.Directory, p.Sketch, name, buildProps, false); err != nil {
			p.logger.WithError(err).Errorf("Build faild for %s", name)
			return err
		}
	}
	return nil
}

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

// private methods
func (p *Project) isQuiet() bool {
	return p.logger.Level == log.InfoLevel
}

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
