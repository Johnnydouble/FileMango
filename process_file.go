package main

import (
	"fmt"
	"os"
	"sync"
)

//use mutex to ensure that writes happen synchronously
var mutex sync.Mutex

func processFile(path string) {
	mutex.Lock()
	var qFile, _ = os.OpenFile("./file_queue.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 600) //rw for user, nothing for group and everyone
	defer qFile.Close()
	defer mutex.Unlock()
	if _, err := qFile.WriteString(path + "\n"); err != nil {
		fmt.Println("WriteString to file queue failed")
		panic(err)
	}
	fmt.Println("Discovered File:", path)

	//todo: remove debug

}
