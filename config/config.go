/**
The MIT License (MIT)

Copyright (c) 2014 Sean Hickey <sean@dulotech.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
 */

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
	conf := &Config{
		Machine: getenvString("MACHINE", DefaultMachine),
		Colors: getenvBool("COLORS", DefaultColors),
		PS1: getenvString("PS1", DefaultPS1),
		PS2: getenvString("PS2", DefaultPS2),
	}

	usr, err := user.Current()
	if err == nil {
		configName := usr.HomeDir + "/.etcdsh"
		configFile, err := os.Open(configName)
		if err == nil {
			decoder := json.NewDecoder(configFile)
			err := decoder.Decode(conf)
			if err != nil {
				panic(err)
			}
		}
	}

	return conf
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
