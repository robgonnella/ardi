package commands

import (
	"github.com/robgonnella/ardi/v2/core/project"
	"github.com/robgonnella/ardi/v2/core/rpc"
	"github.com/robgonnella/ardi/v2/paths"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getListLibrariesCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "libraries",
		Short:   "List all project libraries specified in ardi.json",
		Aliases: []string{"libs"},
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
			projectCore, err := project.New(logger)
			if err != nil {
				return
			}
			projectCore.ListLibraries()
		},
	}
	return listCmd
}

func getListBuildsCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "builds",
		Short:   "List all project builds specified in ardi.json",
		Aliases: []string{"build"},
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
			projectCore, err := project.New(logger)
			if err != nil {
				return
			}
			projectCore.ListBuilds()
		},
	}
	return listCmd
}

func getProjectListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List project attributes saved in ardi.json",
	}
	listCmd.AddCommand(getListLibrariesCmd())
	listCmd.AddCommand(getListBuildsCmd())
	return listCmd
}

func getAddBuildCmd() *cobra.Command {
	var fqbn string
	var sketch string
	var buildProps []string
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add build config to ardi.json",
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
			projectCore, err := project.New(logger)
			if err != nil {
				return
			}
			projectCore.AddBuild(sketch, fqbn, buildProps)
		},
	}
	addCmd.Flags().StringVarP(&fqbn, "fqbn", "f", "", "Specify fully qualified board name")
	addCmd.Flags().StringVarP(&sketch, "sketch", "s", "", "Specify sketch directory")
	addCmd.Flags().StringArrayVarP(&buildProps, "build-prop", "p", []string{}, "Specify build property to compiler")
	return addCmd
}

func getProjectBuildCmd() *cobra.Command {
	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Compile builds specified in ardi.json",
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
			rpc, err := rpc.New(paths.ArdiDataConfig, logger)
			if err != nil {
				return
			}
			projectCore, err := project.New(logger)
			if err != nil {
				return
			}
			projectCore.Build(rpc, args)
		},
	}
	buildCmd.AddCommand(getAddBuildCmd())
	return buildCmd
}

func getProjectCommand() *cobra.Command {
	projectCmd := &cobra.Command{
		Use:   "project",
		Short: "Project related commands",
	}
	projectCmd.AddCommand(getProjectListCmd())
	projectCmd.AddCommand(getProjectBuildCmd())

	return projectCmd
}
