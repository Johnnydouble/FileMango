package watch

import (
	"bytes"
	"github.com/fsnotify/fsnotify"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var watcher *fsnotify.Watcher
var fileTypes []string

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
			ProcessFile(path) //todo: consider moving code from this function into that one or vice versa
			return true
		}

	}
	return false
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
