package core

import (
	"errors"
	"fmt"
	"sort"
	"text/tabwriter"

	"github.com/arduino/arduino-cli/rpc/commands"
	log "github.com/sirupsen/logrus"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
)

// BoardCore module for board commands
type BoardCore struct {
	cli    *cli.Wrapper
	logger *log.Logger
}

// NewBoardCore module instance for board commands
func NewBoardCore(cli *cli.Wrapper, logger *log.Logger) *BoardCore {
	return &BoardCore{
		logger: logger,
		cli:    cli,
	}
}

// FQBNS lists all available boards with associated fqbns
func (c *BoardCore) FQBNS(query string) error {
	platforms, err := c.cli.SearchPlatforms()

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
		err := errors.New("you must install platforms with 'ardi add platform'")
		return err
	}

	sort.Slice(boardList, func(i, j int) bool {
		return boardList[i].GetName() < boardList[j].GetName()
	})

	w := tabwriter.NewWriter(c.logger.Out, 0, 0, 8, ' ', 0)
	defer w.Flush()
	w.Write([]byte("Board\tFQBN\n"))
	for _, board := range boardList {
		w.Write([]byte(fmt.Sprintf("%s\t%s\n", board.GetName(), board.GetFqbn())))
	}
	return nil
}

// Platforms lists all available boards with associated platorms
func (c *BoardCore) Platforms(query string) error {
	platforms, err := c.cli.SearchPlatforms()

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

	w := tabwriter.NewWriter(c.logger.Out, 0, 0, 8, ' ', 0)
	defer w.Flush()
	w.Write([]byte("Board\tPlatform\n"))
	for _, board := range boardList {
		w.Write([]byte(fmt.Sprintf("%s\t%s\n", board.boardName, board.platform)))
	}

	return nil
}
