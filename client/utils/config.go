package utils

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigFile struct {
	Minion struct {
		Name string
	}
	Master struct {
		Host         []string
		KeeplivePing int64
	}
	Log struct {
		Filename string // Location of log file
		LogLevel string // INFO, ERROR, DEBUG
	}
}

var Config *ConfigFile

func NewConfigFile() *ConfigFile {
	tmp := ConfigFile{}
	return &tmp
}

func (c *ConfigFile) ConfigInitialize() {

	f, err := os.Open("config.yml")
	if err != nil {
		fmt.Println("No configuration file found")
		f.Close()
	} else {
		defer f.Close()
		decoder := yaml.NewDecoder(f)
		err = decoder.Decode(&c)
		if err != nil {
			fmt.Printf("Configuration file error: %s", err.Error())
		}
	}
}
