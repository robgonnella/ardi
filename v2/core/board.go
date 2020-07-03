package core

import (
	"errors"
	"fmt"
	"sort"
	"text/tabwriter"

	"github.com/arduino/arduino-cli/rpc/commands"
	log "github.com/sirupsen/logrus"

	"github.com/robgonnella/ardi/v2/rpc"
)

// BoardCore module for board commands
type BoardCore struct {
	client rpc.Client
	logger *log.Logger
}

// NewBoardCore module instance for board commands
func NewBoardCore(client rpc.Client, logger *log.Logger) *BoardCore {
	return &BoardCore{
		logger: logger,
		client: client,
	}
}

// FQBNS lists all available boards with associated fqbns
func (b *BoardCore) FQBNS(query string) error {
	platforms, err := b.client.GetPlatforms()

	if err != nil {
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
		err := errors.New("you must install platforms with 'ardi platform add' or 'ardi project add platform' first")
		return err
	}

	sort.Slice(boardList, func(i, j int) bool {
		return boardList[i].GetName() < boardList[j].GetName()
	})

	w := tabwriter.NewWriter(b.logger.Out, 0, 0, 8, ' ', 0)
	defer w.Flush()
	w.Write([]byte("Board\tFQBN\n"))
	for _, board := range boardList {
		w.Write([]byte(fmt.Sprintf("%s\t%s\n", board.GetName(), board.GetFqbn())))
	}
	return nil
}

// Platforms lists all available boards with associated platorms
func (b *BoardCore) Platforms(query string) error {
	platforms, err := b.client.GetPlatforms()

	if err != nil {
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

	w := tabwriter.NewWriter(b.logger.Out, 0, 0, 8, ' ', 0)
	defer w.Flush()
	w.Write([]byte("Board\tPlatform\n"))
	for _, board := range boardList {
		w.Write([]byte(fmt.Sprintf("%s\t%s\n", board.boardName, board.platform)))
	}

	return nil
}
