/*
Ardi is a command-line tool for compiling, uploading, and watching logs for
your usb connected arduino board. Ardi allows you to develop in an environment
you feel comfortable, without being forced to use arduino's web or desktop IDEs.

Ardi's "--watch" flag allows you to auto re-compile and upload on save, saving
you time and improving efficiency.

Ardi should work for all boards and platforms supported by arduino-cli.
Run "ardi init" to download all supported platforms and index files to ensure
maximum board support. To initialize only for a specific platform, run
"ardi init <platform_id>" or "ardi init <platform_id@version>". To see a list of
supported platforms and associated IDs, run "ardi platform list". To see a list
of all supported boards and their associated platforms and fqbns run
"ardi board list".
(Note board fqbn will only be filled in once platform is initialized)

Once initialized run "ardi go <sketch_dir> --watch --verbose" and ardi will try
to auto detect your board, compile your sketch, upload, watch for changes in
your sketch file, and re-compile and re-upload. You can also run,
"ardi compile <sketch_directory> --fqbn <board_fqbn>" to only compile and
skip uploading.

Ardi also includes a basic library manager. Run "ardi lib init" in your project
directory to initialize it as an ardi project directory. Once initialized,
you can use "ardi lib add <lib_name>" to add libraries,
"ardi lib remove <lib_name>", "ardi lib install" to install missing libraries
defined in ardi.json, and "ardi lib search <searchFilter>" to search existing
libraries.

Ardi stores all its platform data in "~/.ardi/" to avoid any conflicts with
existing "arduino-cli" installations.

Usage:
  ardi [command]

Available Commands:
  board       Board related commands
  clean       Delete all ardi global data
  compile     Compile specified sketch
  go          Compile and upload code to a connected arduino board
  help        Help about any command
  init        Download, install, and update platforms (alias: ardi update)
  lib         Library manager for ardi
  platform    Platform related commands

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
