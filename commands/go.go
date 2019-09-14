package commands

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/robgonnella/ardi/ardi"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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

func getSketch() (string, string) {
	if len(os.Args) < 3 {
		return "", ""
	}

	sketchDir := os.Args[2]

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

func processSketch(baud int) (string, string, int) {
	sketchDir, sketchFile := getSketch()

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

func process(baud int, watchSketch, verbose bool) {
	if verbose {
		logger.SetLevel(log.DebugLevel)
		ardi.SetLogLevel(log.DebugLevel)
	} else {
		logger.SetLevel(log.InfoLevel)
		ardi.SetLogLevel(log.InfoLevel)
	}

	sketchDir, sketchFile, baud := processSketch(baud)

	logFields := log.Fields{"baud": baud, "sketch": sketchDir}
	logWithFields := logger.WithFields(logFields)

	conn, client, rpcInstance := ardi.Initialize()
	defer conn.Close()

	list := ardi.GetTargetList(client, rpcInstance, sketchDir, sketchFile, baud)

	logWithFields.Debug("Parsing target")
	target := ardi.GetTargetInfo(list)

	logWithFields.WithField("target", target).Debug("Found target")

	ardi.Compile(client, rpcInstance, &target)
	if target.CompileError {
		os.Exit(1)
	}
	ardi.Upload(client, rpcInstance, &target)

	if watchSketch {
		ardi.WatchSketch(client, rpcInstance, &target)
	} else {
		ardi.WatchLogs(&target)
	}

}

func getGoCommand() *cobra.Command {

	var baud int
	var watchSketch bool
	var verbose bool

	var goCmd = &cobra.Command{
		Use:   "go [sketch]",
		Short: "Compile and upload code to an arduino board",
		Long: "Compile and upload code to an arduino board. Simply pass the\n" +
			"directory containing the .ino file as the first argument. To watch\n" +
			"your sketch file for changes and auto re-compile & re-upload, use\n" +
			"the --watch flag. You can also specify the baud rate with --baud\n" +
			"(default is 9600).",
		Run: func(cmd *cobra.Command, args []string) {
			process(baud, watchSketch, verbose)
		},
	}
	goCmd.Flags().IntVarP(&baud, "baud", "b", 9600, "specify sketch baud rate")
	goCmd.Flags().BoolVarP(&watchSketch, "watch", "w", false, "watch for changes, recompile and reupload")
	goCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "print all logs")

	return goCmd
}
