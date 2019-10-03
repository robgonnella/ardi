package commands

import (
	"github.com/robgonnella/ardi/ardi"
	"github.com/robgonnella/ardi/arguments"
	"github.com/spf13/cobra"
)

func getCompileCommand() *cobra.Command {
	var fqbn string
	var compileCmd = &cobra.Command{
		Use:   "compile [sketch]",
		Short: "Compile specified sketch",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if !ardi.IsInitialized() {
				logger.Fatal("Ardi has not been initialized. Please run \"ardi init\" first")
			}

			configFile := ardi.GlobalLibConfig
			if ardi.IsProjectDirectory() {
				configFile = ardi.LibConfig
			}

			conn, client, instance := ardi.StartDaemonAndGetConnection(configFile)
			defer conn.Close()

			sketchDir, sketchFile := arguments.GetSketchParts(args[0])
			list := ardi.GetTargetList(client, instance, sketchDir, sketchFile, 9600)
			var target *ardi.TargetInfo

			if len(list) > 0 {
				t := ardi.GetTargetInfo(list)
				target = &t
			} else if fqbn == "" {
				target = &ardi.TargetInfo{
					SketchDir:  sketchDir,
					SketchFile: sketchFile,
					FQBN:       ardi.GetDesiredBoard(client, instance),
				}
			} else {
				target = &ardi.TargetInfo{
					SketchDir:  sketchDir,
					SketchFile: sketchFile,
					FQBN:       fqbn,
				}
			}
			ardi.Compile(client, instance, target)
		},
	}
	compileCmd.Flags().StringVarP(&fqbn, "fqbn", "f", "", "specify fully qualified board name")
	return compileCmd
}
