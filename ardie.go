package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	_ "github.com/arduino/arduino-cli/cli"
	log "github.com/sirupsen/logrus"
)

const cliCmd = "arduino-cli"
const defaultBaud = "9600"

type boardInfo struct {
	FQBN   string
	Device string
}

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
	var boards []byte
	var fqbn string
	var device string
	var sketch string
	var err error
	baud := "9600"

	if len(os.Args) == 1 {
		log.WithError(errors.New("Missing sketch arguemnet")).Fatal("Must provide a sketch name as an argument to upload")
	}

	sketch = os.Args[1]
	sketch = strings.Replace(sketch, "sketches/", "", 1)
	sketch = fmt.Sprintf("sketches/%s", sketch)

	if len(os.Args) == 3 {
		baud = os.Args[2]
	}

	updateIndexCmd := exec.Command(cliCmd, "core", "update-index")
	updateIndexCmd.Stdout = os.Stdout
	updateIndexCmd.Stderr = os.Stderr

	installCoreCmd := exec.Command(cliCmd, "core", "install", "arduino:avr")
	installCoreCmd.Stdout = os.Stdout
	installCoreCmd.Stderr = os.Stderr

	listCmd := exec.Command(cliCmd, "board", "list")

	if err = updateIndexCmd.Run(); err != nil {
		log.WithError(err).Fatal("Failed to update index")
	}

	if err = installCoreCmd.Run(); err != nil {
		log.WithError(err).Fatal("Failed to install avr core")
	}

	if boards, err = listCmd.Output(); err != nil {
		log.WithError(err).Fatal("Failed to get list command output")
	}

	printableList := strings.SplitAfterN(string(boards), "\n", -1)
	printableList = Filter(printableList, func(s string) bool {
		return !strings.Contains(s, "Unknown") && s != ""
	})

	list := strings.Split(string(boards), "\n")
	list = Filter(list, func(s string) bool {
		return !strings.Contains(s, "Unknown") && !strings.Contains(s, "Board Name") && s != ""
	})

	if len(list) == 0 {
		log.WithError(errors.New("No boards detected")).Fatal("Cannot upload sketch")
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
			log.WithError(err).Fatal("Failed to get user input")
		}

		board := strings.Split(list[index], " ")
		device = board[0]
		fqbn = board[len(board)-1]
	}

	compileCmd := exec.Command(cliCmd, "compile", "--fqbn", fqbn, sketch)
	compileCmd.Stdout = os.Stdout
	compileCmd.Stderr = os.Stderr

	if err = compileCmd.Run(); err != nil {
		log.WithError(err).Fatal("Failed to compile")
	}

	uploadCmd := exec.Command(cliCmd, "upload", "-p", device, "--fqbn", fqbn, sketch)
	uploadCmd.Stdout = os.Stdout
	uploadCmd.Stderr = os.Stderr

	if err = uploadCmd.Run(); err != nil {
		log.WithError(err).Fatal("Failed to upload")
	}

	exec.Command("stty", "-F", device, baud).Run()

	watchLogsCmd := exec.Command("cat", device)
	watchLogsCmd.Stdout = os.Stdout
	watchLogsCmd.Stderr = os.Stderr

	watchLogsCmd.Run()
}
