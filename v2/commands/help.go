package commands

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/mitchellh/go-wordwrap"
	"github.com/spf13/cobra"
)

// Help is custom wrapper around help output to make sure lines wrap to 80 char
func Help(cmd *cobra.Command, args []string) {
	wrapLen := uint(80)

	g, err := gocui.NewGui(gocui.OutputNormal)

	if err == nil {
		x, _ := g.Size()
		if x < 80 {
			wrapLen = uint(x)
		}
		g.Close()
	}

	fmt.Printf("%s\n\n", wordwrap.WrapString(cmd.Long, wrapLen))
	fmt.Println("Usage:")
	if cmd.Runnable() {
		fmt.Printf("%s\n\n", wordwrap.WrapString(cmd.CommandPath(), wrapLen))
	}
	if cmd.HasSubCommands() {
		use := cmd.CommandPath() + " [command]"
		fmt.Printf("%s\n\n", wordwrap.WrapString(use, wrapLen))
	}
	if cmd.HasSubCommands() {
		fmt.Println("Available Commands:")
		for _, c := range cmd.Commands() {
			line := fmt.Sprintf("%-*s %s", c.NamePadding(), c.Name(), c.Short)
			fmt.Println(wordwrap.WrapString(line, wrapLen))
		}
		fmt.Println("")
	}
	if cmd.HasFlags() {
		fmt.Printf("Flags:\n")
		fmt.Printf("%s\n\n", wordwrap.WrapString(cmd.Flags().FlagUsages(), wrapLen))
	}
	if cmd.Root().HasFlags() {
		fmt.Printf("Global Flags:\n")
		fmt.Printf("%s\n\n", wordwrap.WrapString(cmd.Root().Flags().FlagUsages(), wrapLen))
	}
	if cmd.HasSubCommands() {
		fmt.Println("Use \"ardi [command] --help\" for more information about a command.")
	}
	fmt.Println("")
}
