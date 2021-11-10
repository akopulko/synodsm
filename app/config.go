package main

import (
	"os/user"
	"path/filepath"

	"gopkg.in/ini.v1"
)

const (
	configFile = ".synodsm" // config file name, stored in $HOME
)

type ConfigFileData struct {
	Server string
	User   string
}

func getUserConfigFilePath() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(currentUser.HomeDir, configFile), nil

}

func saveConfig(server string, user string) error {

	cfgIni := ini.Empty()

	configFileData := &ConfigFileData{
		Server: server,
		User:   user,
	}

	// map ini to struct
	err := ini.ReflectFrom(cfgIni, configFileData)
	if err != nil {
		return err
	}

	configFilePath, err := getUserConfigFilePath()
	if err != nil {
		return err
	}

	err = cfgIni.SaveTo(configFilePath)
	if err != nil {
		return err
	}

	return nil

}

func loadConfig() (*ConfigFileData, error) {

	configFilePath, err := getUserConfigFilePath()
	if err != nil {
		return nil, err
	}

	configFileData := new(ConfigFileData)

	err = ini.MapTo(configFileData, configFilePath)
	if err != nil {
		return nil, err
	}

	return configFileData, nil

}
