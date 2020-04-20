package target

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"

	"github.com/robgonnella/ardi/v2/core/rpc"
)

// Target represents a targeted arduino board
type Target struct {
	Board *rpc.Board
}

// New returns new target
func New(rpc *rpc.RPC, logger *log.Logger, fqbn string, onlyConnected bool) (*Target, error) {
	board, err := getTargetBoard(rpc, logger, fqbn, onlyConnected)
	if err != nil {
		return nil, err
	}
	return &Target{board}, nil
}

func getTargetBoard(server *rpc.RPC, logger *log.Logger, fqbn string, onlyConnected bool) (*rpc.Board, error) {
	var board *rpc.Board
	var err error

	if fqbn != "" {
		board.FQBN = fqbn
		return board, nil
	}

	connectedBoards := server.ConnectedBoards()
	allBoards := server.AllBoards()

	if len(connectedBoards) == 0 {
		if onlyConnected {
			err := errors.New("No connected boards detected")
			logger.WithError(err).Error()
			return nil, err
		}
		board, err = getUserInput(allBoards)
	}

	if len(connectedBoards) == 1 {
		return connectedBoards[0], nil
	}

	if len(connectedBoards) > 1 {
		board, err = getUserInput(connectedBoards)
	}

	if err != nil {
		logger.WithError(err).Error("Failed to parse target board")
		return nil, err
	}

	return board, nil
}

// private helpers
func getUserInput(boardList []*rpc.Board) (*rpc.Board, error) {
	printBoardsWithIndices(boardList)

	var boardIdx int
	fmt.Print("\nEnter number of board for which to compile: ")
	if _, err := fmt.Scanf("%d", &boardIdx); err != nil {
		return nil, err
	}

	if boardIdx < 0 || boardIdx > len(boardList)-1 {
		err := errors.New("Invalid board selection")
		return nil, err
	}

	return boardList[boardIdx], nil
}

func printBoardsWithIndices(boards []*rpc.Board) {
	sort.Slice(boards, func(i, j int) bool {
		return boards[i].Name < boards[j].Name
	})
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', 0)
	defer w.Flush()
	fmt.Fprintln(w, "No.\tName\tFQBN")
	for i, b := range boards {
		fmt.Fprintf(w, "%d\t%s\t%s\n", i, b.Name, b.FQBN)
	}
}
