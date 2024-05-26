package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

var (
	watcherRunning bool
	watcher        *fsnotify.Watcher
	debounceMap    = make(map[string]*time.Timer)
	debounceMutex  sync.Mutex

	lastMoved string
	counter   int
)

func stopWatcher() {
	if !watcherRunning {
		return
	}

	watcherRunning = false
	if err := watcher.Close(); err != nil {
		logrus.Error("Error closing watcher: ", err)
	}
	logrus.Info("Watcher stopped")
}

func startWatcher() {
	logrus.Info("Starting watcher...")
	if watcherRunning {
		return
	}

	watcherRunning = true

	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		logrus.Error("Error creating watcher: ", err)
	}

	for _, watch := range config.Watchlist {
		err = watcher.Add(watch.WatchPath)
		if err != nil {
			logrus.Error("Error adding path to watcher: ", err)
			continue
		}
		logrus.Info("Watching: ", watch.WatchPath)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				handleEvent(event)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logrus.Error("Error: ", err)
			}
		}
	}()

	logrus.Info("Watcher initialized!")
}

func handleEvent(event fsnotify.Event) {
	logrus.Info("Event: ", event)
	logrus.Info("Modified file: ", event.Name)

	if event.Has(fsnotify.Remove) {
		return
	}

	for _, watch := range config.Watchlist {
		filename := filepath.Base(event.Name)

		matched, err := regexp.MatchString(watch.WatchPattern, filename)
		if err != nil {
			logrus.Error("Error matching pattern: ", err)
			continue
		}
		if matched {
			debounceEvent(event.Name, 2*time.Second, func() {
				logrus.Info("Moving file: ", event.Name, " to ", watch.DestinationPath)
				if moveErr := os.Rename(event.Name, path.Join(watch.DestinationPath, filepath.Base(event.Name))); moveErr != nil {
					logrus.Error("Error moving file: ", moveErr)
					return
				}

				lastMoved = event.Name
				counter++
				updateSystray()
			})
		}
	}
}

func debounceEvent(filePath string, delay time.Duration, action func()) {
	debounceMutex.Lock()
	defer debounceMutex.Unlock()

	if timer, found := debounceMap[filePath]; found {
		timer.Stop()
	}

	timer := time.AfterFunc(delay, func() {
		debounceMutex.Lock()
		delete(debounceMap, filePath)
		debounceMutex.Unlock()
		action()
	})

	debounceMap[filePath] = timer
}
