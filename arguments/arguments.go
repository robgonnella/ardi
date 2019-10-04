package arguments

import (
	"bufio"
	"errors"
	"fmt"
	"os"
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

// GetSketchParts sketch directory path and sketch file path.
func GetSketchParts(sketchDir string) (string, string) {
	if !strings.Contains(sketchDir, "/") {
		sketchDir = fmt.Sprintf("sketches/%s", sketchDir)
	}

	if strings.HasSuffix(sketchDir, "/") {
		sketchDir = strings.TrimSuffix(sketchDir, "/")
	}

	sketchParts := strings.Split(sketchDir, "/")
	sketchName := sketchParts[len(sketchParts)-1]
	sketchFile := fmt.Sprintf("%s/%s.ino", sketchDir, sketchName)
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
