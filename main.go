package main

import (
	"FileMango/src/cli"
	"FileMango/src/config"
	"FileMango/src/db"
	"FileMango/src/scheduler"
	"FileMango/src/watch"
	"fmt"
)

var configPath = "./res/config.json"
var queuePath = "./res/file_queue"

// main
func main() {
	run()
}

func run() {
	//make sure only one instance is running
	cli.Single()
	//load the config
	fmt.Print("loading config...") //todo: remove debug?
	config.InitConfig(configPath)
	cfg := config.GetConfig()
	fmt.Println(" done") //todo: remove debug?
	//init queue database
	fmt.Print("loading database...") //todo: remove debug?
	db.Init(queuePath)
	fmt.Println(" done") //todo: remove debug?
	//assemble the filequeue
	go watch.QueueExistingFiles(cfg.UserConfig.Directories)
	fmt.Println("watching...") //todo: remove debug?
	go watch.Create(cfg.UserConfig.Directories)
	//run analysis on the files in the queue
	scheduler.RunAnalysis()
	//keep the program from exiting
	cli.HandleSignal()
}
