package target

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"

	"github.com/robgonnella/ardi/v2/rpc"
)

// Target represents a targeted arduino board
type Target struct {
	Board *rpc.Board
}

// New returns new target
func New(logger *log.Logger, fqbn string, onlyConnected bool) (*Target, error) {
	client, err := rpc.NewClient(logger)
	if err != nil {
		return nil, err
	}
	defer client.Connection.Close()
	board, err := getTargetBoard(client, logger, fqbn, onlyConnected)
	if err != nil {
		return nil, err
	}
	return &Target{board}, nil
}

func getTargetBoard(client *rpc.Client, logger *log.Logger, fqbn string, onlyConnected bool) (*rpc.Board, error) {
	if fqbn != "" {
		return &rpc.Board{FQBN: fqbn}, nil
	}

	connectedBoards := client.ConnectedBoards()
	allBoards := client.AllBoards()

	if len(connectedBoards) == 0 {
		if onlyConnected {
			err := errors.New("No connected boards detected")
			logger.WithError(err).Error()
			return nil, err
		}
		printFQBNs(allBoards)
		fmt.Printf("\nYou must supply fqbn to compile. You can find a list of board fqbns for installed platforms above.\n\n")
		return nil, errors.New("Must provide board fqbn")
	}

	if len(connectedBoards) == 1 {
		return connectedBoards[0], nil
	}

	if len(connectedBoards) > 1 {
		printFQBNs(connectedBoards)
		fmt.Printf("\nYou must supply fqbn to compile. You can find a list of board fqbns for connected boards above.\n\n")
		return nil, errors.New("Must provide board fqbn")
	}

	return nil, errors.New("Error parsing target")
}

// private helpers
func printFQBNs(boardList []*rpc.Board) {
	sort.Slice(boardList, func(i, j int) bool {
		return boardList[i].Name < boardList[j].Name
	})

	printBoardsWithIndices(boardList)
}

func printBoardsWithIndices(boards []*rpc.Board) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', 0)
	defer w.Flush()
	fmt.Fprintln(w, "No.\tName\tFQBN")
	for i, b := range boards {
		fmt.Fprintf(w, "%d\t%s\t%s\n", i, b.Name, b.FQBN)
	}
}
