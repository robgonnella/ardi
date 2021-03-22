package commands

import (
	"github.com/robgonnella/ardi/v2/types"
	"github.com/robgonnella/ardi/v2/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getUploadCmd() *cobra.Command {
	var attach bool
	var fqbn string
	var port string
	var uploadCmd = &cobra.Command{
		Use:   "upload [sketch-dir|build]",
		Short: "Upload pre-compiled sketch build to a connected board",
		Long: "\nUpload pre-compiled sketch build to a connected board. If " +
			"the sketch argument matches a user defined build in ardi.json, the " +
			"build values will be used to find the appropraite build to upload",
		RunE: func(cmd *cobra.Command, args []string) error {
			builds := ardiCore.Config.GetBuilds()

			build := "."
			if len(args) > 0 {
				build = args[0]
			}

			project := &types.Project{}
			var err error

			if ardiBuild, ok := builds[build]; ok {
				project.Directory = ardiBuild.Directory
				project.Sketch = ardiBuild.Sketch
				project.Baud = ardiBuild.Baud
				if fqbn == "" {
					fqbn = ardiBuild.FQBN
				}
			} else {
				project, err = util.ProcessSketch(build)
				if err != nil {
					return err
				}
			}

			board, err := ardiCore.GetTargetBoard(fqbn, port, true)
			if err != nil {
				return err
			}

			fields := logrus.Fields{
				"build":  project.Directory,
				"fqbn":   board.FQBN,
				"device": board.Port,
			}

			logger.WithFields(fields).Info("Uploading...")

			if err := ardiCore.Uploader.Upload(board, project.Directory); err != nil {
				return err
			}

			logger.Info("Upload successful")

			if attach {
				ardiCore.Uploader.Attach(board.Port, project.Baud, nil)
			}

			return nil
		},
	}
	uploadCmd.Flags().BoolVarP(&attach, "attach", "a", false, "Attach to board port and print logs")
	uploadCmd.Flags().StringVarP(&fqbn, "fqbn", "f", "", "The FQBN of the board you want to upload to")
	uploadCmd.Flags().StringVarP(&port, "port", "p", "", "The port your arduino board is connected to")

	return uploadCmd
}
