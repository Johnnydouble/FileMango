package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func HandleFlags() *bool {
	detach := flag.Bool("d", false,
		"Detach started process from the current tty then exit, by default the program starts as the current process")
	flag.Parse()
	return detach
}

func killOtherInstance() {
	fmt.Println("killing...")
	fmt.Println(os.Args[0])
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	path := filepath.Dir(ex)
	fmt.Println(path)
	cmd := exec.Command("pkill", path)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("success")
}
