package core

import (
	"io/ioutil"

	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/types"
	"github.com/robgonnella/ardi/v2/util"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// ArdiYAML represents core module for data config file manipulations
type ArdiYAML struct {
	Config types.DataConfig
	logger *log.Logger
}

// NewArdiYAML returns core yaml module for handling data config file
func NewArdiYAML(logger *log.Logger) (*ArdiYAML, error) {
	config := types.DataConfig{}
	dataConfig, err := ioutil.ReadFile(paths.ArdiProjectDataConfig)
	if err != nil {
		logger.WithError(err).Error("Failed to read data config file")
		return nil, err
	}
	if err := yaml.Unmarshal(dataConfig, &config); err != nil {
		logger.WithError(err).Error("Failed to parse data file")
		return nil, err
	}
	return &ArdiYAML{
		Config: config,
		logger: logger,
	}, nil
}

// AddBoardURL add a board url to data config file
func (a *ArdiYAML) AddBoardURL(url string) error {
	if !util.ArrayContains(a.Config.BoardManager.AdditionalUrls, url) {
		a.logger.Infof("Adding board url: %s", url)
		a.Config.BoardManager.AdditionalUrls = append(a.Config.BoardManager.AdditionalUrls, url)
		return a.write()
	}
	return nil
}

// private methods
func (a *ArdiYAML) write() error {
	newData, err := yaml.Marshal(a.Config)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(paths.ArdiProjectDataConfig, newData, 0644); err != nil {
		return err
	}

	return nil
}
