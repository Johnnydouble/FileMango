package main

import (
	"fmt"
	"os"
)

func processFile(path string) {
	var qFile, _ = os.OpenFile("./fileQueue.txt", os.O_RDWR, 600) //rw for user, nothing for group and everyone
	defer qFile.Close()
	if _, err := qFile.WriteString(path + "\n"); err != nil {
		fmt.Println("WriteString failed, my final message")
		panic("panik!")
	}
	fmt.Println("Discovered File:", path)
}

//func confirmFiletype(path string) {
//	//sections := strings.Split(path, "/")
//
//	cmd := exec.Command("/usr/bin/file", "--brief", path)
//
//	var out bytes.Buffer
//	cmd.Stdout = &out
//
//	err := cmd.Run()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("filetype:", out.String())
//}
