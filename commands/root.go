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
		Long: "A light wrapper around arduino-cli that offers a quick way to upload\n" +
			"sketches and watch logs from command line for a variety of arduino boards.",
	}
}

// Initialize adds all ardi commands to root and returns root command
func Initialize() *cobra.Command {
	rootCmd := getRootCommand()
	initCmd := getInitCommand()
	cleanCmd := getCleanCommand()
	goCmd := getGoCommand()
	compileCmd := getCompileCommand()
	libCmd := getLibCommand()
	platCmd := getPlatformCommand()
	boardCmd := getBoardCommand()

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(goCmd)
	rootCmd.AddCommand(compileCmd)
	rootCmd.AddCommand(libCmd)
	rootCmd.AddCommand(platCmd)
	rootCmd.AddCommand(boardCmd)

	return rootCmd
}
