package testutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	rpccommands "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
	log "github.com/sirupsen/logrus"
)

var here string

func init() {
	here, _ = filepath.Abs(".")
	log.SetOutput(ioutil.Discard)
}

// CleanCoreDir removes test data from core directory
func CleanCoreDir() {
	dataDir := path.Join(here, "../core/.ardi")
	jsonFile := path.Join(here, "../core/ardi.json")
	os.RemoveAll(dataDir)
	os.Remove(jsonFile)
}

// CleanCommandsDir removes project data from commands directory
func CleanCommandsDir() {
	projectDataDir := path.Join(here, "../commands/.ardi")
	projectJSONFile := path.Join(here, "../commands/ardi.json")
	os.RemoveAll(projectDataDir)
	os.Remove(projectJSONFile)
}

// CleanBuilds removes compiled test project builds
func CleanBuilds() {
	os.RemoveAll(path.Join(BlinkProjectDir(), "build"))
	os.RemoveAll(path.Join(BlinkCopyProjectDir(), "build"))
	os.RemoveAll(path.Join(Blink14400ProjectDir(), "build"))
	os.RemoveAll(path.Join(PixieProjectDir(), "build"))
}

// CleanAll removes all test data
func CleanAll() {
	CleanCoreDir()
	CleanCommandsDir()
	CleanBuilds()
}

// ArduinoMegaFQBN returns appropriate fqbn for arduino mega 2560
func ArduinoMegaFQBN() string {
	return "arduino:avr:mega"
}

// Esp8266Platform returns appropriate platform for esp8266
func Esp8266Platform() string {
	return "esp8266:esp8266"
}

// Esp8266WifiduinoFQBN returns appropriate fqbn for esp8266 board
func Esp8266WifiduinoFQBN() string {
	return "esp8266:esp8266:wifiduino"
}

// Esp8266BoardURL returns appropriate board url for esp8266 board
func Esp8266BoardURL() string {
	return "https://arduino.esp8266.com/stable/package_esp8266com_index.json"
}

// GenerateCmdBoard returns a single arduino-cli command Board
func GenerateCmdBoard(name, fqbn string) *rpccommands.Board {
	if fqbn == "" {
		fqbn = fmt.Sprintf("%s-fqbn", name)
	}
	return &rpccommands.Board{Name: name, Fqbn: fqbn}
}

// GenerateCmdBoards generate a list of arduino-cli command boards
func GenerateCmdBoards(n int) []*rpccommands.Board {
	var boards []*rpccommands.Board
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("test-board-%02d", i)
		b := GenerateCmdBoard(name, "")
		boards = append(boards, b)
	}
	return boards
}

// GenerateCmdPlatform generates a single named arduino-cli command platform
func GenerateCmdPlatform(name string, boards []*rpccommands.Board) *rpccommands.Platform {
	return &rpccommands.Platform{
		Id:     name,
		Boards: boards,
	}
}

// GenerateRPCBoard returns a single ardi rpc Board
func GenerateRPCBoard(name, fqbn string) *cli.BoardWithPort {
	if fqbn == "" {
		fqbn = fmt.Sprintf("%s-fqbn", name)
	}
	return &cli.BoardWithPort{
		Name: name,
		FQBN: fqbn,
		Port: "/dev/null",
	}
}

// GenerateRPCBoards generate a list of ardi rpc boards
func GenerateRPCBoards(n int) []*cli.BoardWithPort {
	var boards []*cli.BoardWithPort
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("test-board-%02d", i)
		b := GenerateRPCBoard(name, "")
		boards = append(boards, b)
	}
	return boards
}

// BlinkProjectDir returns path to blink project directory
func BlinkProjectDir() string {
	return path.Join(here, "../test_projects/blink")
}

// BlinkCopyProjectDir returns path to blink project directory
func BlinkCopyProjectDir() string {
	return path.Join(here, "../test_projects/blink2")
}

// Blink14400ProjectDir returns path to blink14400 project directory
func Blink14400ProjectDir() string {
	return path.Join(here, "../test_projects/blink14400")
}

// PixieProjectDir returns path to blink project directory
func PixieProjectDir() string {
	return path.Join(here, "../test_projects/pixie")
}
