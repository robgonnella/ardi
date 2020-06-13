package main

import (
	"log"

	"github.com/robgonnella/ardi/v2/commands"
	"github.com/spf13/cobra/doc"
)

func main() {
	rootCmd := commands.GetRootCmd()
	err := doc.GenMarkdownTree(rootCmd, "./docs")
	if err != nil {
		log.Fatal(err)
	}
}
