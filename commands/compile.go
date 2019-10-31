package commands

import (
	"github.com/robgonnella/ardi/v3/ardi"
	"github.com/robgonnella/ardi/v3/arguments"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getCompileCommand() *cobra.Command {
	var fqbn string
	var verbose bool
	var buildProps []string
	var showProps bool
	var compileCmd = &cobra.Command{
		Use:   "compile [sketch]",
		Short: "Compile specified sketch",
		Run: func(cmd *cobra.Command, args []string) {
			if !ardi.IsInitialized() {
				logger.Fatal("Ardi has not been initialized. Please run \"ardi init\" first")
			}
			if verbose {
				ardi.SetLogLevel(log.DebugLevel)
			} else {
				ardi.SetLogLevel(log.InfoLevel)
			}

			configFile := ardi.GlobalLibConfig
			if ardi.IsProjectDirectory() {
				configFile = ardi.LibConfig
			}

			conn, client, instance := ardi.StartDaemonAndGetConnection(configFile)
			defer conn.Close()

			sketchDir := "."
			sketchFile := ""
			if len(args) > 0 {
				sketchDir = args[0]
			}

			sketchDir, sketchFile = arguments.GetSketchParts(sketchDir)
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
					BuildProps: buildProps,
					ShowProps:  showProps,
				}
			}
			ardi.Compile(client, instance, target)
		},
	}
	compileCmd.Flags().StringVarP(&fqbn, "fqbn", "f", "", "Specify fully qualified board name")
	compileCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print all compilation logs")
	compileCmd.Flags().StringArrayVarP(&buildProps, "build-prop", "p", []string{}, "Specify build property to compiler")
	compileCmd.Flags().BoolVarP(&showProps, "show-props", "s", false, "Show all build properties (does not compile)")

	return compileCmd
}
