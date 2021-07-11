package commands

import (
	"errors"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
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
			defer ardiCore.Uploader.Detach() // noop if not attached

			builds := ardiCore.Config.GetBuilds()

			build := "."
			if len(args) > 0 {
				build = args[0]
			}

			// Ignore errors here as user may have provided fqbn via build to mitigate
			// custom boards that don't show up via auto detect for some reason
			board, _ := ardiCore.Cli.GetTargetBoard(fqbn, port, true)

			project := &types.Project{}
			var err error

			if ardiBuild, ok := builds[build]; ok {
				project.Directory = ardiBuild.Directory
				project.Sketch = ardiBuild.Sketch
				project.Baud = ardiBuild.Baud
				if board == nil {
					board = &cli.BoardWithPort{FQBN: ardiBuild.FQBN, Port: port}
				}
			} else {
				project, err = util.ProcessSketch(build)
				if err != nil {
					return err
				}
			}

			if board == nil || board.FQBN == "" || board.Port == "" {
				return errors.New("no connected boards detected")
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
