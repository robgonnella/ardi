package core

import (
	"errors"
	"fmt"
	"sort"
	"text/tabwriter"

	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/types"
	"github.com/robgonnella/ardi/v2/util"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// ArdiCore represents the core package of ardi
type ArdiCore struct {
	RPCClient rpc.Client
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
	Global             bool
	ArdiConfig         types.ArdiConfig
	ArduinoCliSettings types.ArduinoCliSettings
	Client             rpc.Client
	Logger             *log.Logger
}

// CompileSketchOpts options for compiling sketch directory or file
type CompileSketchOpts struct {
	Sketch              string
	FQBN                string
	BuildPros           []string
	ShowProps           bool
	OnlyConnectedBoards bool
}

// CompileArdiBuildOpts options for compiling a build specified in ardi.json
type CompileArdiBuildOpts struct {
	BuildName           string
	OnlyConnectedBoards bool
}

// NewArdiCore returns a new ardi core
func NewArdiCore(opts NewArdiCoreOpts) *ArdiCore {
	ardiConf := paths.ArdiProjectConfig
	cliConf := paths.ArduinoCliProjectConfig

	if opts.Global {
		ardiConf = paths.ArdiGlobalConfig
		cliConf = paths.ArduinoCliGlobalConfig
	}

	client := opts.Client
	logger := opts.Logger

	compiler := NewCompileCore(client, logger)
	uploader := NewUploadCore(client, logger)

	return &ArdiCore{
		RPCClient: client,
		Config:    NewArdiConfig(ardiConf, opts.ArdiConfig, logger),
		CliConfig: NewArdiYAML(cliConf, opts.ArduinoCliSettings),
		Watcher:   NewWatchCore(compiler, uploader, logger),
		Board:     NewBoardCore(client, logger),
		Compiler:  compiler,
		Uploader:  uploader,
		Lib:       NewLibCore(client, logger),
		Platform:  NewPlatformCore(client, logger),
		logger:    logger,
	}
}

// CompileArdiBuild compiles specified build from ardi.json
func (c *ArdiCore) CompileArdiBuild(buildOpts CompileArdiBuildOpts) (*rpc.CompileOpts, *rpc.Board, error) {
	compileOpts, err := c.Config.GetCompileOpts(buildOpts.BuildName)
	if err != nil {
		return nil, nil, err
	}
	board, err := c.GetTargetBoard(compileOpts.FQBN, buildOpts.OnlyConnectedBoards)
	if err != nil {
		return nil, nil, err
	}
	fields := logrus.Fields{
		"sketch": compileOpts.SketchPath,
		"fqbn":   compileOpts.FQBN,
	}
	c.logger.WithFields(fields).Info("Compiling...")
	if err := c.Compiler.Compile(*compileOpts); err != nil {
		return nil, nil, err
	}

	return compileOpts, board, nil
}

// CompileSketch compiles specified sketch directory or sketch file
func (c *ArdiCore) CompileSketch(sketchOpts CompileSketchOpts) (*rpc.CompileOpts, *rpc.Board, error) {

	board, err := c.GetTargetBoard(sketchOpts.FQBN, sketchOpts.OnlyConnectedBoards)
	if err != nil {
		return nil, nil, err
	}

	project, err := util.ProcessSketch(sketchOpts.Sketch)
	if err != nil {
		return nil, nil, err
	}

	compileOpts := rpc.CompileOpts{
		FQBN:       board.FQBN,
		SketchDir:  project.Directory,
		SketchPath: project.Sketch,
		ExportName: "",
		BuildProps: sketchOpts.BuildPros,
		ShowProps:  sketchOpts.ShowProps,
	}
	fields := logrus.Fields{
		"sketch": compileOpts.SketchPath,
		"fqbn":   compileOpts.FQBN,
	}
	c.logger.WithFields(fields).Info("Compiling...")
	if err := c.Compiler.Compile(compileOpts); err != nil {
		return nil, nil, err
	}

	return &compileOpts, board, nil
}

// GetTargetBoard returns target info for a connected & disconnected boards
func (c *ArdiCore) GetTargetBoard(fqbn string, onlyConnected bool) (*rpc.Board, error) {
	fqbnErr := errors.New("you must specify a board fqbn to compile - you can find a list of board fqbns for installed platforms above")
	connectedBoardsErr := errors.New("No connected boards detected")
	connectedBoards := c.RPCClient.ConnectedBoards()
	allBoards := c.RPCClient.AllBoards()

	if fqbn != "" {
		if onlyConnected {
			for _, b := range connectedBoards {
				if b.FQBN == fqbn {
					return b, nil
				}
			}
			return nil, connectedBoardsErr
		}
		return &rpc.Board{FQBN: fqbn}, nil
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

	return nil, errors.New("Error parsing target")
}

// private helpers
func (c *ArdiCore) printFQBNs(boardList []*rpc.Board, logger *log.Logger) {
	sort.Slice(boardList, func(i, j int) bool {
		return boardList[i].Name < boardList[j].Name
	})

	c.printBoardsWithIndices(boardList, logger)
}

func (c *ArdiCore) printBoardsWithIndices(boards []*rpc.Board, logger *log.Logger) {
	w := tabwriter.NewWriter(logger.Out, 0, 0, 8, ' ', 0)
	defer w.Flush()
	w.Write([]byte("No.\tName\tFQBN\n"))
	for i, b := range boards {
		w.Write([]byte(fmt.Sprintf("%d\t%s\t%s\n", i, b.Name, b.FQBN)))
	}
}
