package core

import (
	"github.com/robgonnella/ardi/v2/rpc"
	log "github.com/sirupsen/logrus"
)

// WatchCore represents core module for adi go commands
type WatchCore struct {
	logger   *log.Logger
	uploader *UploadCore
	compiler *CompileCore
}

// NewWatchCore returns new Project instance
func NewWatchCore(compiler *CompileCore, uploader *UploadCore, logger *log.Logger) *WatchCore {
	return &WatchCore{
		uploader: uploader,
		compiler: compiler,
		logger:   logger,
	}
}

// Watch responds to changes in a given sketch file by automatically
// recompiling and re-uploading.
func (w *WatchCore) Watch(compileOpts rpc.CompileOpts, target Target, baud int) error {
	watcher, err := NewFileWatcher(compileOpts.SketchPath, w.logger)
	if err != nil {
		return err
	}

	port := NewArdiSerialPort(target.Board.Port, baud, w.logger)

	watcher.AddListener(func() {
		port.Stop()
		err := w.compiler.Compile(compileOpts)
		if err != nil {
			w.uploader.Upload(target, compileOpts.SketchDir)
			go port.Watch()
		}
	})

	go port.Watch()
	watcher.Watch()

	return nil
}
