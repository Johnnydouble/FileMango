package watch

import (
	"FileMango/src/config"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
)

func Create(rootDir string) {
	fileTypes = config.GetComputedConfig().Types

	/*CREATE INITIAL FILE QUEUE*/

	/*WATCH FS FOR UPDATES*/

	//create a new file watcher
	watcher, _ = fsnotify.NewWatcher()
	defer func() { _ = watcher.Close() }()

	//starting at the root of the project, walk each file/directory searching for directories
	if err := filepath.Walk(rootDir, watchDir); err != nil {
		fmt.Println("ERROR", err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				handleEvent(event, watcher)

			// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("WATCH ERROR: ", err)
			}
		}
	}()

	<-done
}

//runs as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {
	//todo:handle err from fs watcher.Walk and stat
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if !shouldWatch(path) {
		return nil
	}

	if fi.IsDir() {
		return watcher.Add(path)
	}
	return nil
}

func handleEvent(event fsnotify.Event, watcher *fsnotify.Watcher) {
	switch event.Op {
	case fsnotify.Create:
		handleChange(event)
	case fsnotify.Remove:
		// do nothing
	case fsnotify.Write:
		handleChange(event)
	}
}

func handleChange(event fsnotify.Event) {
	path := event.Name
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Println("File Stopped Existing Between Event and os.Stat:", err)
		return
	}

	queueFile(path)

	_ = watchDir(path, fileInfo, err)
}
