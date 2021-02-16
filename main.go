package main

import (
	"FileMango/src/config"
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
	for _, dir := range cfg.UserConfig.Directories {
		watch.CreateInitialFileQueue(dir)
		fmt.Println("watching...")
		watch.Create(dir)
	}
}
