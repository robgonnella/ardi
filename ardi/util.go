package ardi

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func createDataDir() {
	logger.Debug("Creating data directory")
	_ = os.MkdirAll(DataDir, 0777)
	libDirConfig := LibraryDirConfig{
		ProxyType:      "auto",
		SketchbookPath: ".",
		ArduinoData:    ".",
		BoardManager:   make(map[string]interface{}),
	}
	yamlConfig, _ := yaml.Marshal(&libDirConfig)
	ioutil.WriteFile(GlobalLibConfig, yamlConfig, 0644)
}

func printConnectedBoardsWithIndices(list []TargetInfo) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].BoardName < list[j].BoardName
	})
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', 0)
	defer w.Flush()
	fmt.Fprintln(w, "No.\tBoard\tDevice")
	for i, board := range list {
		fmt.Fprintf(w, "%d\t%s\t%s\n", i, board.FQBN, board.Device)
	}
}

func printSupportedBoardsWithIndices(boards []boardInfo) {
	sort.Slice(boards, func(i, j int) bool {
		return boards[i].Name < boards[j].Name
	})
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', 0)
	defer w.Flush()
	fmt.Fprintln(w, "No.\tName\tFQBN")
	for i, b := range boards {
		fmt.Fprintf(w, "%d\t%s\t%s\n", i, b.Name, b.FQBN)
	}
}

func stopLogs(target *TargetInfo) {
	if target.Stream != nil {
		logWithField := logger.WithField("device", target.Device)
		logWithField.Info("Closing serial port connection")
		if err := target.Stream.Close(); err != nil {
			logWithField.WithError(err).Fatal("Failed to close serial port connection")
		}
		if err := target.Stream.Flush(); err != nil {
			logWithField.WithError(err).Fatal("Failed to flush serial port connection")
		}
		target.Stream = nil
		// block until all logs have stopped
		for {
			if !target.Logging {
				break
			}
		}
	}
}

func waitForPreviousUpload(target *TargetInfo) {
	// block until target is no longer uploading
	for {
		if !target.Uploading {
			break
		}
		logger.Info("Waiting for previous upload to finish...")
		time.Sleep(time.Second)
	}
}

func waitForPreviousCompile(target *TargetInfo) {
	// block until target is no longer compiling
	for {
		if !target.Compiling {
			break
		}
		logger.Info("Waiting for previous compile to finish...")
		time.Sleep(time.Second)
	}
}

func isVerbose() bool {
	return logger.Level == log.DebugLevel
}
