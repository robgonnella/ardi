package core

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/types"
	"github.com/robgonnella/ardi/v2/util"
	log "github.com/sirupsen/logrus"
)

// WatchCore represents core module for adi go commands
type WatchCore struct {
	logger     *log.Logger
	uploader   *UploadCore
	compiler   *CompileCore
	client     rpc.Client
	target     *Target
	project    *types.Project
	buildProps []string
	port       SerialPort
}

// NewWatchCore returns new Project instance
func NewWatchCore(client rpc.Client, compiler *CompileCore, uploader *UploadCore, logger *log.Logger) *WatchCore {
	return &WatchCore{
		client:   client,
		uploader: uploader,
		compiler: compiler,
		logger:   logger,
	}
}

// Init intialize ardi-go core
func (w *WatchCore) Init(port, dir string, props []string) error {
	if w.project == nil {
		project, err := util.ProcessSketch(dir)
		if err != nil {
			return err
		}
		w.project = project
	}

	connectedBoards := w.client.ConnectedBoards()
	allBoards := w.client.AllBoards()

	if w.target == nil {
		targetOpts := NewTargetOpts{
			ConnectedBoards: connectedBoards,
			AllBoards:       allBoards,
			OnlyConnected:   true,
			FQBN:            "",
			Logger:          w.logger,
		}
		target, err := NewTarget(targetOpts)
		if err != nil {
			return err
		}
		w.target = target
	}

	if w.port == nil {
		w.port = NewArdiSerialPort(w.target.Board.Port, w.project.Baud, w.logger)
	}

	w.buildProps = props

	return nil
}

// Upload compiled sketches to the specified board
func (w *WatchCore) Upload() error {
	w.port.Stop()
	w.waitForPreviousCompile()
	w.waitForPreviousUpload()
	w.logger.Info("Uploading...")

	if err := w.uploader.Upload(*w.target, w.project.Directory); err != nil {
		w.logger.WithError(err).Error("Failed to upload sketch")
		return err
	}

	return nil
}

// Compile the specified sketch
func (w *WatchCore) Compile() error {
	w.port.Stop()
	w.waitForPreviousCompile()
	w.waitForPreviousUpload()

	opts := rpc.CompileOpts{
		FQBN:       w.target.Board.FQBN,
		SketchDir:  w.project.Directory,
		SketchPath: w.project.Sketch,
		ExportName: "",
		BuildProps: w.buildProps,
		ShowProps:  false,
	}
	if err := w.compiler.Compile(opts); err != nil {
		w.logger.WithError(err).Error("Failed to compile sketch")
		return err
	}

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
		if !w.uploader.IsUploading() {
			break
		}
		w.logger.Info("Waiting for previous upload to finish...")
		time.Sleep(time.Second)
	}
}

func (w *WatchCore) waitForPreviousCompile() {
	// block until target is no longer compiling
	for {
		if !w.compiler.IsCompiling() {
			break
		}
		w.logger.Info("Waiting for previous compile to finish...")
		time.Sleep(time.Second)
	}
}
