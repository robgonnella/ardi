package commands

import (
	"github.com/robgonnella/ardi/v2/core/compile"
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
		Long:  cyan("\nCompile specified sketch"),
		Short: "Compile specified sketch",
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
			if verbose {
				logger.SetLevel(log.DebugLevel)
			} else {
				logger.SetLevel(log.InfoLevel)
			}

			sketchDir := "."
			if len(args) > 0 {
				sketchDir = args[0]
			}

			compileCore, err := compile.New(logger)
			if err != nil {
				return
			}
			defer compileCore.RPC.Connection.Close()

			compileCore.Compile(sketchDir, fqbn, buildProps, showProps)
		},
	}
	compileCmd.Flags().StringVarP(&fqbn, "fqbn", "f", "", "Specify fully qualified board name")
	compileCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print all compilation logs")
	compileCmd.Flags().StringArrayVarP(&buildProps, "build-prop", "p", []string{}, "Specify build property to compiler")
	compileCmd.Flags().BoolVarP(&showProps, "show-props", "s", false, "Show all build properties (does not compile)")

	return compileCmd
}
