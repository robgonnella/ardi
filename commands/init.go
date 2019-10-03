package commands

import (
	"github.com/robgonnella/ardi/ardi"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getInitCommand() *cobra.Command {
	var verbose bool
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Download, install, and update platforms (alias: ardi update)",
		Long: "Downloads, installs, and updates all available platforms to\n" +
			"support a maximum number of boards. (alias: ardi update)",
		Aliases: []string{"update"},
		Run: func(cmd *cobra.Command, args []string) {
			if verbose {
				ardi.SetLogLevel(log.DebugLevel)
			} else {
				ardi.SetLogLevel(log.InfoLevel)
			}
			platform := ""
			version := ""
			if len(args) > 0 {
				platform = args[0]
			}
			if len(args) > 1 {
				version = args[1]
			}
			logger.Info("Initializing. This may take some time...")
			ardi.Initialize(platform, version)
			logger.Info("Successfully initialized!")
		},
	}

	initCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print all logs")

	return initCmd
}
