package commands

import (
	"strings"

	"github.com/spf13/cobra"
)

func getProjectInitCommand() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize directory as an ardi project",
		Long: "\nDownloads, installs, and updates specified platforms, or " +
			"all platforms if run with no arguments. This command will generate a " +
			"project level data directory for storing dependencies, and an " +
			"ardi.json file for managing builds and dependencies.",
		Aliases: []string{"update"},
		Run: func(cmd *cobra.Command, args []string) {
			ardiCore.Project.Init(port)
		},
	}

	return initCmd
}

func getProjectListPlatformCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "platform",
		Long:    "\nAdd platform(s) to project",
		Short:   "Add platform(s) to project",
		Aliases: []string{"platforms"},
		Run: func(cmd *cobra.Command, args []string) {
			ardiCore.Platform.ListInstalled()
		},
	}
	return listCmd
}

func getProjectListLibrariesCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "libraries",
		Long:    "\nList all project libraries specified in ardi.json",
		Short:   "List all project libraries specified in ardi.json",
		Aliases: []string{"libs"},
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("Libraries specified in ardi.json")
			ardiCore.Project.ListLibraries()
			logger.Info("Installed libraries")
			ardiCore.Lib.ListInstalled()
		},
	}
	return listCmd
}

func getProjectListBuildsCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "builds",
		Long:    "\nList all project builds specified in ardi.json",
		Short:   "List all project builds specified in ardi.json",
		Aliases: []string{"build"},
		Run: func(cmd *cobra.Command, args []string) {
			ardiCore.Project.ListBuilds(args)
		},
	}
	return listCmd
}

func getProjectListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Long:  "\nList project attributes saved in ardi.json",
		Short: "List project attributes saved in ardi.json",
	}
	listCmd.AddCommand(getProjectListPlatformCmd())
	listCmd.AddCommand(getProjectListLibrariesCmd())
	listCmd.AddCommand(getProjectListBuildsCmd())
	return listCmd
}

func getProjectAddPlatformCmd() *cobra.Command {
	addCmd := &cobra.Command{
		Use:     "platform",
		Long:    "\nAdd platform(s) to project",
		Short:   "Add platform(s) to project",
		Aliases: []string{"platforms"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 || strings.ToLower(args[0]) == "all" {
				ardiCore.Platform.AddAll()
				return
			}

			ardiCore.Platform.Add(args)
		},
	}
	return addCmd
}

func getProjectAddBuildCmd() *cobra.Command {
	var name string
	var platform string
	var boardURL string
	var fqbn string
	var sketch string
	var buildProps []string
	addCmd := &cobra.Command{
		Use:   "build",
		Long:  "\nAdd build config to project",
		Short: "Add build config to project",
		Run: func(cmd *cobra.Command, args []string) {
			ardiCore.Project.AddBuild(name, platform, boardURL, sketch, fqbn, buildProps)
		},
	}
	addCmd.Flags().StringVarP(&name, "name", "n", "", "Custom name for the build")
	addCmd.Flags().StringVarP(&fqbn, "fqbn", "f", "", "Specify fully qualified board name")
	addCmd.Flags().StringVarP(&sketch, "sketch", "s", "", "Path to .ino file or sketch directory")
	addCmd.Flags().StringVarP(&platform, "platform", "m", "", "Platform for this build \"package:architecture@version\" (optional)")
	addCmd.Flags().StringVarP(&boardURL, "board-url", "u", "", "Custom board url (optional)")
	addCmd.Flags().StringArrayVarP(&buildProps, "build-prop", "p", []string{}, "Specify build property to compiler (optional)")
	return addCmd
}

func getProjectAddLibCmd() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "lib",
		Long:  "\nAdd libraries to project",
		Short: "Add libraries to project\\e[0m",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, l := range args {
				name, vers, err := ardiCore.Lib.Add(l)
				if err != nil {
					logger.WithError(err).Error("Failed to install all libraries")
					return
				}
				if err := ardiCore.Project.AddLibrary(name, vers); err != nil {
					logger.WithError(err).Error("Failed to save libary to ardi.json")
					return
				}
			}
		},
	}
	return addCmd
}

func getProjectAddCmd() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Long:  "\nAdd libraries and builds to project",
		Short: "Add libraries and builds to project",
	}
	addCmd.AddCommand(getProjectAddPlatformCmd())
	addCmd.AddCommand(getProjectAddBuildCmd())
	addCmd.AddCommand(getProjectAddLibCmd())
	return addCmd
}

func getProjectRemovePlatformCmd() *cobra.Command {
	removeCmd := &cobra.Command{
		Use:     "platform",
		Long:    "\nRemove platform(s) from project",
		Short:   "Remove platform(s) from project",
		Aliases: []string{"platforms"},
		Run: func(cmd *cobra.Command, args []string) {
			ardiCore.Platform.Remove(args)
		},
	}
	return removeCmd
}

func getProjectRemoveBuildCmd() *cobra.Command {
	removeCmd := &cobra.Command{
		Use:     "build",
		Long:    "\nRemove build config from project",
		Short:   "Remove build config from project",
		Aliases: []string{"builds"},
		Run: func(cmd *cobra.Command, args []string) {
			ardiCore.Project.RemoveBuild(args)
		},
	}
	return removeCmd
}

func getProjectRemoveLibCmd() *cobra.Command {
	removeCmd := &cobra.Command{
		Use:   "lib",
		Long:  "\nRemove libraries from project",
		Short: "Remove libraries from project",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, l := range args {
				if err := ardiCore.Lib.Remove(l); err != nil {
					logger.WithError(err).Errorf("Failed to uninstall library: %s", l)
					return
				}
				if err := ardiCore.Lib.Remove(l); err != nil {
					logger.WithError(err).Error("Failed to remove library from ardi.json: %s", l)
					return
				}
			}
		},
	}
	return removeCmd
}

func getProjectRemoveCmd() *cobra.Command {
	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove libraries and builds from project",
		Long:  "\nRemove libraries and builds from project",
	}
	removeCmd.AddCommand(getProjectRemovePlatformCmd())
	removeCmd.AddCommand(getProjectRemoveBuildCmd())
	removeCmd.AddCommand(getProjectRemoveLibCmd())
	return removeCmd
}

func getProjectBuildCmd() *cobra.Command {
	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Compile builds specified in ardi.json",
		Long:  "\nCompile builds specified in ardi.json",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				ardiCore.Project.BuildAll()
				return
			}

			ardiCore.Project.BuildList(args)
		},
	}
	return buildCmd
}

func getProjectInstallCmd() *cobra.Command {
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install all libraries specified in ardi.json",
		Long:  "\nInstall all libraries specified in ardi.json",
		Run: func(cmd *cobra.Command, args []string) {
			for _, l := range ardiCore.Project.GetLibraries() {
				lib, vers, err := ardiCore.Lib.Add(l)
				if err != nil {
					logger.WithError(err).Errorf("Failed to install %s", l)
				} else {
					logger.Infof("Installed: %s@%s", lib, vers)
				}
			}
		},
	}
	return installCmd
}

func getProjectCommand() *cobra.Command {
	projectCmd := &cobra.Command{
		Use:   "project",
		Short: "Project manager",
		Long: "\nProject manager allowing you to store versioned dependencies " +
			"and build configurations for each project.\n" +
			"See \"ardi help project init\" for more details",
	}
	projectCmd.AddCommand(getProjectInitCommand())
	projectCmd.AddCommand(getProjectListCmd())
	projectCmd.AddCommand(getProjectAddCmd())
	projectCmd.AddCommand(getProjectRemoveCmd())
	projectCmd.AddCommand(getProjectBuildCmd())
	projectCmd.AddCommand(getProjectInstallCmd())

	return projectCmd
}
