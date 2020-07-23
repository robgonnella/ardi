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
	port     SerialPort
	watcher  *FileWatcher
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
func (w *WatchCore) Watch(compileOpts rpc.CompileOpts, target Target, baud int, port SerialPort) error {
	if w.watcher != nil {
		w.watcher.Close()
		w.watcher = nil
	}

	watcher, err := NewFileWatcher(compileOpts.SketchPath, w.logger)
	if err != nil {
		return err
	}
	w.watcher = watcher

	if w.port != nil {
		w.port.Stop()
		w.port = nil
	}

	if port == nil {
		port = NewArdiSerialPort(target.Board.Port, baud, w.logger)
	} else {
		port.Stop()
	}
	w.port = port

	watcher.AddListener(func() {
		port.Stop()
		err := w.compiler.Compile(compileOpts)
		if err == nil {
			w.logger.Info("Reuploading")
			w.uploader.Upload(target, compileOpts.SketchDir)
			w.logger.Info("Upload successful")
			go port.Watch()
		}
	})

	go port.Watch()
	watcher.Watch()

	return nil
}

// Stop deletes watcher and unattaches from port
func (w *WatchCore) Stop() {
	if w.watcher != nil {
		w.watcher.Close()
		w.watcher = nil
	}
	if w.port != nil {
		w.port.Stop()
		w.port = nil
	}
}
