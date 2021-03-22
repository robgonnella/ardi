package core

import (
	"io/ioutil"
	"sync"

	"github.com/robgonnella/ardi/v2/types"
	"github.com/robgonnella/ardi/v2/util"
	"gopkg.in/yaml.v2"
)

// ArdiYAML represents core module for data config file manipulations
type ArdiYAML struct {
	Config   types.ArduinoCliSettings
	confPath string
	mux      sync.Mutex
}

// NewArdiYAML returns core yaml module for handling data config file
func NewArdiYAML(confPath string, initalConfig types.ArduinoCliSettings) *ArdiYAML {
	return &ArdiYAML{
		Config:   initalConfig,
		confPath: confPath,
		mux:      sync.Mutex{},
	}
}

// AddBoardURL add a board url to data config file
func (a *ArdiYAML) AddBoardURL(url string) error {
	if !util.ArrayContains(a.Config.BoardManager.AdditionalUrls, url) {
		a.Config.BoardManager.AdditionalUrls = append(a.Config.BoardManager.AdditionalUrls, url)
		return a.write()
	}
	return nil
}

// RemoveBoardURL remove a board url from data config file
func (a *ArdiYAML) RemoveBoardURL(url string) error {
	if util.ArrayContains(a.Config.BoardManager.AdditionalUrls, url) {
		newList := []string{}
		for _, u := range a.Config.BoardManager.AdditionalUrls {
			if u != url {
				newList = append(newList, u)
			}
		}
		a.Config.BoardManager.AdditionalUrls = newList
		return a.write()
	}
	return nil
}

// private methods
func (a *ArdiYAML) write() error {
	a.mux.Lock()
	defer a.mux.Unlock()

	newData, err := yaml.Marshal(a.Config)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(a.confPath, newData, 0644); err != nil {
		return err
	}

	return nil
}
