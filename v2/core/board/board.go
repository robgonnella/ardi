package board

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/arduino/arduino-cli/rpc/commands"
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

// FQBNS all available boards with optional search filter
func (b *Board) FQBNS(query string) error {
	platforms, err := b.Client.GetPlatforms(query)

	if err != nil {
		b.logger.WithError(err).Error("Platform search error")
		return err
	}

	var boardList []*commands.Board

	for _, plat := range platforms {
		for _, board := range plat.GetBoards() {
			if board.GetFqbn() != "" {
				boardList = append(boardList, board)
			}
		}
	}

	if len(boardList) == 0 {
		b.logger.Info("You must install platforms with \"ardi init\" first")
		return nil
	}

	sort.Slice(boardList, func(i, j int) bool {
		return boardList[i].GetName() < boardList[j].GetName()
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', 0)
	defer w.Flush()
	fmt.Fprintln(w, "Board\tFQBN")
	for _, board := range boardList {
		fmt.Fprintf(w, "%s\t%s\n", board.GetName(), board.GetFqbn())
	}
	return nil
}

// Platforms all available boards with optional search filter
func (b *Board) Platforms(query string) error {
	platforms, err := b.Client.GetPlatforms(query)

	if err != nil {
		b.logger.WithError(err).Error("Platform search error")
		return err
	}

	type boardAndPlat struct {
		boardName string
		platform  string
	}

	var boardList []boardAndPlat
	for _, plat := range platforms {
		for _, board := range plat.GetBoards() {
			boardList = append(
				boardList,
				boardAndPlat{
					boardName: board.GetName(),
					platform:  plat.GetID(),
				},
			)
		}
	}

	sort.Slice(boardList, func(i, j int) bool {
		return boardList[i].boardName < boardList[j].boardName
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', 0)
	defer w.Flush()
	fmt.Fprintln(w, "Board\tPlatform")
	for _, board := range boardList {
		fmt.Fprintf(w, "%s\t%s\n", board.boardName, board.platform)
	}

	return nil
}
