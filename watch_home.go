package main

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func watchHome(rootDir string) {
	fmt.Println(rootDir)
	/*CREATE INITIAL FILE QUEUE*/

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
				handleEvent(event, watcher)

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

	if fi.IsDir() {
		return watcher.Add(path)
	}

	return nil
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

	fi, _ := os.Stat(path)
	if fi.IsDir() == true {
		return false
	}

	//allow file names with certain fileTypes
	for _, fileType := range fileTypes {
		if fileType == getFileType(path) {
			processFile(path)
			return true
		}

	}
	return false
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

func getFileType(path string) string {

	// Open File
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Get the content
	contentType, err := GetFileContentType(f)

	fallbackType := "application/octet-stream"
	if err != nil || contentType == fallbackType {
		contentType = handleUnknownContentType(f)
	}
	if contentType == "" {
		contentType = fallbackType
	}

	return contentType
}

func GetFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// default content type is "application/octet-stream"
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func handleUnknownContentType(f *os.File) string {
	//todo: change this to a pipeline model so that 'file' doesnt have to start each time its needed
	//todo: investigate why "," is being appended to the ASCII.text content type some of the time
	cmd := exec.Command("/usr/bin/file", "--brief", f.Name())

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	result := strings.Split(out.String(), " ")

	resultStr := ""
	switch len(result) {
	case 0:
	case 1:
		resultStr = result[0]
	default:
		resultStr = result[0] + "." + result[1]
	}

	return "filemango/" + resultStr
}
