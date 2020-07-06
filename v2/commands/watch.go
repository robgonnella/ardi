package commands

import (
	"github.com/spf13/cobra"
)

func process(sketchDir string, buildProps []string) error {
	if err := ardiCore.Watch.Init(port, sketchDir, buildProps); err != nil {
		logger.WithError(err).Error("Failed to initialize ardi watch core")
		return err
	}

	if err := ardiCore.Watch.Compile(); err != nil {
		logger.WithError(err).Error("Failed to compile")
		return err
	}

	if err := ardiCore.Watch.Upload(); err != nil {
		logger.WithError(err).Error("Failed to upload")
		return err
	}

	ardiCore.Watch.WatchSketch()
	return nil
}

func getGoCommand() *cobra.Command {
	var buildProps []string

	var goCmd = &cobra.Command{
		Use:   "watch [sketch]",
		Short: "Compile, upload, and watch",
		Long: "\nCompile and upload code to an arduino board. Simply pass the " +
			"directory containing the .ino file as the first argument. Ardi will " +
			"automatically watch your sketch file for changes and auto re-compile " +
			"& re-upload for you. Baud will be automatically be detected from " +
			"sketch file.",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sketchDir := args[0]
			return process(sketchDir, buildProps)
		},
	}

	goCmd.Flags().StringArrayVarP(&buildProps, "build-prop", "p", []string{}, "Specify build property to compiler")

	return goCmd
}
