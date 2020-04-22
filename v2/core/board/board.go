package board

import (
	"fmt"
	"os"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"

	"github.com/robgonnella/ardi/v2/rpc"
)

// Board module for board commands
type Board struct {
	Client *rpc.Client
	logger *log.Logger
}

// New module instance for board commands
func New(logger *log.Logger) (*Board, error) {
	client, err := rpc.NewClient(logger)
	if err != nil {
		return nil, err
	}

	return &Board{
		logger: logger,
		Client: client,
	}, nil
}

// List all available boards with optional search filter
func (b *Board) List(query string) error {
	platforms, err := b.Client.GetPlatforms(query)

	if err != nil {
		b.logger.WithError(err).Error("Platform search error")
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', 0)
	defer w.Flush()
	fmt.Fprintln(w, "Board\tPlatform\tFQBN")
	for _, plat := range platforms {
		for _, board := range plat.GetBoards() {
			fmt.Fprintf(w, "%s\t%s\t%s\n", board.GetName(), plat.GetID(), board.GetFqbn())
		}
	}
	return nil
}
