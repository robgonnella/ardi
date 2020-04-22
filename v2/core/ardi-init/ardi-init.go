package ardiinit

import (
	"fmt"
	"strings"
	"time"

	"github.com/robgonnella/ardi/v2/rpc"
	log "github.com/sirupsen/logrus"
)

// Init represents core module for init commands
type Init struct {
	logger *log.Logger
	Client *rpc.Client
}

// New init module instance
func New(logger *log.Logger) (*Init, error) {
	client, err := rpc.NewClient(logger)
	if err != nil {
		return nil, err
	}
	return &Init{
		logger: logger,
		Client: client,
	}, nil
}

// Initialize ardi by downloading platform cores to data directory
func (i *Init) Initialize(platform, version string) error {
	var err error
	errMsg := "Failed to initialize ardi"
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
		err = i.Client.InstallAllPlatforms()
	} else {
		platParts := strings.Split(platform, ":")
		platPackage := platParts[0]
		arch := platParts[len(platParts)-1]
		err = i.Client.InstallPlatform(platPackage, arch, version)
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
