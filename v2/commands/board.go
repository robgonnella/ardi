package commands

import (
	"github.com/robgonnella/ardi/v2/core/board"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getBoardFQBNSCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "fqbns",
		Long:    cyan("\nList supported board fqbns"),
		Short:   "List supported board fqbns",
		Aliases: []string{"fqbn"},
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
			query := ""
			if len(args) > 0 {
				query = args[0]
			}

			boardCore, err := board.New(logger)
			if err != nil {
				return
			}
			defer boardCore.Client.Connection.Close()
			boardCore.FQBNS(query)
		},
	}
	return listCmd
}

func getBoardPlatformsCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "platforms",
		Long:    cyan("\nList boards with their associated platform"),
		Short:   "List boards with their associated platform",
		Aliases: []string{"platform"},
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
			query := ""
			if len(args) > 0 {
				query = args[0]
			}

			boardCore, err := board.New(logger)
			if err != nil {
				return
			}
			defer boardCore.Client.Connection.Close()
			boardCore.Platforms(query)
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
	boardCmd.AddCommand(getBoardPlatformsCmd())
	boardCmd.AddCommand(getBoardFQBNSCmd())

	return boardCmd
}
