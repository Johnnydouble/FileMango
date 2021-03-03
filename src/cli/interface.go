package cli

import (
	"FileMango/src/db"
	"flag"
	"fmt"
	"github.com/postfinance/single"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
)

var one *single.Single

func Single() {
	var err error
	one, err = single.New("FileMango", single.WithLockPath("/tmp"))
	if err != nil {
		log.Fatal(err)
	}
	err = one.Lock()
	if err != nil {
		log.Fatal(err)
	}
}

func HandleSignal() {
	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt)
	<-killSignal
	fmt.Print("Exiting... ")
	db.Close()
	_ = one.Unlock()
	fmt.Print("done.")
	os.Exit(0)
}

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
