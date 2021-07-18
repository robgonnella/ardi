package commands

import (
	"github.com/spf13/cobra"
)

func getAddPlatformCmd(env *CommandEnv) *cobra.Command {
	addCmd := &cobra.Command{
		Use:     "platforms",
		Long:    "\nAdd platform(s) to project",
		Short:   "Add platform(s) to project",
		Aliases: []string{"platform"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, p := range args {
				env.Logger.Infof("Adding platform: %s", p)
				installed, vers, err := env.ArdiCore.Platform.Add(p)
				if err != nil {
					env.Logger.WithError(err).Errorf("Failed to add arduino platform %s", p)
					return err
				}
				if err := env.ArdiCore.Config.AddPlatform(installed, vers); err != nil {
					return err
				}
				env.Logger.Info("Updated config")
			}
			return nil
		},
	}
	return addCmd
}

func getAddBuildCmd(env *CommandEnv) *cobra.Command {
	var name string
	var fqbn string
	var sketch string
	var baud int
	var buildProps []string
	addCmd := &cobra.Command{
		Use:   "build",
		Long:  "\nAdd build config to project",
		Short: "Add build config to project",
		RunE: func(cmd *cobra.Command, args []string) error {
			return env.ArdiCore.Config.AddBuild(name, sketch, fqbn, baud, buildProps)
		},
	}
	addCmd.Flags().StringVarP(&name, "name", "n", "", "Custom name for the build")
	addCmd.Flags().StringVarP(&fqbn, "fqbn", "f", "", "Specify fully qualified board name")
	addCmd.Flags().StringVarP(&sketch, "sketch", "s", "", "Path to .ino file or sketch directory")
	addCmd.Flags().IntVarP(&baud, "baud", "b", 0, "Specify baud rate for build")
	addCmd.Flags().StringArrayVarP(&buildProps, "build-prop", "p", []string{}, "Specify build property to compiler")
	addCmd.MarkFlagRequired("name")
	addCmd.MarkFlagRequired("fqbn")
	addCmd.MarkFlagRequired("sketch")

	return addCmd
}

func getAddLibCmd(env *CommandEnv) *cobra.Command {
	addCmd := &cobra.Command{
		Use:     "libraries",
		Long:    "\nAdd libraries to project",
		Short:   "Add libraries to project",
		Aliases: []string{"libs", "lib", "library"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, l := range args {
				env.Logger.Infof("Adding library: %s", l)
				name, vers, err := env.ArdiCore.Lib.Add(l)
				if err != nil {
					env.Logger.WithError(err).Errorf("Failed to install library %s", l)
					return err
				}
				env.Logger.Infof("Successfully installed %s@%s", name, vers)
				if err := env.ArdiCore.Config.AddLibrary(name, vers); err != nil {
					env.Logger.WithError(err).Error("Failed to save libary to ardi.json")
					return err
				}
			}
			return nil
		},
	}
	return addCmd
}

func getAddBoardURLCmd(env *CommandEnv) *cobra.Command {
	addCmd := &cobra.Command{
		Use:     "board-url",
		Long:    "\nAdd board urls to project",
		Short:   "Add board urls to project",
		Aliases: []string{"board-urls"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, u := range args {
				if err := env.ArdiCore.Config.AddBoardURL(u); err != nil {
					return err
				}
				if err := env.ArdiCore.CliConfig.AddBoardURL(u); err != nil {
					return err
				}
			}
			return nil
		},
	}
	return addCmd
}

func getAddCmd(env *CommandEnv) *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Long:  "\nAdd project dependencies",
		Short: "Add project dependencies",
	}
	addCmd.AddCommand(getAddPlatformCmd(env))
	addCmd.AddCommand(getAddBuildCmd(env))
	addCmd.AddCommand(getAddLibCmd(env))
	addCmd.AddCommand(getAddBoardURLCmd(env))
	return addCmd
}
