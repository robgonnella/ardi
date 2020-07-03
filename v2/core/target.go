package core

import (
	"errors"
	"fmt"
	"sort"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"

	"github.com/robgonnella/ardi/v2/rpc"
)

// Target represents a targeted arduino board
type Target struct {
	Board *rpc.Board
}

// NewTarget returns new target
func NewTarget(connectedBoards, allBoards []*rpc.Board, fqbn string, onlyConnected bool, logger *log.Logger) (*Target, error) {
	board, err := getTargetBoard(connectedBoards, allBoards, fqbn, onlyConnected, logger)
	if err != nil {
		return nil, err
	}
	return &Target{board}, nil
}

func getTargetBoard(connectedBoards, allBoards []*rpc.Board, fqbn string, onlyConnected bool, logger *log.Logger) (*rpc.Board, error) {
	if fqbn != "" {
		return &rpc.Board{FQBN: fqbn}, nil
	}

	fqbnErr := errors.New("you must specify a board fqbn to compile - you can find a list of board fqbns for installed platforms above")

	if len(connectedBoards) == 0 {
		if onlyConnected {
			err := errors.New("No connected boards detected")
			return nil, err
		}
		printFQBNs(allBoards, logger)
		return nil, fqbnErr
	}

	if len(connectedBoards) == 1 {
		return connectedBoards[0], nil
	}

	if len(connectedBoards) > 1 {
		printFQBNs(connectedBoards, logger)
		return nil, fqbnErr
	}

	return nil, errors.New("Error parsing target")
}

// private helpers
func printFQBNs(boardList []*rpc.Board, logger *log.Logger) {
	sort.Slice(boardList, func(i, j int) bool {
		return boardList[i].Name < boardList[j].Name
	})

	printBoardsWithIndices(boardList, logger)
}

func printBoardsWithIndices(boards []*rpc.Board, logger *log.Logger) {
	w := tabwriter.NewWriter(logger.Out, 0, 0, 8, ' ', 0)
	defer w.Flush()
	w.Write([]byte("No.\tName\tFQBN\n"))
	for i, b := range boards {
		w.Write([]byte(fmt.Sprintf("%d\t%s\t%s\n", i, b.Name, b.FQBN)))
	}
}
