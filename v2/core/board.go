package core

import (
	"errors"
	"fmt"
	"sort"
	"text/tabwriter"

	"github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	log "github.com/sirupsen/logrus"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
)

// BoardCore module for board commands
type BoardCore struct {
	cli    *cli.Wrapper
	logger *log.Logger
}

// BoardCoreOption represents options for the BoardCore
type BoardCoreOption = func(c *BoardCore)

// NewBoardCore module instance for board commands
func NewBoardCore(logger *log.Logger, options ...BoardCoreOption) *BoardCore {
	c := &BoardCore{
		logger: logger,
	}

	for _, o := range options {
		o(c)
	}

	return c
}

// WithBoardCliWrapper allows an injectable cli wrapper
func WithBoardCliWrapper(wrapper *cli.Wrapper) BoardCoreOption {
	return func(c *BoardCore) {
		c.cli = wrapper
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
					platform:  plat.GetId(),
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
