package commands

import (
	"github.com/robgonnella/ardi/v2/core/platform"
	"github.com/spf13/cobra"
)

func getPlatformListCmd() *cobra.Command {
	var all bool
	var installed bool
	listCmd := &cobra.Command{
		Use:   "list",
		Long:  "\nList platforms",
		Short: "List platforms",
		Run: func(cmd *cobra.Command, args []string) {
			platformCore, err := platform.New(client, logger)
			if err != nil {
				return
			}
			if all || (!all && !installed) {
				platformCore.ListAll()
			}
			if installed {
				platformCore.ListInstalled()
			}
		},
	}
	listCmd.Flags().BoolVarP(&all, "all", "a", false, "List all platforms")
	listCmd.Flags().BoolVarP(&installed, "installed", "i", false, "List only installed platforms")
	return listCmd
}

func getPlatformAddCmd() *cobra.Command {
	var all bool
	addCmd := &cobra.Command{
		Use:   "add",
		Long:  "\nInstall platforms",
		Short: "Install platforms",
		Run: func(cmd *cobra.Command, args []string) {
			platformCore, err := platform.New(client, logger)
			if err != nil {
				return
			}

			if all {
				platformCore.AddAll()
				return
			}

			platformCore.Add(args)
		},
	}

	addCmd.Flags().BoolVarP(&all, "all", "a", false, "Install all platforms")
	return addCmd
}

func getPlatformRemoveCmd() *cobra.Command {
	removeCmd := &cobra.Command{
		Use:   "remove",
		Long:  "\nRemove installed platforms",
		Short: "Remove installed platforms",
		Run: func(cmd *cobra.Command, args []string) {
			platformCore, err := platform.New(client, logger)
			if err != nil {
				return
			}

			platformCore.Remove(args)
		},
	}

	return removeCmd
}

func getPlatformCommand() *cobra.Command {
	platCmd := &cobra.Command{
		Use: "platform",
		Long: "\nPlatform manager allowing addition and removal of specified " +
			"platforms either globally or at the project level. Default is " +
			"project level, use \"--global\" to manage global platforms. For " +
			"project specific platform commands see \"ardi help project platform\".",
		Short:   "Platform manager",
		Aliases: []string{"platforms"},
	}
	platCmd.AddCommand(getPlatformListCmd())
	platCmd.AddCommand(getPlatformAddCmd())
	platCmd.AddCommand(getPlatformRemoveCmd())
	return platCmd
}
