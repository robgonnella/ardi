package core

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/robgonnella/ardi/v2/rpc"
	log "github.com/sirupsen/logrus"
)

// WatchCore represents core module for adi go commands
type WatchCore struct {
	logger     *log.Logger
	client     rpc.Client
	target     *Target
	project    *ProjectCore
	buildProps []string
	port       *SerialCore
	compiling  bool
	uploading  bool
}

// NewWatchCore returns new Project instance
func NewWatchCore(client rpc.Client, logger *log.Logger) *WatchCore {
	return &WatchCore{
		client: client,
		logger: logger,
	}
}

// Init intialize ardi-go core
func (w *WatchCore) Init(port, dir string, props []string) error {
	if w.project == nil {
		proj := NewProjectCore(w.client, w.logger)
		err := proj.SetConfigHelpers()
		if err != nil {
			return err
		}
		if err := proj.ProcessSketch(dir); err != nil {
			return err
		}
		w.project = proj
	}

	connectedBoards := w.client.ConnectedBoards()
	allBoards := w.client.AllBoards()

	if w.target == nil {
		target, err := NewTarget(connectedBoards, allBoards, "", true, w.logger)
		if err != nil {
			return err
		}
		w.target = target
	}

	if w.port == nil {
		w.port = NewSerialCore(w.target.Board.Port, w.project.Baud, w.logger)
	}

	return nil
}

// Upload compiled sketches to the specified board
func (w *WatchCore) Upload() error {
	w.port.Stop()
	w.waitForPreviousCompile()
	w.waitForPreviousUpload()
	w.logger.Info("Uploading...")

	fqbn := w.target.Board.FQBN
	device := w.target.Board.Port
	sketchDir := w.project.Directory

	w.uploading = true
	if err := w.client.Upload(fqbn, sketchDir, device); err != nil {
		w.logger.WithError(err).Error("Failed to upload sketch")
		w.uploading = false
		return err
	}

	w.uploading = false
	return nil
}

// Compile the specified sketch
func (w *WatchCore) Compile() error {
	w.port.Stop()
	w.waitForPreviousCompile()
	w.waitForPreviousUpload()

	fqbn := w.target.Board.FQBN
	sketchDir := w.project.Directory
	sketch := w.project.Sketch
	buildProps := w.buildProps

	w.compiling = true
	opts := rpc.CompileOpts{
		FQBN:       fqbn,
		SketchDir:  sketchDir,
		SketchPath: sketch,
		ExportName: "",
		BuildProps: buildProps,
		ShowProps:  false,
	}
	if err := w.client.Compile(opts); err != nil {
		w.logger.WithError(err).Error("Failed to compile sketch")
		w.compiling = false
		return err
	}

	w.compiling = false
	return nil
}

// WatchSketch responds to changes in a given sketch file by automatically
// recompiling and re-uploading.
func (w *WatchCore) WatchSketch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		w.logger.WithError(err).Error("Failed to watch directory for changes")
	}
	defer watcher.Close()

	err = watcher.Add(w.project.Sketch)
	if err != nil {
		w.logger.WithError(err).Error("Failed to watch directory for changes")
	}

	go w.WatchLogs()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				break
			}
			w.logger.Debugf("event: %+v", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				w.logger.Debugf("modified file: %s", event.Name)
				err := w.Compile()
				if err != nil {
					w.Upload()
					go w.WatchLogs()
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			w.logger.WithError(err).Warn("Watch error")
		}
	}
}

// WatchLogs connects to a serial port at a specified baud rate and prints
// any logs received.
func (w *WatchCore) WatchLogs() {
	baud := w.project.Baud
	device := w.target.Board.Port

	logFields := log.Fields{"baud": baud, "device": device}

	w.port.Stop()
	w.waitForPreviousCompile()
	w.waitForPreviousUpload()

	w.logger.WithFields(logFields).Info("Watching logs...")
	w.port.Watch()
}

// private helpers
func (w *WatchCore) waitForPreviousUpload() {
	// block until target is no longer uploading
	for {
		if !w.uploading {
			break
		}
		w.logger.Info("Waiting for previous upload to finish...")
		time.Sleep(time.Second)
	}
}

func (w *WatchCore) waitForPreviousCompile() {
	// block until target is no longer compiling
	for {
		if !w.compiling {
			break
		}
		w.logger.Info("Waiting for previous compile to finish...")
		time.Sleep(time.Second)
	}
}
