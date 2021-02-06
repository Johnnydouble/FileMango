package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"os"
	"os/user"
)

var configPath = "./config.json"
var watcher *fsnotify.Watcher
var fileTypes = initFileTypes()

// main
func main() {
	config := loadConfig(configPath)
	config.UserConfig.Directories = append(config.UserConfig.Directories, getUserHome()+"/")
	fmt.Printf("%+v", config)
	for _, dir := range config.UserConfig.Directories {
		createInitialFileQueue(dir)
		watchHome(dir)
	}
}

func getUserHome() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func initFileTypes() []string {
	supportedExtensionsFilepath := "./supportedTypes"
	extFile, _ := os.Open(supportedExtensionsFilepath)
	defer extFile.Close()
	scanner := bufio.NewScanner(io.Reader(extFile))
	extensions := make([]string, 0)
	for scanner.Scan() {
		extensions = append(extensions, scanner.Text())
	}
	return extensions
}

//information to be provided in the configuration file
type Config struct {
	UserConfig   UserConfig
	ModuleConfig ModuleConfig
}

type UserConfig struct {
	Directories []string
}
type ModuleConfig struct {
	Modules   []string
	FileTypes []FileType
}

//a supported filetype and the
type FileType struct {
	Type    string
	Modules []string
}

func loadConfig(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
