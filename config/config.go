package config

import (
	"encoding/json"
	"os"
	"os/user"
	"strconv"
)

const (
	EnvPrefix      = "ETCDSH_"
	DefaultMachine = "http://127.0.0.1:4001"
	DefaultColors  = false
	DefaultPS1     = "\\u@etcd:\\w\\$ "
	DefaultPS2     = "> "
)

// Represents configuration file values.
type Config struct {
	Machine string
	Colors  bool
	PS1     string
	PS2     string
}

// Creates a new Config instance.
func New() *Config {
	c := new(Config)
	c.Machine = getenvString("MACHINE", DefaultMachine)
	c.Colors = getenvBool("COLORS", DefaultColors)
	c.PS1 = getenvString("PS1", DefaultPS1)
	c.PS2 = getenvString("PS2", DefaultPS2)

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

// getenv returns the value of an environment variable as a string or the default when the variable
// is not set. The EnvPrefix constant is automatically prepended to the key.
func getenvString(key, def string) string {
	val := os.Getenv(EnvPrefix + key)
	if val != "" {
		def = val
	}

	return def
}

// getenv returns the value of an environment variable as a bool or the default when the variable
// is not set. The EnvPrefix constant is automatically prepended to the key.
func getenvBool(key string, def bool) bool {
	val := os.Getenv(EnvPrefix + key)
	var err error
	if val != "" {
		def, err = strconv.ParseBool(val)
		if err != nil {
			def = false
		}
	}

	return def
}
