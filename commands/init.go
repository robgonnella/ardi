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
		Short: "Download and install platforms",
		Long: "Downloads and installs all available platforms to support\n" +
			"a maximum number of boards.",
		Run: func(cmd *cobra.Command, args []string) {
			if verbose {
				ardi.SetLogLevel(log.DebugLevel)
			} else {
				ardi.SetLogLevel(log.InfoLevel)
			}
			logger.Info("Initializing. This may take some time...")
			conn, _, _ := ardi.Initialize()
			defer conn.Close()
			logger.Info("Successfully initialized!")
		},
	}

	initCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "print all logs")

	return initCmd
}
