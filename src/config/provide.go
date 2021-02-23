package config

import (
	"log"
	"os/user"
)

var configObj Config
var computedConfigObj ComputedConfig
var fileTypes []FileAssociation

func InitConfig(cfgFile string) {
	configObj = loadConfig(cfgFile)
	for i, dir := range configObj.UserConfig.Directories {
		if dir == "~" {
			configObj.UserConfig.Directories[i] = getUserHome() + "/"
		}
	}

	computedConfigObj = computeConfig()
	fileTypes = computedConfigObj.FileTypes
}

func GetConfig() Config {
	return configObj
}

func GetComputedConfig() ComputedConfig {
	return computedConfigObj
}

func GetFileTypes() []FileAssociation {
	return fileTypes
}

//non public functions

func getUserHome() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}
