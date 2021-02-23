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
	cfg := ComputedConfig{}
	modules := configObj.ModuleConfig.Modules
	for _, module := range modules {
		for _, fileType := range module.FileTypes {
			cfg.invertTree(fileType, module.Name, module.Path)
		}
	}
	return cfg
}

func (f *ComputedConfig) invertTree(ftype, name, path string) {
	for i, item := range f.FileTypes {
		//if an entry for this file type exists append this module to it and return
		if item.Type == ftype {
			f.FileTypes[i].ModuleNames = append(f.FileTypes[i].ModuleNames, name)
			f.FileTypes[i].ModulePaths = append(f.FileTypes[i].ModulePaths, path)
			return
		}
	}
	//if an entry for this file type does not exist create one
	f.FileTypes = append(f.FileTypes, FileAssociation{
		ftype,
		[]string{name},
		[]string{path},
	})
}

func (f *ComputedConfig) contains(x string) bool {
	for _, item := range f.FileTypes {
		if item.Type == x {
			return true
		}
	}
	return false
}

//structure representing the config file
type Config struct {
	UserConfig struct {
		Directories []string `json:"Directories"`
	} `json:"UserConfig"`
	ModuleConfig struct {
		Modules []struct {
			Name      string   `json:"Name"`
			Path      string   `json:"Path"`
			FileTypes []string `json:"FileTypes"`
		} `json:"Modules"`
	} `json:"ModuleConfig"`
}

//runtime configuration struct
type ComputedConfig struct {
	FileTypes []FileAssociation
}

type FileAssociation struct {
	Type        string
	ModuleNames []string
	ModulePaths []string
}
