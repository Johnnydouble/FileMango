package main

import (
	"FileMango_AAS/src/config"
	"FileMango_AAS/src/watch"
	"fmt"
)

var configPath = "./res/config.json"

// main
func main() {
	config.InitConfig(configPath)
	cfg := config.GetConfig()
	ComputedCfg := config.GetComputedConfig()
	fmt.Printf("config:%+v\ncomputed config:%+v\n", cfg, ComputedCfg)
	for _, dir := range cfg.UserConfig.Directories {
		watch.CreateInitialFileQueue(dir)
		fmt.Println("watching...")
		watch.WatchHome(dir)
	}
}
