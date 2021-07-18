package commands

import (
	"github.com/robgonnella/ardi/v2/core"
	"github.com/spf13/cobra"
)

func getAttachCmd(env *CommandEnv) *cobra.Command {
	var port string
	var baud int
	var attachCmd = &cobra.Command{
		Use:   "attach",
		Short: "Attach and print board logs",
		Long:  "\nAttach and print board logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if port == "" {
				board, err := env.ArdiCore.Cli.GetTargetBoard("", "", true)
				if err != nil {
					return err
				}
				port = board.Port
			}
			serialPort := core.NewArdiSerialPort(port, baud, env.Logger)
			defer serialPort.Close()
			return serialPort.Watch()
		},
	}

	attachCmd.Flags().StringVarP(&port, "port", "p", "", "The port your arduino board is connected to")
	attachCmd.Flags().IntVarP(&baud, "baud", "b", 9600, "Specify baud rate of port")
	return attachCmd
}
