package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/robgonnella/ardi/v2/core/rpc"
	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/types"
)

// Lib core module for lib commands
type Lib struct {
	logger *log.Logger
	RPC    *rpc.RPC
}

// BuildJSON represents the build properties in ardi.json
type BuildJSON struct {
	Props map[string]string `json:"props"`
}

// ArdiJSON represents the ardi.json file
type ArdiJSON struct {
	Libraries map[string]string    `json:"libraries"`
	Builds    map[string]BuildJSON `json:"builds"`
}

// New Lib instance
func New(dataConfig string, logger *log.Logger) (*Lib, error) {
	rpc, err := rpc.New(dataConfig, logger)
	if err != nil {
		return nil, err
	}
	return &Lib{
		logger: logger,
		RPC:    rpc,
	}, nil
}

// Init intialized current directory as an ardi project directory
func (l *Lib) Init() {
	dataConfig := types.DataConfig{
		ProxyType:      "auto",
		SketchbookPath: ".",
		ArduinoData:    ".",
		BoardManager:   make(map[string]interface{}),
	}
	yamlConfig, _ := yaml.Marshal(&dataConfig)
	ioutil.WriteFile(paths.ArdiDataConfig, yamlConfig, 0644)
	buildConfig := ArdiJSON{make(map[string]string), make(map[string]BuildJSON)}
	jsonConfig, _ := json.MarshalIndent(&buildConfig, "\n", " ")
	ioutil.WriteFile(paths.ArdiBuildConfig, jsonConfig, 0644)
	l.logger.Info("Directory initialized")
}

// Search all available libraries with optional search filter
func (l *Lib) Search(searchArg string) error {
	libraries, err := l.RPC.SearchLibraries(searchArg)
	if err != nil {
		return err
	}

	sort.Slice(libraries, func(i, j int) bool {
		return libraries[i].GetName() < libraries[j].GetName()
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "Library\tLatest\tReleases")
	for _, lib := range libraries {
		releases := []string{}
		for _, rel := range lib.GetReleases() {
			releases = append(releases, rel.GetVersion())
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", lib.GetName(), lib.GetLatest().GetVersion(), strings.Join(releases, ", "))
	}
	return nil
}

// Add library for project
func (l *Lib) Add(name, version string) error {
	if !isProjectDirectory() {
		err := errors.New("Directory not initialized. Run \"ardi lib init\" first")
		l.logger.WithError(err).Error("Failed to add library")
		return nil
	}

	l.logger.Infof("Installing library: %s %s", name, version)
	installedVersion, err := l.RPC.InstallLibrary(name, version)
	if err != nil {
		l.logger.WithError(err).Errorf("Failed to install %s", name)
		return err
	}

	config := ArdiJSON{}
	errMsg := fmt.Sprintf("Failed to add library: %s", name)

	buildConfig, err := ioutil.ReadFile(paths.ArdiBuildConfig)
	if err != nil {
		l.logger.WithError(err).Error(errMsg)
		return err
	}

	if err := json.Unmarshal(buildConfig, &config); err != nil {
		l.logger.WithError(err).Error(errMsg)
		return err
	}

	config.Libraries[name] = installedVersion

	if err := writeAridJSON(config); err != nil {
		l.logger.WithError(err).Error(errMsg)
		return err
	}

	return nil
}

// Remove library either globally or for project
func (l *Lib) Remove(name string) error {
	if !isProjectDirectory() {
		err := errors.New("Directory not initialized. Run \"ardi lib init\" first")
		l.logger.WithError(err).Error("Cannot remove library")
		return err
	}

	l.logger.Infof("Removing library: %s", name)
	if err := l.RPC.UninstallLibrary(name); err != nil {
		return err
	}
	config := ArdiJSON{}
	errMsg := fmt.Sprintf("Failed to remove library: %s", name)

	buildConfig, err := ioutil.ReadFile(paths.ArdiBuildConfig)
	if err != nil {
		l.logger.WithError(err).Error(errMsg)
		return err
	}

	if err := json.Unmarshal(buildConfig, &config); err != nil {
		l.logger.WithError(err).Error(errMsg)
		return err
	}

	for lib := range config.Libraries {
		if lib == name {
			delete(config.Libraries, lib)
		}
	}

	if err := writeAridJSON(config); err != nil {
		l.logger.WithError(err).Error(errMsg)
		return err
	}

	return nil
}

// Install all libraries specified in ardi.json
func (l *Lib) Install() error {
	if !isProjectDirectory() {
		err := errors.New("Directory not initialized. Run \"ardi lib init\" first")
		l.logger.WithError(err).Error("Cannot install libraries")
		return err
	}

	config := ArdiJSON{}
	errMsg := "Failed to install libraries"

	buildConfig, err := ioutil.ReadFile(paths.ArdiBuildConfig)
	if err != nil {
		l.logger.WithError(err).Error(errMsg)
		return err
	}

	if err := json.Unmarshal(buildConfig, &config); err != nil {
		l.logger.WithError(err).Error(errMsg)
		return err
	}

	for name, version := range config.Libraries {
		if err := l.Add(name, version); err != nil {
			return err
		}
	}

	return nil
}

//private helpers
func isProjectDirectory() bool {
	_, err := os.Stat(paths.ArdiBuildConfig)
	return !os.IsNotExist(err)
}

func writeAridJSON(config ArdiJSON) error {
	newData, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(paths.ArdiBuildConfig, newData, 0644); err != nil {
		return err
	}

	return nil
}
