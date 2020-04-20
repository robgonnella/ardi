package ardiinit

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/robgonnella/ardi/v2/core/rpc"
	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/types"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Init represents core module for init commands
type Init struct {
	logger *log.Logger
	RPC    *rpc.RPC
}

// New init module instance
func New(logger *log.Logger) (*Init, error) {
	rpc, err := rpc.New(paths.ArdiGlobalDataConfig, logger)
	if err != nil {
		logger.WithError(err).Error("Failed to initialize ardi")
		return nil, err
	}
	return &Init{
		logger: logger,
		RPC:    rpc,
	}, nil
}

// Initialize ardi by downloading platform cores to data directory
func (i *Init) Initialize(platform, version string) error {
	var err error
	errMsg := "Failed to initialize ardi"
	if err = initializeDataDirectory(); err != nil {
		i.logger.WithError(err).Error(errMsg)
		return err
	}

	i.logger.Info("Initializing. This may take some time...")
	quit := make(chan bool, 1)

	// Show simple "processing" indicator if not logging verbosely
	if i.isQuiet() {
		i.logger.Info("Installing platforms...")
		ticker := time.NewTicker(2 * time.Second)
		go func() {
			for {
				select {
				case <-ticker.C:
					fmt.Print(".")
				case <-quit:
					ticker.Stop()
				}
			}
		}()
	}

	if platform == "" {
		err = i.RPC.InstallAllPlatforms()
	} else {
		platParts := strings.Split(platform, ":")
		platPackage := platParts[0]
		arch := platParts[len(platParts)-1]
		err = i.RPC.InstallPlatform(platPackage, arch, version)
	}

	quit <- true
	fmt.Println("")
	if err != nil {
		i.logger.WithError(err).Error(errMsg)
		return err
	}

	i.logger.Info("Successfully initialized!")
	fmt.Println("")
	return nil
}

// private methods
func (i *Init) isQuiet() bool {
	return i.logger.Level == log.InfoLevel
}

// private helpers
func initializeDataDirectory() error {
	if _, err := os.Stat(paths.ArdiGlobalDataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(paths.ArdiGlobalDataDir, 0744); err != nil {
			return err
		}
	}

	if _, err := os.Stat(paths.ArdiGlobalDataConfig); os.IsNotExist(err) {
		dataConfig := types.DataConfig{
			ProxyType:      "auto",
			SketchbookPath: ".",
			ArduinoData:    ".",
			BoardManager:   make(map[string]interface{}),
		}
		yamlConfig, _ := yaml.Marshal(&dataConfig)
		if err := ioutil.WriteFile(paths.ArdiGlobalDataConfig, yamlConfig, 0644); err != nil {
			return err
		}
	}

	return nil
}
