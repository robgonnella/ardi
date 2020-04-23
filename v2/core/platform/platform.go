package platform

import (
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
	Client *rpc.Client
}

// New platform module instance
func New(logger *log.Logger) (*Platform, error) {
	client, err := rpc.NewClient(logger)
	if err != nil {
		return nil, err
	}
	return &Platform{
		logger: logger,
		Client: client,
	}, nil
}

// List all available platforms or filter with a search arg
func (p *Platform) List(query string) error {
	platforms, err := p.Client.GetPlatforms(query)
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
