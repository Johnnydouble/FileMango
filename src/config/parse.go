package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

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

func computeConfig() ComputedConfig {
	fileTypes := configObj.ModuleConfig.FileTypes
	var Types []string
	var Modules []string

	for _, modType := range fileTypes {
		Types = append(Types, modType.Type)
	}

	for _, fileType := range fileTypes {
		modList := fileType.Modules
		for _, mod := range modList {
			if !Contains(Modules, mod) {
				Modules = append(Modules, mod)
			}
		}
	}

	return ComputedConfig{Types, Modules}
}

func Contains(list []string, x string) bool {
	for _, item := range list {
		if item == x {
			return true
		}
	}
	return false
}

//Config File Data Structures
type Config struct {
	UserConfig   UserConfig
	ModuleConfig ModuleConfig
}

type UserConfig struct {
	Directories []string
}

type ModuleConfig struct {
	FileTypes []FileType
}

//a supported filetype and the
type FileType struct {
	Type    string
	Modules []string
}

//runtime configuration struct
type ComputedConfig struct {
	Types   []string
	Modules []string
}
