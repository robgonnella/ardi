package commands

import (
	"github.com/robgonnella/ardi/v2/version"

	"github.com/spf13/cobra"
)

func getVersionCmd(env *CommandEnv) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Long:  "\nPrints current version of ardi",
		Short: "Prints current version of ardi",
		Run: func(cmd *cobra.Command, args []string) {
			ardiVersion := version.VERSION
			arduinoCliVersion := env.ArdiCore.Cli.ClientVersion()
			env.Logger.Infoln("")
			env.Logger.Infof("ardi: v%s", ardiVersion)
			env.Logger.Infof("arduino-cli: %s", arduinoCliVersion)
			env.Logger.Infoln("")
		},
	}
}
