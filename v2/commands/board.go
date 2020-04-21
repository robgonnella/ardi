package commands

import (
	"github.com/robgonnella/ardi/v2/core/board"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getBoardListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "list",
		Long:    cyan("\nList all available boards"),
		Short:   "List all available boards",
		Aliases: []string{"search"},
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
			query := ""
			if len(args) > 0 {
				query = args[0]
			}

			board, err := board.New(logger)
			if err != nil {
				return
			}
			defer board.RPC.Connection.Close()
			board.List(query)
		},
	}
	return listCmd
}

func getBoardCommand() *cobra.Command {
	boardCmd := &cobra.Command{
		Use:   "board",
		Long:  cyan("\nBoard related commands"),
		Short: "Board related commands",
	}
	boardCmd.AddCommand(getBoardListCmd())

	return boardCmd
}
