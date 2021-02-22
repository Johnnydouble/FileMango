package main

import (
	"FileMango/src/config"
	"FileMango/src/scheduler"
	"FileMango/src/watch"
	"fmt"
)

var configPath = "./res/config.json"

// main
func main() {
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
	fmt.Scanln() //wait for keypress to exit
}
