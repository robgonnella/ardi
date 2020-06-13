/*
Ardi is a command-line tool for ardiuno that enables you to properly version and
manage project builds, and provides tools to help facilitate the development
process.

Things ardi can fo for you:

• Manage versioned platforms and libraries on a per-project basis
• Store user defined build configurations with a mechanism for easily running
  consistent and repeatable builds.
• Enable running your builds in a CI pipeline
• Compile and upload to an auto discovered connected board
• Watche a sketch for changes and auto recompil / reupload to a connected board
• Print various info about platforms and boards
• Search and print available libraries and versions

Ardi should work for all boards and platforms supported by arduino-cli.
*/
package main

import (
	"github.com/robgonnella/ardi/v2/commands"
)

func main() {
	rootCmd := commands.GetRootCmd()
	rootCmd.Execute()
}
