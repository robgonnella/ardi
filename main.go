/*
ardi is a command-line tool for compiling, uploading code, and
watching logs for your usb connected arduino board. This allows you to
develop in an environment you feel comfortable in, without needing to
use arduino's web or desktop IDEs.

Usage:
  ardi [command]

Available Commands:
  clean       Delete all ardi data
  go          Compile and upload code to an arduino board
  help        Help about any command
  init        Download and install platforms

Flags:
  -h, --help   help for ardi

Use "ardi [command] --help" for more information about a command.
*/
package main

import "github.com/robgonnella/ardi/commands"

func main() {
	rootCmd := commands.Initialize()
	rootCmd.Execute()
}
