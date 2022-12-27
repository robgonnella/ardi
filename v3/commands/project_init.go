package commands

import (
	"github.com/robgonnella/ardi/v3/util"
	"github.com/spf13/cobra"
)

func newProjectInitCmd(env *CommandEnv) *cobra.Command {
	initCmd := &cobra.Command{
		Use:     "project-init",
		Aliases: []string{"init"},
		Short:   "Initialize directory as an ardi project",
		Long:    "\nInitialize directory as an ardi project",
		RunE: func(cmd *cobra.Command, args []string) error {
			return util.InitProjectDirectory()
		},
	}

	return initCmd
}
