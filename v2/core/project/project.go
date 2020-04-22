package project

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	ardijson "github.com/robgonnella/ardi/v2/core/ardi-json"
	"github.com/robgonnella/ardi/v2/rpc"
	log "github.com/sirupsen/logrus"
)

// Project represents an arduino project
type Project struct {
	Sketch    string
	Directory string
	Baud      int
	Client    *rpc.Client
	ardiJSON  *ardijson.ArdiJSON
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

	return &Project{
		Client:   client,
		ardiJSON: ardiJSON,
		logger:   logger,
	}, nil
}

// ProcessSketch to find directory, filepath, and baud
func (p *Project) ProcessSketch(sketchDir string) error {
	if sketchDir == "" {
		msg := "Must provide a sketch directory as an argument"
		err := errors.New("Missing directory argument")
		p.logger.WithError(err).Error(msg)
		return err
	}

	// Guard in case someone tries to pass full path to .ino file
	sketchDir = path.Dir(sketchDir)

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
func (p *Project) ListBuilds() {
	p.ardiJSON.ListBuilds()
}

// ListLibraries specified in ardi.json
func (p *Project) ListLibraries() {
	p.ardiJSON.ListLibraries()
}

// AddBuild to ardi.json build specifications
func (p *Project) AddBuild(name, path, fqbn string, buildProps []string) {
	p.ardiJSON.AddBuild(name, path, fqbn, buildProps)
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
		for _, sketch := range builds {
			for _, build := range p.ardiJSON.Config.Builds {
				if sketch == build.Path {
					buildProps := []string{}
					for prop, instruction := range build.Props {
						buildProps = append(buildProps, fmt.Sprintf("%s=%s", prop, instruction))
					}
					p.logger.Infof("Building %s", sketch)
					if err := p.Client.Compile(build.FQBN, sketch, buildProps, false); err != nil {
						p.logger.WithError(err).Errorf("Build failed for %s", sketch)
						return err
					}
					break
				}
			}
		}
		return nil
	}
	// Build all
	for _, build := range p.ardiJSON.Config.Builds {
		buildProps := []string{}
		for prop, instruction := range build.Props {
			buildProps = append(buildProps, fmt.Sprintf("%s=%s", prop, instruction))
		}
		p.logger.Infof("Building %s", build.Path)
		if err := p.Client.Compile(build.FQBN, build.Path, buildProps, false); err != nil {
			p.logger.WithError(err).Errorf("Build faild for %s", build.Path)
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
