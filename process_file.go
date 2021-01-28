package main

import (
	//"fmt"
	//"github.com/pkg/xattr"
	//"os"
	//"path/filepath"
	"bytes"
	"fmt"
	"os/exec"
)

func processFile(path string) {
	confirmFiletype(path)
}

func confirmFiletype(path string) {
	//sections := strings.Split(path, "/")

	cmd := exec.Command("/usr/bin/file", "--brief", path)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("filetype:", out.String())
}
