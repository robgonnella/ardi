/*
Ardi is a command-line tool for ardiuno that enables you to properly version and
manage project builds, and provides tools to help facilitate the development
process.

Things ardi can fo for you:

• Manage versioned platforms and libraries on a per-project basis

• Store user defined build config for consistent and repeatable builds.

• Enable running your builds in a CI pipeline

• Compile and upload to an auto discovered connected board

• Watch sketch for changes and auto recompile / reupload

• Print various info about platforms and boards

• Search and print available libraries and versions

Ardi should work for all boards and platforms supported by arduino-cli.
*/
package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/robgonnella/ardi/v2/commands"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := log.New()
	rootCmd := commands.GetRootCmd(logger)
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		logger.WithError(err).Error("Command failed")
		os.Exit(1)
	}
}
