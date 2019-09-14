/*
Ardi is a command-line tool for compiling, uploading, and watching logs for
your usb connected arduino board. Ardi allows you to develop in an environment
you feel comfortable, without being forced to use arduino's web or desktop IDEs.

Ardi's `--watch` flag allows you to auto re-compile and upload on save, saving
you time and improving efficiency.

Ardi should work for all boards and platforms supported by arduino-cli.
Run `ardi init` to download all supported platforms and indexes to ensure
maximum board support.

Once initialized run `ardi go <sketch_dir> --watch --verbose` and ardi will try
to auto detect your board, compile your sketch, upload, watch for changes in
your sketch file, and re-compile and re-upload.

Ardi stores all its data in a `.ardi` directory in the users home directory
to avoid any conflicts with existing `arduino-cli` installations.

Usage:
  ardi [command]

Description:
  A light wrapper around arduino-cli that offers a quick way to upload
  sketches and watch logs from command line for a variety of arduino boards.

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
