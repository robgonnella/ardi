package core

import (
	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/types"
	log "github.com/sirupsen/logrus"
)

// ArdiCore represents the core package of ardi
type ArdiCore struct {
	RPCClient rpc.Client
	Config    *ArdiJSON
	CliConfig *ArdiYAML
	Watch     *WatchCore
	Board     *BoardCore
	Compiler  *CompileCore
	Uploader  *UploadCore
	Lib       *LibCore
	Platform  *PlatformCore
}

// NewArdiCoreOpts options fore creating new ardi core
type NewArdiCoreOpts struct {
	Global             bool
	ArdiConfig         types.ArdiConfig
	ArduinoCliSettings types.ArduinoCliSettings
	Client             rpc.Client
	Logger             *log.Logger
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
		Config:    NewArdiJSON(ardiConf, opts.ArdiConfig, logger),
		CliConfig: NewArdiYAML(cliConf, opts.ArduinoCliSettings),
		Watch:     NewWatchCore(client, compiler, uploader, logger),
		Board:     NewBoardCore(client, logger),
		Compiler:  compiler,
		Uploader:  uploader,
		Lib:       NewLibCore(client, logger),
		Platform:  NewPlatformCore(client, logger),
	}
}
