package commands

import (
	"strings"

	"github.com/robgonnella/ardi/v2/core/lib"
	"github.com/robgonnella/ardi/v2/core/project"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getProjectInitCommand() *cobra.Command {
	var verbose bool
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize directory as an ardi project",
		Long: cyan("\nDownloads, installs, and updates specified platforms, or\n" +
			"all platforms if not specified, creates project data directory, and\n" +
			"creates project level ardi.json"),
		Aliases: []string{"update"},
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()

			projectCore, err := project.New(logger)
			if err != nil {
				return
			}
			defer projectCore.Client.Connection.Close()

			if verbose {
				logger.SetLevel(log.DebugLevel)
			} else {
				logger.SetLevel(log.InfoLevel)
			}

			platform := ""
			version := ""
			if len(args) > 0 {
				platParts := strings.Split(args[0], "@")
				if len(platParts) > 0 {
					platform = platParts[0]
				}
				if len(platParts) > 1 {
					version = platParts[1]
				}
			}

			projectCore.Init(platform, version)
		},
	}

	initCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print all logs")

	return initCmd
}

func getProjectListLibrariesCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "libraries",
		Long:    cyan("\nList all project libraries specified in ardi.json"),
		Short:   "List all project libraries specified in ardi.json",
		Aliases: []string{"libs"},
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
			projectCore, err := project.New(logger)
			if err != nil {
				return
			}
			defer projectCore.Client.Connection.Close()
			projectCore.ListLibraries()
		},
	}
	return listCmd
}

func getProjectListBuildsCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "builds",
		Long:    cyan("\nList all project builds specified in ardi.json"),
		Short:   "List all project builds specified in ardi.json",
		Aliases: []string{"build"},
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
			projectCore, err := project.New(logger)
			if err != nil {
				return
			}
			defer projectCore.Client.Connection.Close()
			projectCore.ListBuilds()
		},
	}
	return listCmd
}

func getProjectListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Long:  cyan("\nList project attributes saved in ardi.json"),
		Short: "List project attributes saved in ardi.json",
	}
	listCmd.AddCommand(getProjectListLibrariesCmd())
	listCmd.AddCommand(getProjectListBuildsCmd())
	return listCmd
}

func getProjectAddBuildCmd() *cobra.Command {
	var name string
	var fqbn string
	var path string
	var buildProps []string
	addCmd := &cobra.Command{
		Use:   "build",
		Long:  cyan("\nAdd build config to project"),
		Short: "Add build config to project",
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
			projectCore, err := project.New(logger)
			if err != nil {
				return
			}
			defer projectCore.Client.Connection.Close()
			projectCore.AddBuild(name, path, fqbn, buildProps)
		},
	}
	addCmd.Flags().StringVarP(&name, "name", "n", "", "Custom name for the build")
	addCmd.Flags().StringVarP(&fqbn, "fqbn", "f", "", "Specify fully qualified board name")
	addCmd.Flags().StringVarP(&path, "path", "i", "", "Path to .ino file or sketch directory")
	addCmd.Flags().StringArrayVarP(&buildProps, "build-prop", "p", []string{}, "Specify build property to compiler")
	return addCmd
}

func getProjectAddLibCmd() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "lib",
		Long:  cyan("\nAdd libraries to project"),
		Short: "Add libraries to project\\e[0m",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
			libCore, err := lib.New(logger)
			if err != nil {
				return
			}
			defer libCore.Client.Connection.Close()
			libCore.Add(args)
		},
	}
	return addCmd
}

func getProjectAddCmd() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Long:  cyan("\nAdd libraries and builds to project"),
		Short: "Add libraries and builds to project",
	}
	addCmd.AddCommand(getProjectAddBuildCmd())
	addCmd.AddCommand(getProjectAddLibCmd())
	return addCmd
}

func getProjectRemoveBuildCmd() *cobra.Command {
	removeCmd := &cobra.Command{
		Use:     "build",
		Long:    cyan("\nRemove build config from project"),
		Short:   "Remove build config from project",
		Aliases: []string{"builds"},
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
			projectCore, err := project.New(logger)
			if err != nil {
				return
			}
			defer projectCore.Client.Connection.Close()
			projectCore.RemoveBuild(args)
		},
	}
	return removeCmd
}

func getProjectRemoveLibCmd() *cobra.Command {
	removeCmd := &cobra.Command{
		Use:   "lib",
		Long:  cyan("\nRemove libraries from project"),
		Short: "Remove libraries from project",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
			libCore, err := lib.New(logger)
			if err != nil {
				return
			}
			defer libCore.Client.Connection.Close()
			libCore.Remove(args)
		},
	}
	return removeCmd
}

func getProjectRemoveCmd() *cobra.Command {
	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove libraries and builds from project",
		Long:  cyan("\nRemove libraries and builds from project"),
	}
	removeCmd.AddCommand(getProjectRemoveBuildCmd())
	removeCmd.AddCommand(getProjectRemoveLibCmd())
	return removeCmd
}

func getProjectBuildCmd() *cobra.Command {
	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Compile builds specified in ardi.json",
		Long:  cyan("\nCompile builds specified in ardi.json"),
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
			projectCore, err := project.New(logger)
			if err != nil {
				return
			}
			defer projectCore.Client.Connection.Close()
			projectCore.Build(args)
		},
	}
	return buildCmd
}

func getProjectCommand() *cobra.Command {
	projectCmd := &cobra.Command{
		Use:   "project",
		Short: "Project related commands",
		Long:  cyan("\nProject related commands"),
	}
	projectCmd.AddCommand(getProjectInitCommand())
	projectCmd.AddCommand(getProjectListCmd())
	projectCmd.AddCommand(getProjectAddCmd())
	projectCmd.AddCommand(getProjectRemoveCmd())
	projectCmd.AddCommand(getProjectBuildCmd())

	return projectCmd
}
