package watch

import (
	"FileMango/src/config"
	"os"
	"path/filepath"
)

func QueueExistingFiles(Directories []string) {
	fileTypes = getFieldSlice(config.GetFileTypes())
	/*OPEN OR CREATE QUEUE FILE*/

	for _, dir := range Directories {
		_ = filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
			queueFile(path)
			return err
		})
	}
}
