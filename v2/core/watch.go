package core

import (
	"errors"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
	log "github.com/sirupsen/logrus"
)

// WatchCore represents core module for adi go commands
type WatchCore struct {
	logger           *log.Logger
	uploader         *UploadCore
	compiler         *CompileCore
	port             SerialPort
	watcher          *FileWatcher
	board            *cli.BoardWithPort
	compileOpts      *cli.CompileOpts
	baud             int
	processingUpdate bool
}

// WatchCoreTargets targets for watching, recompiling, and reuploading
type WatchCoreTargets struct {
	Board       *cli.BoardWithPort
	CompileOpts *cli.CompileOpts
	Baud        int
	Port        SerialPort
}

// NewWatchCore returns new Project instance
func NewWatchCore(compiler *CompileCore, uploader *UploadCore, logger *log.Logger) *WatchCore {
	return &WatchCore{
		uploader:         uploader,
		compiler:         compiler,
		logger:           logger,
		processingUpdate: false,
	}
}

// SetTargets sets the board and compile options for the watcher
func (w *WatchCore) SetTargets(targets WatchCoreTargets) error {
	if w.watcher != nil {
		w.watcher.Stop()
		w.watcher = nil
	}

	board := targets.Board
	compileOpts := targets.CompileOpts
	baud := targets.Baud

	watcher, err := NewFileWatcher(compileOpts.SketchPath, w.logger)
	if err != nil {
		return err
	}
	w.watcher = watcher

	if targets.Port != nil {
		w.port = targets.Port
	} else {
		w.port = NewArdiSerialPort(board.Port, baud, w.logger)
	}

	w.board = board
	w.compileOpts = compileOpts
	w.baud = baud
	w.watcher.AddListener(w.onFileChange)
	return nil
}

// Watch responds to changes in a given sketch file by automatically
// recompiling and re-uploading.
func (w *WatchCore) Watch() error {
	if !w.hasTargets() {
		return errors.New("must call SetTargets before calling watch")
	}

	go w.port.Watch()
	return w.watcher.Watch()
}

// Stop deletes watcher and unattaches from port
func (w *WatchCore) Stop() {
	if w.port != nil {
		w.port.Close()
		w.port = nil
	}

	if w.watcher != nil {
		w.watcher.Close()
		w.watcher = nil
	}

	w.baud = 0
	w.board = nil
	w.compileOpts = nil
}

// private
func (w *WatchCore) onFileChange() {
	if !w.hasTargets() {
		err := errors.New("watch targets have gone missing")
		w.logger.WithError(err).Error()
		return
	}

	if w.processingUpdate {
		w.logger.Debug("already processing file change...")
		return
	}

	w.processingUpdate = true
	w.watcher.Stop()
	w.port.Close()

	defer func() {
		w.processingUpdate = false
		if w.watcher != nil {
			w.watcher.Restart()
		}
	}()

	err := w.compiler.Compile(*w.compileOpts)
	if err != nil {
		return
	}
	err = w.uploader.Upload(w.board, w.compileOpts.SketchDir)
	if err != nil {
		return
	}

	if w.port != nil {
		go w.port.Watch()
	}
}

func (w *WatchCore) hasTargets() bool {
	return w.port != nil && w.board != nil && w.compileOpts != nil && w.baud != 0 && w.watcher != nil
}
