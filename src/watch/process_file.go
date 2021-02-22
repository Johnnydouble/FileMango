package watch

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
)

//use mutex to ensure that writes happen synchronously
var mutex sync.Mutex

func ProcessFile(path string) {
	mutex.Lock()
	defer mutex.Unlock()
	var qFile, _ = os.OpenFile("./res/file_queue.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 600) //rw for user, nothing for group and everyone
	defer qFile.Close()

	//todo: Remove duplicates
	scanner := bufio.NewScanner(io.Reader(qFile))

	for scanner.Scan() {
		fmt.Println(scanner.Text())
		if scanner.Text() == path {
			fmt.Println("duplicate event detected")
			return
		}
	}

	if _, err := qFile.WriteString(path + "\n"); err != nil {
		fmt.Println("WriteString to file queue failed")
		panic(err)
	}
	fmt.Println("Discovered File:", path)

}
