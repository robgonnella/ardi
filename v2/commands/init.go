package commands

import (
	"strings"

	ardiInitCore "github.com/robgonnella/ardi/v2/core/ardi-init"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getInitCommand() *cobra.Command {
	var verbose bool
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Download, install, and update platforms (alias: ardi update)",
		Long: cyan("\nDownloads, installs, and updates all available platforms to\n" +
			"support a maximum number of boards. (alias: ardi update)"),
		Aliases: []string{"update"},
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()

			ardiInit, err := ardiInitCore.New(logger)
			if err != nil {
				return
			}
			defer ardiInit.RPC.Connection.Close()

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

			ardiInit.Initialize(platform, version)
		},
	}

	initCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print all logs")

	return initCmd
}
