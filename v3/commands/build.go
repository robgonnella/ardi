package commands

import (
	"errors"

	"github.com/spf13/cobra"
)

func newBuildCmd(env *CommandEnv) *cobra.Command {
	var all bool
	var showProps bool

	var buildCmd = &cobra.Command{
		Use:   "build",
		Long:  "\nCompiles builds defined in ardi.json",
		Short: "Compiles builds defined in ardi.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireProjectInit(); err != nil {
				return err
			}

			ardiBuilds := env.ArdiCore.Config.GetBuilds()

			if len(ardiBuilds) == 0 {
				return errors.New("no builds defined in ardi.json")
			}

			if all {
				for name := range ardiBuilds {
					opts, err := env.ArdiCore.Config.GetCompileOpts(name)

					if err != nil {
						return err
					}

					opts.ShowProps = showProps

					if err := env.ArdiCore.Compiler.Compile(*opts); err != nil {
						return err
					}
				}
				return nil
			}

			for _, build := range args {
				opts, err := env.ArdiCore.Config.GetCompileOpts(build)

				if err != nil {
					return err
				}

				opts.ShowProps = showProps

				if err := env.ArdiCore.Compiler.Compile(*opts); err != nil {
					return err
				}
			}

			return nil
		},
	}

	buildCmd.Flags().BoolVarP(&all, "all", "a", false, "Compile all builds specified in ardi.json")
	buildCmd.Flags().BoolVarP(&showProps, "show-props", "s", false, "Show all build properties (does not compile)")

	return buildCmd
}
