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
	logger *log.Logger
	client rpc.Client
}

// NewPlatformCore platform module instance
func NewPlatformCore(client rpc.Client, logger *log.Logger) *PlatformCore {
	return &PlatformCore{
		logger: logger,
		client: client,
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

	p.logger.Info("------INSTALLED PLATFORMS------")
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
func (p *PlatformCore) Add(platform string) error {
	if platform == "" {
		return errors.New("Empty platform list")
	}

	if err := p.client.InstallPlatform(platform); err != nil {
		return err
	}

	return nil
}

// AddAll installs all platforms
func (p *PlatformCore) AddAll() error {
	return p.client.InstallAllPlatforms()
}

// Remove uninstalls specified platforms
func (p *PlatformCore) Remove(platform string) error {
	if platform == "" {
		return errors.New("Empty platform list")
	}

	if err := p.client.UninstallPlatform(platform); err != nil {
		return err
	}

	return nil
}
