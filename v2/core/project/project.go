package project

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

// Project represents an arduino project
type Project struct {
	Sketch    string
	Directory string
	Baud      int
}

// New returns new Project instance
func New(sketchDir string, logger *log.Logger) (*Project, error) {
	if sketchDir == "" {
		msg := "Must provide a sketch directory as an argument"
		err := errors.New("Missing directory argument")
		logger.WithError(err).Error(msg)
		return nil, err
	}

	// Guard in case someone tries to pass full path to .ino file
	sketchDir = path.Dir(sketchDir)

	sketchFile, err := findSketch(sketchDir, logger)
	if err != nil {
		return nil, err
	}

	sketchBaud := parseSketchBaud(sketchFile, logger)
	if sketchBaud != 0 {
		fmt.Println("")
		logger.WithField("detected baud", sketchBaud).Info("Detected baud rate from sketch file.")
		fmt.Println("")
	}

	return &Project{
		Sketch:    sketchFile,
		Directory: sketchDir,
		Baud:      sketchBaud,
	}, nil
}

// helpers
func findSketch(directory string, logger *log.Logger) (string, error) {
	sketchFile := ""

	d, err := os.Open(directory)
	if err != nil {
		logger.WithError(err).Error("Failed to open sketch directory")
		return "", err
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		logger.WithError(err).Error("Cannot process .ino file")
		return "", err
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ".ino" {
				sketchFile = path.Join(directory, file.Name())
			}
		}
	}
	if sketchFile == "" {
		msg := fmt.Sprintf("Failed to find .ino file in %s", directory)
		logger.Error(msg)
		return "", errors.New(msg)
	}

	if sketchFile, err = filepath.Abs(sketchFile); err != nil {
		msg := "Could not resolve sketch file path"
		logger.WithError(err).Error(msg)
		return "", errors.New(msg)
	}

	return sketchFile, nil
}

func parseSketchBaud(sketch string, logger *log.Logger) int {
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
