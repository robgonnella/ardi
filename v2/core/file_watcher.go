package core

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

// Listener represents and action to run on file change
type Listener = func()

// FileWatcher watches sketch files for changes and runs user defined actions
type FileWatcher struct {
	file      string
	watcher   *fsnotify.Watcher
	destroyed bool
	pause     chan bool
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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	pause := make(chan bool, 1)

	fileWatcher := &FileWatcher{
		file:      file,
		watcher:   watcher,
		pause:     pause,
		listeners: []Listener{},
		logger:    logger,
	}

	go func() {
		<-sigs
		logger.Debug("gracefully shutting down file watcher")
		fileWatcher.Close()
	}()

	return fileWatcher, nil
}

// AddListener adds an action to run on file change
func (f *FileWatcher) AddListener(l Listener) {
	f.listeners = append(f.listeners, l)
}

// Watch watches the file for changes and runs user defined actions
func (f *FileWatcher) Watch() error {
	if f.destroyed {
		f.logger.Debug("file watcher has been destroyed")
		return nil
	}

	f.logger.Infof("Watching %s for changes", f.file)
	for {
		if f.destroyed {
			f.logger.Debug("file watcher has been destroyed")
			return nil
		}

		select {
		case <-f.pause:
			f.pause <- true
			return nil
		case event, ok := <-f.watcher.Events:
			if !ok {
				err := errors.New("Failed to watch")
				f.logger.WithError(err).Debug()
				return err
			}
			f.logger.Debugf("event: %+v", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				f.logger.Debugf("modified file: %s", event.Name)
				for _, l := range f.listeners {
					l()
				}
			}
		case err, ok := <-f.watcher.Errors:
			if !ok {
				return nil
			}
			f.logger.WithError(err).Warn("Watch error")
			return err
		}
	}
}

// Close closes the os watcher, watcher can not be restarted
func (f *FileWatcher) Close() {
	if f.watcher != nil {
		f.watcher.Close()
		f.watcher.Remove(f.file)
		f.destroyed = true
	}
}

// Stop stops watching file with ability to restart at any time
func (f *FileWatcher) Stop() {
	f.pause <- true
	<-f.pause
	f.watcher.Remove(f.file)
	f.logger.WithField("file", f.file).Debug("file watcher paused")
}

// Restart restarts existing watcher
func (f *FileWatcher) Restart() {
	f.watcher.Add(f.file)
	f.Watch()
}
