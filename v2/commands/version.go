package commands

import (
	"github.com/robgonnella/ardi/v2/version"

	"github.com/spf13/cobra"
)

func getVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Long:  "\nPrints current version of ardi",
		Short: "Prints current version of ardi",
		Run: func(cmd *cobra.Command, args []string) {
			ardiVersion := version.VERSION
			arduinoCliVersion := ardiCore.RPCClient.ClientVersion()
			logger.Infoln("")
			logger.Infof("ardi: v%s", ardiVersion)
			logger.Infof("arduino-cli: %s", arduinoCliVersion)
			logger.Infoln("")
		},
	}
}
