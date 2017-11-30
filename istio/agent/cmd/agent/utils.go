package main

import (
	"context"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/golang/glog"
)

// FSWatcher watches files and directories on the file system for changes.
type FSWatcher struct {
	names    []string
	changeCh chan string
	watcher  *fsnotify.Watcher
}

// NewFSWatcher creates a new watcher.
func NewFSWatcher(names []string) (*FSWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create fsnotify.Watcher: %v", err)
	}

	for _, name := range names {
		err := watcher.Add(name)
		if err != nil {
			return nil, fmt.Errorf("failed to add a name to fsnotify.Watcher: %v", err)
		}
	}

	return &FSWatcher{
		names:    names,
		changeCh: make(chan string),
		watcher:  watcher,
	}, nil
}

// Run starts the watcher.
func (w *FSWatcher) Run(ctx context.Context) {
	for {
		select {
		case event := <-w.watcher.Events:
			glog.V(2).Infof("Got a filesystem event: %v", event.String())
			w.changeCh <- event.String()
		case err := <-w.watcher.Errors:
			glog.Errorf("Error when watching for filesystem changes: %v", err)
		case <-ctx.Done():
			glog.V(2).Info("Terminating FSWatcher")
			w.watcher.Close()
			return
		}
	}

}

// Changes returns a chanel through which the watcher notifies about changes on the file system.
func (w *FSWatcher) Changes() <-chan string {
	return w.changeCh
}
