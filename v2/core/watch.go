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
	fileWatcher      *FileWatcher
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
}

// WatchCoreOption represents options for WatchCore
type WatchCoreOption = func(c *WatchCore)

// NewWatchCore returns new Project instance
func NewWatchCore(logger *log.Logger, options ...WatchCoreOption) *WatchCore {
	c := &WatchCore{
		logger:           logger,
		processingUpdate: false,
	}

	for _, o := range options {
		o(c)
	}

	return c
}

// WithWatchCoreUploader allows an injectable UploadCore
func WithWatchCoreUploader(uploader *UploadCore) WatchCoreOption {
	return func(c *WatchCore) {
		c.uploader = uploader
	}
}

// WithWatchCoreCompiler allows an injectable CompileCore
func WithWatchCoreCompiler(compiler *CompileCore) WatchCoreOption {
	return func(c *WatchCore) {
		c.compiler = compiler
	}
}

// SetTargets sets the board and compile options for the watcher
func (w *WatchCore) SetTargets(targets WatchCoreTargets) error {
	if w.fileWatcher != nil {
		w.fileWatcher.Stop()
		w.fileWatcher = nil
	}

	board := targets.Board
	compileOpts := targets.CompileOpts
	baud := targets.Baud

	watcher, err := NewFileWatcher(compileOpts.SketchPath, w.logger)
	if err != nil {
		return err
	}
	w.fileWatcher = watcher

	w.uploader.SetPortTargets(board.Port, baud)

	w.board = board
	w.compileOpts = compileOpts
	w.baud = baud
	w.fileWatcher.AddListener(w.onFileChange)
	return nil
}

// Watch responds to changes in a given sketch file by automatically
// recompiling and re-uploading.
func (w *WatchCore) Watch() error {
	if !w.hasTargets() {
		return errors.New("must call SetTargets before calling watch")
	}

	go w.uploader.Attach()
	return w.fileWatcher.Watch()
}

// Stop deletes watcher and unattaches from port
func (w *WatchCore) Stop() {
	w.uploader.Detach()

	if w.fileWatcher != nil {
		w.fileWatcher.Close()
		w.fileWatcher = nil
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
	w.fileWatcher.Stop()
	w.uploader.Detach()

	defer func() {
		w.processingUpdate = false
		if w.fileWatcher != nil {
			w.fileWatcher.Restart()
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

	w.uploader.SetPortTargets(w.board.Port, w.baud)
	go w.uploader.Attach()
}

func (w *WatchCore) hasTargets() bool {
	return w.board != nil && w.compileOpts != nil && w.baud != 0 && w.fileWatcher != nil
}
