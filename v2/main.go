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

	log "github.com/sirupsen/logrus"

	"github.com/robgonnella/ardi/v2/cli-wrapper"
	"github.com/robgonnella/ardi/v2/commands"
	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/util"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := log.New()

	ardiConfig, svrSettings := util.GetAllSettings()
	cliSettingsPath := util.GetCliSettingsPath()

	if util.IsProjectDirectory() {
		if err := util.WriteAllSettings(ardiConfig, svrSettings); err != nil {
			logger.WithError(err).Fatal("Failed to write settings files")
		}
	}

	coreOpts := core.NewArdiCoreOpts{
		Ctx:                ctx,
		Logger:             logger,
		CliSettingsPath:    cliSettingsPath,
		ArdiConfig:         *ardiConfig,
		ArduinoCliSettings: *svrSettings,
	}

	arduinoCli := cli.NewArduinoCli()
	withArduinoCli := core.WithArduinoCli(arduinoCli)
	ardiCore := core.NewArdiCore(coreOpts, withArduinoCli)

	env := &commands.CommandEnv{
		ArdiCore: ardiCore,
		Logger:   logger,
	}

	rootCmd := commands.GetRootCmd(env)
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		logger.WithError(err).Fatal("Command failed")
	}
}
