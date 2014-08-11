package config

import (
	"encoding/json"
	"os"
	"os/user"
)

const (
	DefaultMachine = "http://127.0.0.1:4001"
)

// Represents configuration file values.
type Config struct {
	Machine string
}

// Creates a new Config instance.
func New() *Config {
	c := new(Config)
	c.Machine = DefaultMachine

	usr, err := user.Current()
	if err == nil {
		configName := usr.HomeDir + "/.etcdsh"
		configFile, err := os.Open(configName)
		if err == nil {
			decoder := json.NewDecoder(configFile)
			err := decoder.Decode(c)
			if err != nil {
				panic(err)
			}
		}
	}

	return c
}
