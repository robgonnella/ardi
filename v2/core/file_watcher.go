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

	// Run select statement for fsnotify events in separate goroutine
	// as recommended here : https://github.com/fsnotify/fsnotify#faq
	go func() {
		for {
			if f.watcher == nil {
				f.logger.Debug("file watcher already closed")
				return
			}

			select {
			case event, ok := <-f.watcher.Events:
				if !ok {
					f.logger.Error("unknown file watch error")
					return
				}

				f.logger.Debugf("event: %+v", event)

				writeEvt := event.Op&fsnotify.Write == fsnotify.Write
				removeEvt := event.Op&fsnotify.Remove == fsnotify.Remove

				// Remove evts are fired as the last event in atomic updates
				// i.e. mv file.tmp file
				// In this case we want to make sure the file still exists after
				// the remove by re-adding it to the watcher
				if removeEvt {
					if err := f.watcher.Add(f.file); err != nil {
						f.logger.WithError(err).Error("file watcher error")
						f.close <- true
						return
					}
				}

				if writeEvt || removeEvt {
					f.executeListeners()
				}
			case err, _ := <-f.watcher.Errors:
				f.logger.WithError(err).Error("Watch error")
				f.close <- true
				return
			}
		}
	}()

	// Block and wait for requests
	for {
		select {
		case <-f.stop:
			if f.watcher != nil {
				f.watcher.Remove(f.file)
				f.stop <- true
			}
		case <-f.restart:
			if f.watcher != nil {
				f.watcher.Add(f.file)
				f.restart <- true
			}
		case <-f.close:
			if f.watcher != nil {
				f.watcher.Remove(f.file)
				f.watcher.Close()
				f.watcher = nil
			}
			return nil
		}
	}
}

func (f *FileWatcher) executeListeners() {
	for _, l := range f.listeners {
		go l()
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
