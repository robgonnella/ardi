package commands

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/mitchellh/go-wordwrap"
	"github.com/spf13/cobra"
)

func help(cmd *cobra.Command, args []string) {
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
	fmt.Printf("Usage:\n")
	fmt.Printf("%s\n\n", wordwrap.WrapString(cmd.Use, wrapLen))
	fmt.Printf("Flags:\n")
	fmt.Printf("%s\n\n", wordwrap.WrapString(cmd.Flags().FlagUsages(), wrapLen))
	fmt.Printf("Global Flags:\n")
	fmt.Printf("%s\n\n", wordwrap.WrapString(cmd.Root().Flags().FlagUsages(), wrapLen))
}
