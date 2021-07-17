package core

import (
	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/types"
	"github.com/robgonnella/ardi/v2/util"
	log "github.com/sirupsen/logrus"
)

// ArdiCore represents the core package of ardi
type ArdiCore struct {
	Cli       *cli.Wrapper
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
	Cli                *cli.Wrapper
	Logger             *log.Logger
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

// GetCompileOptsFromArgs return a list of compile opts from the given args
func (c *ArdiCore) GetCompileOptsFromArgs(fqbn string, buildProps []string, showProps bool, args []string) ([]*cli.CompileOpts, error) {
	ardiBuilds := c.Config.GetBuilds()
	_, defaultExists := ardiBuilds["default"]

	opts := []*cli.CompileOpts{}

	if len(args) == 0 {
		if len(ardiBuilds) == 1 {
			for k := range ardiBuilds {
				c.logger.Infof("Using build definition: %s", k)
				compileOpts, err := c.Config.GetCompileOpts(k)
				if err != nil {
					return nil, err
				}
				opts = append(opts, compileOpts)
			}
		} else if defaultExists {
			c.logger.Info("Using build definition: default")
			compileOpts, err := c.Config.GetCompileOpts("default")
			if err != nil {
				return nil, err
			}
			opts = append(opts, compileOpts)
		} else {
			c.logger.Info("Using ino file in current directory")
			project, err := util.ProcessSketch(".")
			if err != nil {
				return nil, err
			}
			compileOpts := &cli.CompileOpts{
				FQBN:       fqbn,
				BuildProps: buildProps,
				ShowProps:  showProps,
				SketchDir:  project.Directory,
				SketchPath: project.Sketch,
			}
			opts = append(opts, compileOpts)
		}

		return opts, nil
	}

	if len(args) == 1 {
		sketch := args[0]
		if _, ok := ardiBuilds[sketch]; ok {
			c.logger.Infof("Using build definition: %s", sketch)
			compileOpts, err := c.Config.GetCompileOpts(sketch)
			if err != nil {
				return nil, err
			}
			opts = append(opts, compileOpts)
			return opts, nil
		}

		c.logger.Info("Using ino file in current directory")
		project, err := util.ProcessSketch(sketch)
		if err != nil {
			return nil, err
		}

		compileOpts := &cli.CompileOpts{
			FQBN:       fqbn,
			BuildProps: buildProps,
			ShowProps:  showProps,
			SketchDir:  project.Directory,
			SketchPath: project.Sketch,
		}

		opts = append(opts, compileOpts)

		return opts, nil
	}

	for _, buildName := range args {
		c.logger.Infof("Using build definition: %s", buildName)
		compileOpts, err := c.Config.GetCompileOpts(buildName)
		if err != nil {
			return nil, err
		}
		opts = append(opts, compileOpts)
	}

	return opts, nil
}
