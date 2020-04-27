package platform

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"

	"github.com/robgonnella/ardi/v2/rpc"
)

// Platform module for platform commands
type Platform struct {
	logger *log.Logger
	client *rpc.Client
}

// New platform module instance
func New(client *rpc.Client, logger *log.Logger) (*Platform, error) {
	if err := client.UpdateIndexFiles(); err != nil {
		logger.WithError(err).Error("Failed to update index files")
		return nil, err
	}
	return &Platform{
		logger: logger,
		client: client,
	}, nil
}

// ListInstalled lists only installed platforms
func (p *Platform) ListInstalled() error {
	platforms, err := p.client.GetInstalledPlatforms()
	if err != nil {
		return err
	}

	sort.Slice(platforms, func(i, j int) bool {
		return platforms[i].GetName() < platforms[j].GetName()
	})

	p.logger.Info("------INSTALLED PLATFORMS------")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', 0)
	defer w.Flush()
	fmt.Fprintln(w, "Platform\tID\tInstalled")
	for _, plat := range platforms {
		fmt.Fprintf(w, "%s\t%s\t%s\n", plat.GetName(), plat.GetID(), plat.GetInstalled())
	}
	return nil
}

// ListAll lists all available platforms
func (p *Platform) ListAll() error {
	platforms, err := p.client.GetPlatforms()
	if err != nil {
		return err
	}

	sort.Slice(platforms, func(i, j int) bool {
		return platforms[i].GetName() < platforms[j].GetName()
	})

	p.logger.Info("------AVAILABLE PLATFORMS------")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', 0)
	defer w.Flush()
	fmt.Fprintln(w, "Platform\tID\tLatest")
	for _, plat := range platforms {
		fmt.Fprintf(w, "%s\t%s\t%s\n", plat.GetName(), plat.GetID(), plat.GetLatest())
	}
	return nil
}

// Add installs specified platforms
func (p *Platform) Add(platforms []string) error {
	if len(platforms) == 0 {
		err := errors.New("Empty platform list")
		p.logger.WithError(err).Error()
		return err
	}

	for _, platform := range platforms {
		if err := p.client.InstallPlatform(platform); err != nil {
			p.logger.WithError(err).Error()
			return err
		}
	}

	return nil
}

// AddAll installs all platforms
func (p *Platform) AddAll() error {
	return p.client.InstallAllPlatforms()
}

// Remove uninstalls specified platforms
func (p *Platform) Remove(platforms []string) error {
	if len(platforms) == 0 {
		err := errors.New("Empty platform list")
		p.logger.WithError(err).Error()
		return err
	}

	for _, platform := range platforms {
		if err := p.client.UninstallPlatform(platform); err != nil {
			p.logger.WithError(err).Error()
			return err
		}
	}

	return nil
}
