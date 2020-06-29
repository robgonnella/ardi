package core

import (
	"github.com/robgonnella/ardi/v2/rpc"
	log "github.com/sirupsen/logrus"
)

// ArdiCore represents the core package of ardi
type ArdiCore struct {
	Watch    *WatchCore
	Board    *BoardCore
	Compiler *CompileCore
	Lib      *LibCore
	Platform *PlatformCore
	Project  *ProjectCore
}

// NewArdiCore returns a new ardi core
func NewArdiCore(client rpc.Client, logger *log.Logger) *ArdiCore {
	return &ArdiCore{
		Watch:    NewWatchCore(client, logger),
		Board:    NewBoardCore(client, logger),
		Compiler: NewCompileCore(client, logger),
		Lib:      NewLibCore(client, logger),
		Platform: NewPlatformCore(client, logger),
		Project:  NewProjectCore(client, logger),
	}
}
