package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	arduino "github.com/arduino/arduino-cli/cli"
	log "github.com/sirupsen/logrus"
)

var cli = arduino.ArduinoCli
var logger = log.New()

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func main() {
	var boards string
	var fqbn string
	var device string
	var sketch string
	var err error
	baud := "9600"

	if len(os.Args) == 1 {
		logger.WithError(errors.New("Missing sketch arguemnet")).Fatal("Must provide a sketch name as an argument to upload")
	}

	sketch = os.Args[1]
	sketch = strings.Replace(sketch, "sketches/", "", 1)
	sketch = fmt.Sprintf("sketches/%s", sketch)

	if len(os.Args) == 3 {
		baud = os.Args[2]
	}

	cli.SetArgs([]string{"core", "update-index"})
	if err = cli.Execute(); err != nil {
		logger.WithError(err).Fatal("Failed to update index")
	}

	cli.SetArgs([]string{"core", "install", "arduino:avr"})
	if err = cli.Execute(); err != nil {
		logger.WithError(err).Fatal("Failed to install arduino:avr core")
	}

	out := os.Stdout
	buf := new(bytes.Buffer)
	r, w, _ := os.Pipe()
	os.Stdout = w
	cli.SetArgs([]string{"board", "list"})
	if err = cli.Execute(); err != nil {
		os.Stdout = out
		logger.WithError(err).Fatal("Failed to get board list")
	}
	os.Stdout = out
	w.Close()
	buf.ReadFrom(r)
	r.Close()

	boards = buf.String()

	printableList := strings.SplitAfterN(boards, "\n", -1)
	printableList = Filter(printableList, func(s string) bool {
		return !strings.Contains(s, "Unknown") && s != ""
	})

	list := strings.Split(string(boards), "\n")
	list = Filter(list, func(s string) bool {
		return !strings.Contains(s, "Unknown") && !strings.Contains(s, "Board Name") && s != ""
	})

	if len(list) == 0 {
		err = errors.New("No boards detected")
		logger.WithError(err).Fatal("Cannot upload sketch")
	} else if len(list) == 1 {
		board := strings.Split(list[0], " ")
		device = board[0]
		fqbn = board[len(board)-1]
	} else {
		for i, line := range printableList {
			if i == 0 {
				fmt.Printf("\n   %s", line)
			} else {
				fmt.Printf("%d: %s", i-1, line)
			}
		}
		var index int
		fmt.Print("\nEnter number of board to upload to: ")
		fmt.Scanf("%d", &index)

		if err != nil {
			logger.WithError(err).Fatal("Failed to get user input")
		}

		board := strings.Split(list[index], " ")
		device = board[0]
		fqbn = board[len(board)-1]
	}

	cli.SetArgs([]string{"compile", "--fqbn", fqbn, sketch})
	if err = cli.Execute(); err != nil {
		logger.WithError(err).Fatal("Failed to compile")
	}

	cli.SetArgs([]string{"upload", "-p", device, "--fqbn", fqbn, sketch})
	if err = cli.Execute(); err != nil {
		logger.WithError(err).Fatal("Failed to upload")
	}

	exec.Command("stty", "-F", device, baud).Run()

	watchLogsCmd := exec.Command("cat", device)
	watchLogsCmd.Stdout = os.Stdout
	watchLogsCmd.Stderr = os.Stderr

	watchLogsCmd.Run()
}
