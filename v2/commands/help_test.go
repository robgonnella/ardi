package commands_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/robgonnella/ardi/v2/commands"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestHelpCommand(t *testing.T) {
	t.Run("should wrap lines to 80 char", func(st *testing.T) {
		originalOut := os.Stdout
		r, w, _ := os.Pipe()

		st.Cleanup(func() {
			os.Stdout = originalOut
			w.Close()
			r.Close()
		})

		os.Stdout = w
		cmd := &cobra.Command{
			Use:     "somecmd [args]",
			Short:   "Longer than 80char Longer than 80char Longer than 80char Longer than 80char Longer than 80char Longer than 80char Longer than 80char Longer than 80char",
			Long:    "\nLonger than 80char Longer than 80char Longer than 80char Longer than 80char Longer than 80char Longer than 80char Longer than 80char Longer than 80char Longer than 80char",
			Aliases: []string{"many", "aliases", "so", "what", "about", "even", "more", "how", "can", "we", "get", "words", "to", "continue", "forever", "end", "ever", "without", "repeating"},
			Run:     func(cmd *cobra.Command, args []string) {},
		}

		commands.Help(cmd, []string{})
		w.Close()

		var buf bytes.Buffer
		io.Copy(&buf, r)
		r.Close()

		split := strings.Split(buf.String(), "\n")
		for _, line := range split {
			assert.True(st, len(line) <= 80)
		}
	})
}
