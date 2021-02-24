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
		if dir[(len(dir)-1)] != '/' {
			configObj.UserConfig.Directories[i] = configObj.UserConfig.Directories[i] + "/"
		}
		if dir[0] == '~' {
			configObj.UserConfig.Directories[i] = getUserHome() +
				configObj.UserConfig.Directories[i][1:len(configObj.UserConfig.Directories[i])]
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
