package commands

import (
	"errors"

	"github.com/spf13/cobra"
)

func getCompileCmd(env *CommandEnv) *cobra.Command {
	var all bool
	var fqbn string
	var buildProps []string
	var showProps bool
	var watch bool
	var compileCmd = &cobra.Command{
		Use: "compile [sketch|build(s)]",
		Long: "\nCompile sketches and builds for specified boards. When " +
			"compileing for a sketch, you must provide the board FQBN. If left " +
			"unspecified, a list of available choices will be be printed. If the " +
			"sketch argument matches a user defined build in ardi.json, the values " +
			"defined in build will be used to compile",
		Short:   "Compile specified sketch or build(s)",
		Aliases: []string{"build"},
		RunE: func(cmd *cobra.Command, args []string) error {
			defer env.ArdiCore.Compiler.StopWatching() // noop if not watching

			board, _ := env.ArdiCore.Cli.GetTargetBoard(fqbn, "", true)

			if board != nil {
				fqbn = board.FQBN
			}

			if all {
				if watch {
					return errors.New("cannot watch all builds. You can only watch one build at a time")
				}

				ardiBuilds := env.ArdiCore.Config.GetBuilds()

				for name := range ardiBuilds {
					opts, err := env.ArdiCore.Config.GetCompileOpts(name)
					if err != nil {
						return err
					}
					if err := env.ArdiCore.Compiler.Compile(*opts); err != nil {
						return err
					}
				}
				return nil
			}

			opts, err := env.ArdiCore.GetCompileOptsFromArgs(fqbn, buildProps, showProps, args)
			if err != nil {
				return err
			}

			if len(opts) > 1 && watch {
				return errors.New("cannot specifify watch with mutiple builds. You can only watch one build at a time")
			}

			for _, compileOpts := range opts {
				if err := env.ArdiCore.Compiler.Compile(*compileOpts); err != nil {
					return err
				}
			}

			if watch {
				return env.ArdiCore.Compiler.WatchForChanges(*opts[0])
			}

			return nil
		},
	}
	compileCmd.Flags().BoolVarP(&all, "all", "a", false, "Compile all builds specified in ardi.json")
	compileCmd.Flags().StringVarP(&fqbn, "fqbn", "f", "", "Specify fully qualified board name")
	compileCmd.Flags().StringArrayVarP(&buildProps, "build-prop", "p", []string{}, "Specify build property to compiler")
	compileCmd.Flags().BoolVarP(&showProps, "show-props", "s", false, "Show all build properties (does not compile)")
	compileCmd.Flags().BoolVarP(&watch, "watch", "w", false, "Watch sketch file for changes and recompile")

	return compileCmd
}
