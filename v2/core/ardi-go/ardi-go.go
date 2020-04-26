package ardigo

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/robgonnella/ardi/v2/core/project"
	"github.com/robgonnella/ardi/v2/core/serial"
	"github.com/robgonnella/ardi/v2/core/target"
	"github.com/robgonnella/ardi/v2/rpc"
	log "github.com/sirupsen/logrus"
)

// ArdiGo represents core module for adi go commands
type ArdiGo struct {
	logger     *log.Logger
	Client     *rpc.Client
	target     *target.Target
	project    *project.Project
	buildProps []string
	port       *serial.Port
	compiling  bool
	uploading  bool
}

// New returns new Project instance
func New(sketchDir string, buildProps []string, logger *log.Logger) (*ArdiGo, error) {
	proj, err := project.New(logger)
	if err != nil {
		return nil, err
	}
	if err := proj.ProcessSketch(sketchDir); err != nil {
		return nil, err
	}

	client, err := rpc.NewClient(logger)
	if err != nil {
		return nil, err
	}

	target, err := target.New(logger, "", true)
	if err != nil {
		return nil, err
	}

	port := serial.New(target.Board.Port, proj.Baud, logger)

	return &ArdiGo{
		Client:     client,
		project:    proj,
		target:     target,
		buildProps: buildProps,
		port:       port,
		logger:     logger,
	}, nil
}

// Upload compiled sketches to the specified board
func (a *ArdiGo) Upload() error {
	a.port.Stop()
	a.waitForPreviousCompile()
	a.waitForPreviousUpload()
	a.logger.Info("Uploading...")

	fqbn := a.target.Board.FQBN
	device := a.target.Board.Port
	sketchDir := a.project.Directory

	a.uploading = true
	if err := a.Client.Upload(fqbn, sketchDir, device); err != nil {
		a.logger.WithError(err).Error("Failed to upload sketch")
		a.uploading = false
		return err
	}

	a.uploading = false
	return nil
}

// Compile the specified sketch
func (a *ArdiGo) Compile() error {
	a.port.Stop()
	a.waitForPreviousCompile()
	a.waitForPreviousUpload()

	fqbn := a.target.Board.FQBN
	sketchDir := a.project.Directory
	sketch := a.project.Sketch
	buildProps := a.buildProps

	a.compiling = true
	if err := a.Client.Compile(fqbn, sketchDir, sketch, "", buildProps, false); err != nil {
		a.logger.WithError(err).Error("Failed to compile sketch")
		a.compiling = false
		return err
	}

	a.compiling = false
	return nil
}

// WatchLogs connects to a serial port at a specified baud rate and prints
// any logs received.
func (a *ArdiGo) WatchLogs() {
	baud := a.project.Baud
	device := a.target.Board.Port

	logFields := log.Fields{"baud": baud, "device": device}

	a.port.Stop()
	a.waitForPreviousCompile()
	a.waitForPreviousUpload()

	a.logger.WithFields(logFields).Info("Watching logs...")
	a.port.Watch()
}

// WatchSketch responds to changes in a given sketch file by automatically
// recompiling and re-uploading.
func (a *ArdiGo) WatchSketch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		a.logger.WithError(err).Error("Failed to watch directory for changes")
	}
	defer watcher.Close()

	err = watcher.Add(a.project.Sketch)
	if err != nil {
		a.logger.WithError(err).Error("Failed to watch directory for changes")
	}

	go a.WatchLogs()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				break
			}
			a.logger.Debugf("event: %+v", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				a.logger.Debugf("modified file: %s", event.Name)
				err := a.Compile()
				if err != nil {
					a.Upload()
					go a.WatchLogs()
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			a.logger.WithError(err).Warn("Watch error")
		}
	}

}

// private helpers
func (a *ArdiGo) waitForPreviousUpload() {
	// block until target is no longer uploading
	for {
		if !a.uploading {
			break
		}
		a.logger.Info("Waiting for previous upload to finish...")
		time.Sleep(time.Second)
	}
}

func (a *ArdiGo) waitForPreviousCompile() {
	// block until target is no longer compiling
	for {
		if !a.compiling {
			break
		}
		a.logger.Info("Waiting for previous compile to finish...")
		time.Sleep(time.Second)
	}
}
