package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Config represent a set of configurations for this app.
type Config struct {
	Host     string `json:"host"`
	Port     int16  `json:"port"`
	Database string `json:"database"`
}

const configPath = "config.json"

var config Config

func readConfigFile() Config {
	var confFile *os.File
	var byteValue []byte
	var err error
	var conf Config

	defaultConf := Config{
		Host:     "localhost",
		Port:     3000,
		Database: "mongodb://localhost/api-exercise",
	}

	confFile, err = os.Open(configPath)
	if err != nil {
		log.Println("Cannot read config.json, use default configurations.")
		return defaultConf
	}
	defer confFile.Close()

	byteValue, err = ioutil.ReadAll(confFile)
	if err != nil {
		log.Println("Cannot read config.json, use default configurations.")
		return defaultConf
	}

	err = json.Unmarshal(byteValue, &conf)
	if err != nil {
		log.Println("Cannot read config.json, use default configurations.")
		return defaultConf
	}

	// Check required fields are exist
	if conf.Host == "" {
		log.Println("Field \"host\" undefined, use default value.")
		conf.Host = defaultConf.Host
	}
	if conf.Port == 0 {
		log.Println("Field \"port\" undefined, use default value.")
		conf.Port = defaultConf.Port
	}
	if conf.Database == "" {
		log.Println("Field \"database\" undefined, use default value.")
		conf.Database = defaultConf.Database
	}

	return conf
}

func init() {
	config = readConfigFile()
}

// GetConfig get app configurations
func GetConfig() *Config {
	return &config
}
