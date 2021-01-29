package main

import (
	"bufio"
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"os"
	"os/user"
)

var watcher *fsnotify.Watcher
var fileTypes = initFileTypes()

// main
func main() {
	rootDir := getUserHome() + "/"
	watchHome(rootDir)
	createInitialFileQueue(rootDir)
}

func getUserHome() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func initFileTypes() []string {
	supportedExtensionsFilepath := "./supportedTypes"
	extFile, _ := os.Open(supportedExtensionsFilepath)
	defer extFile.Close()
	scanner := bufio.NewScanner(io.Reader(extFile))
	extensions := make([]string, 0)
	for scanner.Scan() {
		extensions = append(extensions, scanner.Text())
	}
	return extensions
}
