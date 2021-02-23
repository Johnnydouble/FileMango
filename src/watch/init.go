package watch

import (
	"FileMango/src/config"
	"os"
	"path/filepath"
)

func CreateInitialFileQueue(rootDir string) {
	fileTypes = getFieldSlice(config.GetFileTypes())
	/*OPEN OR CREATE QUEUE FILE*/

	_ = filepath.Walk(rootDir, func(path string, fi os.FileInfo, err error) error {
		queueFile(path)
		return err
	})

}
