package core

import (
	"errors"
	"fmt"
	"sort"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
)

// PlatformCore module for platform commands
type PlatformCore struct {
	logger      *log.Logger
	cli         *cli.Wrapper
	initialized bool
}

// NewPlatformCore platform module instance
func NewPlatformCore(cli *cli.Wrapper, logger *log.Logger) *PlatformCore {
	return &PlatformCore{
		logger:      logger,
		cli:         cli,
		initialized: false,
	}
}

// ListInstalled lists only installed platforms
func (c *PlatformCore) ListInstalled() error {
	platforms, err := c.cli.GetInstalledPlatforms()
	if err != nil {
		return err
	}

	sort.Slice(platforms, func(i, j int) bool {
		return platforms[i].GetName() < platforms[j].GetName()
	})

	w := tabwriter.NewWriter(c.logger.Out, 0, 0, 8, ' ', 0)
	defer w.Flush()
	w.Write([]byte("Platform\tID\tInstalled\n"))
	for _, plat := range platforms {
		w.Write([]byte(fmt.Sprintf("%s\t%s\t%s\n", plat.GetName(), plat.GetId(), plat.GetInstalled())))
	}
	return nil
}

// ListAll lists all available platforms
func (c *PlatformCore) ListAll() error {
	c.init()

	platforms, err := c.cli.SearchPlatforms()
	if err != nil {
		return err
	}

	sort.Slice(platforms, func(i, j int) bool {
		return platforms[i].GetName() < platforms[j].GetName()
	})

	c.logger.Info("------AVAILABLE PLATFORMS------")
	w := tabwriter.NewWriter(c.logger.Out, 0, 0, 8, ' ', 0)
	defer w.Flush()
	w.Write([]byte("Platform\tID\tLatest\n"))
	for _, plat := range platforms {
		w.Write([]byte(fmt.Sprintf("%s\t%s\t%s\n", plat.GetName(), plat.GetId(), plat.GetLatest())))
	}
	return nil
}

// Add installs specified platforms
func (c *PlatformCore) Add(platform string) (string, string, error) {
	c.init()

	if platform == "" {
		return "", "", errors.New("empty platform list")
	}

	installed, vers, err := c.cli.InstallPlatform(platform)
	if err != nil {
		return "", "", err
	}

	c.logger.Infof("Installed Platform: %s %s", installed, vers)
	return installed, vers, nil
}

// Remove uninstalls specified platforms
func (c *PlatformCore) Remove(platform string) (string, error) {
	if platform == "" {
		return "", errors.New("empty platform list")
	}

	removed, err := c.cli.UninstallPlatform(platform)
	if err != nil {
		return "", err
	}

	return removed, nil
}

// private
func (c *PlatformCore) init() error {
	if !c.initialized {
		if err := c.cli.UpdatePlatformIndex(); err != nil {
			return err
		}
		c.initialized = true
	}
	return nil
}
