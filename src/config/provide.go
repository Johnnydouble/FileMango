package config

import (
	"log"
	"os/user"
)

var configObj Config
var computedConfigObj ComputedConfig

func InitConfig(cfgFile string) {
	configObj = loadConfig(cfgFile)
	for i, dir := range configObj.UserConfig.Directories {
		if dir == "~" {
			configObj.UserConfig.Directories[i] = getUserHome() + "/"
		}
	}

	computedConfigObj = computeConfig()
}

func GetConfig() Config {
	return configObj
}

func GetComputedConfig() ComputedConfig {
	return computedConfigObj
}

//non public functions

func getUserHome() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}
