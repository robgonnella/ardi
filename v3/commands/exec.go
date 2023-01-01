package commands

import (
	"github.com/arduino/arduino-cli/cli"
	"github.com/robgonnella/ardi/v3/paths"
	"github.com/spf13/cobra"
)

func filter[T any](slice []T, f func(T) bool) []T {
	var n []T
	for _, e := range slice {
		if f(e) {
			n = append(n, e)
		}
	}
	return n
}

func newExecCmd(env *CommandEnv) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec",
		Short: "Execute arduino-cli command",
		Long:  "\nExecutes an arudion-cli command. All arduino-cli options are supported",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireProjectInit(); err != nil {
				return err
			}

			args = filter(args, func(a string) bool {
				return a != "arduino-cli"
			})
			args = append(args, "--config-file", paths.ArduinoCliProjectConfig)

			arduinoCliCmd := cli.NewCommand()
			arduinoCliCmd.SetArgs(args)

			return arduinoCliCmd.Execute()
		},
	}

	return cmd
}
