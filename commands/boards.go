package commands

import (
	"github.com/robgonnella/ardi/ardi"
	"github.com/spf13/cobra"
)

func getBoardListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "list",
		Short:   "List all available boards",
		Aliases: []string{"search"},
		Run: func(cmd *cobra.Command, args []string) {
			board := ""
			if len(args) > 0 {
				board = args[0]
			}
			ardi.ListBoards(board)
		},
	}
	return listCmd
}

func getBoardCommand() *cobra.Command {
	boardCmd := &cobra.Command{
		Use:   "board",
		Short: "Board related commands",
	}
	boardCmd.AddCommand(getBoardListCmd())

	return boardCmd
}
