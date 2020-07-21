package core

import (
	"fmt"
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
	listeners []Listener
	sigs      chan os.Signal
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

	return &FileWatcher{
		file:      file,
		watcher:   watcher,
		listeners: []Listener{},
		sigs:      sigs,
		logger:    logger,
	}, nil
}

// AddListener adds an action to run on file change
func (f *FileWatcher) AddListener(l Listener) {
	f.listeners = append(f.listeners, l)
}

// Watch watches the file for changes and runs user defined actions
func (f *FileWatcher) Watch() {
	defer f.Close()

	go func() {
		<-f.sigs
		fmt.Println()
		fmt.Println("gracefully shutting down file watcher")
		f.Close()
	}()

	f.logger.Infof("Watching %s for changes", f.file)

	for {
		select {
		case event, ok := <-f.watcher.Events:
			if !ok {
				break
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
				return
			}
			f.logger.WithError(err).Warn("Watch error")
		}
	}
}

// Close closes the os watcher
func (f *FileWatcher) Close() {
	if f.watcher != nil {
		f.watcher.Close()
	}
}
