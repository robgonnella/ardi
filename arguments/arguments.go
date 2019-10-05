package arguments

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

var logger = log.New()

func parseBaudRate(sketch string) int {
	var baud int
	rgx := regexp.MustCompile(`Serial\.begin\((\d+)\);`)
	file, err := os.Open(sketch)
	if err != nil {
		// Log the error and return 0 for baud to let script continue
		// with either default value or value specified from command-line.
		logger.WithError(err).
			WithField("sketch", sketch).
			Info("Failed to read sketch")
		return baud
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := scanner.Text()
		if match := rgx.MatchString(text); match {
			stringBaud := strings.TrimSpace(rgx.ReplaceAllString(text, "$1"))
			if baud, err = strconv.Atoi(stringBaud); err != nil {
				// set baud to 0 and let script continue with either default
				// value or value specified from command-line.
				logger.WithError(err).Info("Failed to parse baud rate from sketch")
				baud = 0
			}
			break
		}
	}

	return baud
}

func findInoFile(sketchDir string) string {
	sketchFile := ""

	d, err := os.Open(sketchDir)
	if err != nil {
		logger.WithError(err).Fatal("Failed to open sketch directory")
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ".ino" {
				sketchFile = path.Join(sketchDir, file.Name())
			}
		}
	}
	if sketchFile == "" {
		logger.Fatalf("Failed to find .ino file in %s", sketchDir)
	}

	if sketchFile, err = filepath.Abs(sketchFile); err != nil {
		logger.WithError(err).Fatal("Could not resolve sketch file path")
	}

	return sketchFile
}

// GetSketchParts sketch directory path and sketch file path.
func GetSketchParts(sketchDir string) (string, string) {
	sketchFile := findInoFile(sketchDir)
	sketchDir = path.Dir(sketchFile)
	return sketchDir, sketchFile
}

// ProcessSketch reads arguments, finds sketch file, determines appropriate baud
// rate and returns sketch directory path, sketch file path, and baud rate.
func ProcessSketch(sketchDir string, baud int) (string, string, int) {
	sketchDir, sketchFile := GetSketchParts(sketchDir)

	if sketchDir == "" {
		logger.WithError(errors.New("Missing sketch argument")).Fatal("Must provide a sketch name as an argument to upload")
	}
	parsedBaud := parseBaudRate(sketchFile)

	if parsedBaud != 0 && parsedBaud != baud {
		fmt.Println("")
		logger.Info("Detected a different baud rate from sketch file.")
		logger.WithField("detected baud", parsedBaud).Info("Using detected baud rate")
		fmt.Println("")
		baud = parsedBaud
	}

	return sketchDir, sketchFile, baud
}
