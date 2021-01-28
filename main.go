package main

import (
	"bufio"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/xattr"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var watcher *fsnotify.Watcher
var extensions = initFileExtensions()

// main
func main() {
	fmt.Println(extensions)

	/*CREATE INITIAL FILE QUEUE*/
	rootDir := "/home/chrisdev/"
	createInitialFileQueue(rootDir)

	/*WATCH FS FOR UPDATES*/

	fmt.Println("watching...")

	// creates a new file watcher
	watcher, _ = fsnotify.NewWatcher()
	defer func() { _ = watcher.Close() }()

	// starting at the root of the project, walk each file/directory searching for
	// directories
	if err := filepath.Walk(rootDir, watchDir); err != nil {
		fmt.Println("ERROR", err)
	}

	//
	done := make(chan bool)

	//
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				watchStuff(event, watcher)

				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	<-done
}

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {
	//handle err from fs watcher.Walk and stat
	if err != nil {
		log.Fatal(err)
	}

	if !shouldWatch(path) {
		return nil
	}

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.IsDir() {
		return watcher.Add(path)
	}

	return nil
}

func watchStuff(event fsnotify.Event, watcher *fsnotify.Watcher) {
	fmt.Println("EVENT:", event)

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
	fileInfo, err := os.Stat(path) // file path
	if err != nil {
		fmt.Println("File Stopped Existing Between Event and os.Stat:", err)
		return
	}
	/*TESTING SECTION*/
	queueFile(path)
	//testing xattrs
	list, _ := xattr.List(path)
	fmt.Println("xattr(s):", list)
	/*END TESTING SECTION*/

	_ = watchDir(path, fileInfo, err)
}

func shouldWatch(path string) bool {
	dir, _ := filepath.Split(path)
	sections := strings.Split(dir, "/")

	//deny certain home prefixes
	deniedHomePrefixes := []string{".", "snap", ".zsh_history.LOCK"}
	for _, prefix := range deniedHomePrefixes {
		if strings.HasPrefix(sections[3], prefix) {
			return false
		}
	}

	//deny certain path sections
	deniedSections := []string{".cache", ".git", ".fingerprint"}
	for _, section := range sections {
		for _, deniedSection := range deniedSections {
			if section == deniedSection {
				return false
			}
		}
	}
	return true
}

func createInitialFileQueue(rootDir string) {
	/*OPEN OR CREATE QUEUE FILE*/

	_ = filepath.Walk(rootDir, func(path string, fi os.FileInfo, err error) error {
		queueFile(path)
		return err
	})

}

func queueFile(path string) bool {
	if !shouldWatch(path) {
		return false
	}
	sections := strings.Split(path, "/")

	//allow file names with certain extensions
	for _, extension := range extensions {
		if strings.HasSuffix(sections[len(sections)-1], extension) {
			var qFile, _ = os.OpenFile("./fileQueue", os.O_RDWR, 600) //rw for user, nothing for group and everyone
			defer qFile.Close()
			if _, err := qFile.WriteString(path + "\n"); err != nil {
				fmt.Println("WriteString failed, my final message")
				panic("panik!")
			}
			fmt.Println("Discovered File:", path)
			//eventually this should be replaced with something that uses the fileQueue to ensure that unexpected stops are handled gracefully
			go processFile(path)
			return true
		}

	}
	return false
}

func initFileExtensions() []string {
	supportedExtensionsFilepath := "./supportedExtensions"
	extFile, _ := os.Open(supportedExtensionsFilepath)
	defer extFile.Close()
	scanner := bufio.NewScanner(io.Reader(extFile))
	extensions := make([]string, 0)
	for scanner.Scan() {
		extensions = append(extensions, scanner.Text())
	}
	return extensions
}
