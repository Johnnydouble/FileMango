package db

import (
	"fmt"
	"github.com/prologic/bitcask"
	"time"
)

var db *bitcask.Bitcask
var Ambassador = NewAmb()

func Init(path string) {
	db, _ = bitcask.Open(path)
	//db.Put([]byte("Hello"), []byte("World"))
	//_ = db.Put([]byte(path), []byte("World"))
}

func FoldQueue(f func(key []byte) error) error {
	return db.Fold(f)
}

func QueueFile(path string) {
	if !db.Has([]byte(path)) {
		_ = db.Put([]byte(path), []byte(""))
	} else {
		fmt.Println("Write to db failed")
		return
	}
	//scheduler.AddJob(path)
	Ambassador.Path <- path
}

func DequeueFile(path string) {
	err := db.Delete([]byte(path))
	if err != nil {
		fmt.Print("Failed to remove file from queue, retrying... ")
		time.Sleep(100 * time.Millisecond)
		err := db.Delete([]byte(path))
		if err != nil {
			fmt.Println("failed.")
			fmt.Println(err)
		} else {
			fmt.Println("succeeded.")
		}
	}
}

type ambassador struct {
	Path chan string
}

func NewAmb() ambassador {
	ambChan := make(chan string)
	return ambassador{ambChan}
}
