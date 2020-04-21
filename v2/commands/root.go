package commands

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var logger = log.New()

func getRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "ardi",
		Short: "Ardi uploads sketches and prints logs for a variety of arduino boards.",
		Long: cyan("\nA light wrapper around arduino-cli that offers a quick way to upload\n" +
			"sketches and watch logs from command line for a variety of arduino boards."),
	}
}

// Initialize adds all ardi commands to root and returns root command
func Initialize(version string) *cobra.Command {
	rootCmd := getRootCommand()

	rootCmd.AddCommand(getVersionCommand(version))
	rootCmd.AddCommand(getInitCommand())
	rootCmd.AddCommand(getCleanCommand())
	rootCmd.AddCommand(getGoCommand())
	rootCmd.AddCommand(getCompileCommand())
	rootCmd.AddCommand(getLibCommand())
	rootCmd.AddCommand(getPlatformCommand())
	rootCmd.AddCommand(getBoardCommand())
	rootCmd.AddCommand(getProjectCommand())

	return rootCmd
}
