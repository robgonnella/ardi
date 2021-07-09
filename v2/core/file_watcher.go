package core

import (
	"errors"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

// Listener represents and action to run on file change
type Listener = func()

// FileWatcher watches sketch files for changes and runs user defined actions
type FileWatcher struct {
	file      string
	stop      chan bool
	close     chan bool
	restart   chan bool
	watcher   *fsnotify.Watcher
	listeners []Listener
	logger    *log.Logger
}

// NewFileWatcher returns a new sketch watcher instance
func NewFileWatcher(file string, logger *log.Logger) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = watcher.Add(file)
	if err != nil {
		return nil, err
	}

	return &FileWatcher{
		file:      file,
		watcher:   watcher,
		stop:      make(chan bool, 1),
		close:     make(chan bool, 1),
		restart:   make(chan bool, 1),
		listeners: []Listener{},
		logger:    logger,
	}, nil
}

// AddListener adds an action to run on file change
func (f *FileWatcher) AddListener(l Listener) {
	f.listeners = append(f.listeners, l)
}

// Watch watches the file for changes and runs user defined actions
func (f *FileWatcher) Watch() error {
	if f.watcher == nil {
		return errors.New("file watcher already closed")
	}

	f.logger.Infof("Watching %s for changes", f.file)

	for {
		if f.watcher == nil {
			f.logger.Debug("file watcher already closed")
			return nil
		}

		select {
		case <-f.stop:
			f.watcher.Remove(f.file)
			f.stop <- true
		case <-f.restart:
			f.watcher.Add(f.file)
			f.restart <- true
		case <-f.close:
			f.watcher.Remove(f.file)
			f.watcher.Close()
			f.watcher = nil
			return nil
		case event, ok := <-f.watcher.Events:
			if !ok {
				f.logger.Error("unknown file watch error")
				return nil
			}
			f.logger.Debugf("event: %+v", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				f.logger.Debugf("modified file: %s", event.Name)
				for _, l := range f.listeners {
					go l()
				}
			}
		case err, ok := <-f.watcher.Errors:
			if !ok {
				f.logger.Debug("unknown file watch error")
				return nil
			}
			f.watcher.Close()
			f.watcher = nil
			f.logger.WithError(err).Error("Watch error")
			return err
		}
	}
}

// Stop stops the os watcher
func (f *FileWatcher) Stop() {
	loggerWithField := f.logger.WithField("file", f.file)
	if f.watcher != nil {
		loggerWithField.Info("Stopping file watcher")
		f.stop <- true
		<-f.stop
	}
	loggerWithField.Info("File watcher stopped")
}

// Restart restarts the file watcher after being stopped
func (f *FileWatcher) Restart() {
	loggerWithField := f.logger.WithField("file", f.file)
	if f.watcher != nil {
		loggerWithField.Info("Restarting file watcher")
		f.restart <- true
		<-f.restart
	}
	loggerWithField.Info("File watcher restarted")
}

// Close fully closes file watcher (cannot be restarted)
func (f *FileWatcher) Close() {
	loggerWithField := f.logger.WithField("file", f.file)
	if f.watcher != nil {
		loggerWithField.Info("Closing file watcher")
		f.close <- true
		<-f.close
	}
	loggerWithField.Info("file watcher closed")
}
