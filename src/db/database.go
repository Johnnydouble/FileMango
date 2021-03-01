package db

import (
	"fmt"
	"github.com/prologic/bitcask"
	"os"
	"sync"
	"time"
)

var db *bitcask.Bitcask
var Ambassador = newAmb()
var mutex sync.Mutex
var path string

func Init(pathi string) {
	path = pathi
	mutex.Lock()
	defer mutex.Unlock()

	//create db if doesnt exist
	_, err := os.Stat(path)
	if err == nil {
		newDB, _ := bitcask.Open(path, bitcask.WithAutoRecovery(true), bitcask.WithSync(true))
		newDB.Merge()
		newDB.Close()
	}

	//open db proper
	db, _ = bitcask.Open(path, bitcask.WithAutoRecovery(true), bitcask.WithSync(true))
	db.Merge()

}

func FoldQueue(f func(key []byte) error) error {
	mutex.Lock()
	defer mutex.Unlock()
	return db.Fold(f)
}

func QueueFile(path string) {
	mutex.Lock()
	defer mutex.Unlock()
	defer db.Sync()
	if !db.Has([]byte(path)) {
		_ = db.Put([]byte(path), []byte(""))
	} else {
		fmt.Println("Write to db failed", path)
		return
	}
	//scheduler.AddJob(path)
	Ambassador.Path <- path
}

func DequeueFile(path string) {
	mutex.Lock()
	defer mutex.Unlock()
	defer db.Sync()
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

func newAmb() ambassador {
	ambChan := make(chan string)
	return ambassador{ambChan}
}
