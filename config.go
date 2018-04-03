package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	//Croc Congig
	Croc = viper.New()

	//API Config
	API = viper.New()
)

func setupConfig() {

	//Setting API config settings
	API.SetConfigName("crocgodyl")
	API.AddConfigPath("configs/")
	API.WatchConfig()

	API.OnConfigChange(func(e fsnotify.Event) {
		writeLog("info", "crocgodyl config changed", nil)
	})

	if err := API.ReadInConfig(); err != nil {
		writeLog("fatal", "Could not load crocgodyl configuration.", err)
		return
	}

	//Setting API config settings
	API.SetConfigName("api")
	API.AddConfigPath("configs/")
	API.WatchConfig()

	API.OnConfigChange(func(e fsnotify.Event) {
		writeLog("info", "API config changed", nil)
	})

	if err := API.ReadInConfig(); err != nil {
		writeLog("fatal", "Could not load API configuration.", err)
		return
	}
}

func getCrocConfigString(req string) string {
	return Croc.GetString(req)
}

func getAPIConfigString(req string) string {
	return API.GetString(req)
}
