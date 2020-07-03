package commands

import (
	"github.com/spf13/cobra"
)

func getBoardFQBNSCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "fqbns",
		Long:    "\nList boards with associated fqbns",
		Short:   "List boards with associated fqbns",
		Aliases: []string{"fqbn"},
		Run: func(cmd *cobra.Command, args []string) {
			query := ""
			if len(args) > 0 {
				query = args[0]
			}
			if err := ardiCore.Board.FQBNS(query); err != nil {
				logger.WithError(err).Error("Failed to print board fqbns")
			}
		},
	}
	return listCmd
}

func getBoardPlatformsCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "platforms",
		Long:    "\nList boards with their associated platform",
		Short:   "List boards with their associated platform",
		Aliases: []string{"platform"},
		Run: func(cmd *cobra.Command, args []string) {
			query := ""
			if len(args) > 0 {
				query = args[0]
			}
			if err := ardiCore.Board.Platforms(query); err != nil {
				logger.WithError(err).Error("Failed to print board platforms")
			}
		},
	}
	return listCmd
}

func getBoardCommand() *cobra.Command {
	boardCmd := &cobra.Command{
		Use: "board",
		Long: "\nBoard helper allowing you to see which boards belong to " +
			"which platforms, and the FQBN associated with each board",
		Short: "Board helper",
	}
	boardCmd.AddCommand(getBoardPlatformsCmd())
	boardCmd.AddCommand(getBoardFQBNSCmd())

	return boardCmd
}
