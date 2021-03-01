package main

import (
	"FileMango/src/config"
	"FileMango/src/scheduler"
	"FileMango/src/watch"
	"fmt"
	"github.com/postfinance/single"
	"log"
)

var configPath = "./res/config.json"

// main
func main() {
	run()
}

func run() {
	one, err := single.New("FileMango", single.WithLockPath("/tmp"))
	if err != nil {
		log.Fatal(err)
	}
	// lock and defer unlocking
	err = one.Lock()
	if err != nil {
		log.Fatal(err)
	}
	defer one.Unlock()
	config.InitConfig(configPath)
	cfg := config.GetConfig()
	ComputedCfg := config.GetComputedConfig()
	fmt.Printf("config:%+v\ncomputed config:%+v\n", cfg, ComputedCfg)
	//todo: refactor create initial and create to take arrays of directories and replace this for loop with a call to each
	for _, dir := range cfg.UserConfig.Directories {
		go watch.CreateInitialFileQueue(dir)
		fmt.Println("watching...")
		go watch.Create(dir)
	}
	scheduler.RunAnalysis() //non functional right now
	//todo: WARNING: MAY CAUSE ISSUES WHEN DAEMONIZED
	_, _ = fmt.Scanln() //wait for keypress to exit
}
