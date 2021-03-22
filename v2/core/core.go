package core

import (
	"errors"
	"fmt"
	"sort"
	"text/tabwriter"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/types"
	"github.com/robgonnella/ardi/v2/util"
	log "github.com/sirupsen/logrus"
)

// ArdiCore represents the core package of ardi
type ArdiCore struct {
	Cli       cli.Cli
	Config    *ArdiConfig
	CliConfig *ArdiYAML
	Watcher   *WatchCore
	Board     *BoardCore
	Compiler  *CompileCore
	Uploader  *UploadCore
	Lib       *LibCore
	Platform  *PlatformCore
	logger    *log.Logger
}

// NewArdiCoreOpts options fore creating new ardi core
type NewArdiCoreOpts struct {
	ArdiConfig         types.ArdiConfig
	ArduinoCliSettings types.ArduinoCliSettings
	Cli                cli.Cli
	Logger             *log.Logger
}

// CompileSketchOpts options for compiling sketch directory or file
type CompileSketchOpts struct {
	Sketch    string
	FQBN      string
	BuildPros []string
	ShowProps bool
}

// NewArdiCore returns a new ardi core
func NewArdiCore(opts NewArdiCoreOpts) *ArdiCore {
	ardiConf := paths.ArdiProjectConfig
	cliConf := paths.ArduinoCliProjectConfig

	cli := opts.Cli
	logger := opts.Logger

	compiler := NewCompileCore(cli, logger)
	uploader := NewUploadCore(cli, logger)

	return &ArdiCore{
		Cli:       cli,
		Config:    NewArdiConfig(ardiConf, opts.ArdiConfig, logger),
		CliConfig: NewArdiYAML(cliConf, opts.ArduinoCliSettings),
		Watcher:   NewWatchCore(compiler, uploader, logger),
		Board:     NewBoardCore(cli, logger),
		Compiler:  compiler,
		Uploader:  uploader,
		Lib:       NewLibCore(cli, logger),
		Platform:  NewPlatformCore(cli, logger),
		logger:    logger,
	}
}

// CompileArdiBuild compiles specified build from ardi.json
func (c *ArdiCore) CompileArdiBuild(buildName string) (*cli.CompileOpts, error) {
	compileOpts, err := c.Config.GetCompileOpts(buildName)
	if err != nil {
		return nil, err
	}

	fields := log.Fields{
		"sketch": compileOpts.SketchPath,
		"fqbn":   compileOpts.FQBN,
	}
	c.logger.WithFields(fields).Info("Compiling...")

	if err := c.Compiler.Compile(*compileOpts); err != nil {
		return nil, err
	}

	return compileOpts, nil
}

// CompileSketch compiles specified sketch directory or sketch file
func (c *ArdiCore) CompileSketch(sketchOpts CompileSketchOpts) (*cli.CompileOpts, error) {
	project, err := util.ProcessSketch(sketchOpts.Sketch)
	if err != nil {
		return nil, err
	}

	compileOpts := cli.CompileOpts{
		FQBN:       sketchOpts.FQBN,
		SketchDir:  project.Directory,
		SketchPath: project.Sketch,
		BuildProps: sketchOpts.BuildPros,
		ShowProps:  sketchOpts.ShowProps,
	}

	fields := log.Fields{
		"sketch": compileOpts.SketchPath,
		"fqbn":   compileOpts.FQBN,
	}
	c.logger.WithFields(fields).Info("Compiling...")

	if err := c.Compiler.Compile(compileOpts); err != nil {
		return nil, err
	}

	return &compileOpts, nil
}

// GetTargetBoard returns target info for a connected & disconnected boards
func (c *ArdiCore) GetTargetBoard(fqbn, port string, onlyConnected bool) (*cli.Board, error) {
	if fqbn != "" && port != "" {
		return &cli.Board{
			FQBN: fqbn,
			Port: port,
		}, nil
	}

	fqbnErr := errors.New("you must specify a board fqbn to compile - you can find a list of board fqbns for installed platforms above")
	connectedBoardsErr := errors.New("no connected boards detected")
	connectedBoards := c.Cli.ConnectedBoards()
	allBoards := c.Cli.AllBoards()

	if fqbn != "" {
		if onlyConnected {
			for _, b := range connectedBoards {
				if b.FQBN == fqbn {
					return b, nil
				}
			}
			return nil, connectedBoardsErr
		}
		return &cli.Board{FQBN: fqbn}, nil
	}

	if len(connectedBoards) == 0 {
		if onlyConnected {
			return nil, connectedBoardsErr
		}
		c.printFQBNs(allBoards, c.logger)
		return nil, fqbnErr
	}

	if len(connectedBoards) == 1 {
		return connectedBoards[0], nil
	}

	if len(connectedBoards) > 1 {
		c.printFQBNs(connectedBoards, c.logger)
		return nil, fqbnErr
	}

	return nil, errors.New("error parsing target")
}

// private helpers
func (c *ArdiCore) printFQBNs(boardList []*cli.Board, logger *log.Logger) {
	sort.Slice(boardList, func(i, j int) bool {
		return boardList[i].Name < boardList[j].Name
	})

	c.printBoardsWithIndices(boardList, logger)
}

func (c *ArdiCore) printBoardsWithIndices(boards []*cli.Board, logger *log.Logger) {
	w := tabwriter.NewWriter(logger.Out, 0, 0, 8, ' ', 0)
	defer w.Flush()
	w.Write([]byte("No.\tName\tFQBN\n"))
	for i, b := range boards {
		w.Write([]byte(fmt.Sprintf("%d\t%s\t%s\n", i, b.Name, b.FQBN)))
	}
}
