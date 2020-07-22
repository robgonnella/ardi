package core

import (
	"errors"
	"fmt"
	"sort"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"

	"github.com/robgonnella/ardi/v2/rpc"
)

// PlatformCore module for platform commands
type PlatformCore struct {
	logger      *log.Logger
	client      rpc.Client
	initialized bool
}

// NewPlatformCore platform module instance
func NewPlatformCore(client rpc.Client, logger *log.Logger) *PlatformCore {
	return &PlatformCore{
		logger:      logger,
		client:      client,
		initialized: false,
	}
}

// ListInstalled lists only installed platforms
func (p *PlatformCore) ListInstalled() error {
	platforms, err := p.client.GetInstalledPlatforms()
	if err != nil {
		return err
	}

	sort.Slice(platforms, func(i, j int) bool {
		return platforms[i].GetName() < platforms[j].GetName()
	})

	w := tabwriter.NewWriter(p.logger.Out, 0, 0, 8, ' ', 0)
	defer w.Flush()
	w.Write([]byte("Platform\tID\tInstalled\n"))
	for _, plat := range platforms {
		w.Write([]byte(fmt.Sprintf("%s\t%s\t%s\n", plat.GetName(), plat.GetID(), plat.GetInstalled())))
	}
	return nil
}

// ListAll lists all available platforms
func (p *PlatformCore) ListAll() error {
	p.init()

	platforms, err := p.client.GetPlatforms()
	if err != nil {
		return err
	}

	sort.Slice(platforms, func(i, j int) bool {
		return platforms[i].GetName() < platforms[j].GetName()
	})

	p.logger.Info("------AVAILABLE PLATFORMS------")
	w := tabwriter.NewWriter(p.logger.Out, 0, 0, 8, ' ', 0)
	defer w.Flush()
	w.Write([]byte("Platform\tID\tLatest\n"))
	for _, plat := range platforms {
		w.Write([]byte(fmt.Sprintf("%s\t%s\t%s\n", plat.GetName(), plat.GetID(), plat.GetLatest())))
	}
	return nil
}

// Add installs specified platforms
func (p *PlatformCore) Add(platform string) (string, string, error) {
	p.init()

	if platform == "" {
		return "", "", errors.New("Empty platform list")
	}

	installed, vers, err := p.client.InstallPlatform(platform)
	if err != nil {
		return "", "", err
	}

	return installed, vers, nil
}

// Remove uninstalls specified platforms
func (p *PlatformCore) Remove(platform string) (string, error) {
	if platform == "" {
		return "", errors.New("Empty platform list")
	}

	removed, err := p.client.UninstallPlatform(platform)
	if err != nil {
		return "", err
	}

	return removed, nil
}

// private
func (p *PlatformCore) init() error {
	if !p.initialized {
		if err := p.client.UpdatePlatformIndex(); err != nil {
			return err
		}
		p.initialized = true
	}
	return nil
}
