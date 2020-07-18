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

// NewTargetOpts options for creating a new target for compiling and uploading
type NewTargetOpts struct {
	FQBN            string
	ConnectedBoards []*rpc.Board
	AllBoards       []*rpc.Board
	OnlyConnected   bool
	Logger          *log.Logger
}

// NewTarget returns new target
func NewTarget(opts NewTargetOpts) (*Target, error) {
	board, err := getTargetBoard(opts)
	if err != nil {
		return nil, err
	}
	return &Target{board}, nil
}

func getTargetBoard(opts NewTargetOpts) (*rpc.Board, error) {
	if opts.FQBN != "" {
		return &rpc.Board{FQBN: opts.FQBN}, nil
	}

	fqbnErr := errors.New("you must specify a board fqbn to compile - you can find a list of board fqbns for installed platforms above")

	if len(opts.ConnectedBoards) == 0 {
		if opts.OnlyConnected {
			err := errors.New("No connected boards detected")
			return nil, err
		}
		printFQBNs(opts.AllBoards, opts.Logger)
		return nil, fqbnErr
	}

	if len(opts.ConnectedBoards) == 1 {
		return opts.ConnectedBoards[0], nil
	}

	if len(opts.ConnectedBoards) > 1 {
		printFQBNs(opts.ConnectedBoards, opts.Logger)
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
