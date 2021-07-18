package main

import (
	"log"

	logrus "github.com/sirupsen/logrus"

	"github.com/robgonnella/ardi/v2/commands"
	"github.com/spf13/cobra/doc"
)

func main() {
	logger := logrus.New()
	cmdEnv := &commands.CommandEnv{Logger: logger}
	rootCmd := commands.GetRootCmd(cmdEnv)
	err := doc.GenMarkdownTree(rootCmd, "./docs")
	if err != nil {
		log.Fatal(err)
	}
}
