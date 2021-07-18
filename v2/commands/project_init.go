package commands

import (
	"github.com/robgonnella/ardi/v2/util"
	"github.com/spf13/cobra"
)

func getProjectInitCmd(env *CommandEnv) *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "project-init",
		Short: "Initialize directory as an ardi project",
		Long:  "\nInitialize directory as an ardi project",
		RunE: func(cmd *cobra.Command, args []string) error {
			return util.InitProjectDirectory()
		},
	}

	return initCmd
}
