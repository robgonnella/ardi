package commands

import (
	"github.com/robgonnella/ardi/ardi"
	"github.com/spf13/cobra"
)

func getPlatformListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "list",
		Short:   "List all available platforms",
		Aliases: []string{"search"},
		Run: func(cmd *cobra.Command, args []string) {
			platform := ""
			if len(args) > 0 {
				platform = args[0]
			}
			ardi.ListPlatforms(platform)
		},
	}
	return listCmd
}

func getPlatformCommand() *cobra.Command {
	platCmd := &cobra.Command{
		Use:   "platform",
		Short: "Platform related commands",
	}
	platCmd.AddCommand(getPlatformListCmd())

	return platCmd
}
