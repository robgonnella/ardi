package commands

import (
	"errors"

	"github.com/spf13/cobra"
)

func getPlatformListCmd() *cobra.Command {
	var all bool
	var installed bool
	listCmd := &cobra.Command{
		Use:   "list",
		Long:  "\nList platforms",
		Short: "List platforms",
		RunE: func(cmd *cobra.Command, args []string) error {
			if all || (!all && !installed) {
				if err := ardiCore.Platform.ListAll(); err != nil {
					logger.WithError(err).Error("Failed to list arduino platforms")
					return err
				}
				return nil
			}

			if installed {
				if err := ardiCore.Platform.ListInstalled(); err != nil {
					logger.WithError(err).Error("Failed to list installed arduino platforms")
					return err
				}
			}
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			if all {
				if err := ardiCore.Platform.AddAll(); err != nil {
					logger.WithError(err).Error("Failed to install arduino platforms")
					return err
				}
				return nil
			}
			if len(args) == 0 {
				msg := "No platforms specified"
				logger.Error(msg)
				return errors.New(msg)
			}
			for _, p := range args {
				if _, _, err := ardiCore.Platform.Add(p); err != nil {
					logger.WithError(err).Error("Failed to install arduino platforms")
					return err
				}
			}
			return nil
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
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, p := range args {
				if _, err := ardiCore.Platform.Remove(p); err != nil {
					logger.WithError(err).Errorf("Failed to remove arduino platform %s", p)
					return err
				}
			}
			return nil
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
