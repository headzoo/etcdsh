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

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/coreos/go-etcd/etcd"
	"github.com/headzoo/etcdsh/config"
	"github.com/headzoo/etcdsh/etcdsh"
	"github.com/headzoo/etcdsh/handlers"
)

// Main method.
func main() {
	conf := config.New()

	help, version := false, false
	flag.BoolVar(&help, "help", false, "Prints command line options and exit.")
	flag.BoolVar(&version, "version", false, "Prints the etcdsh version and exit.")
	flag.StringVar(&conf.Machine, "machine", conf.Machine, "Connect to this etcd server.")
	flag.StringVar(&conf.PS1, "ps1", conf.PS1, "First prompt format")
	flag.StringVar(&conf.PS2, "ps2", conf.PS2, "Second prompt format")
	flag.BoolVar(&conf.Colors, "colors", conf.Colors, "Use colors in display.")
	flag.Parse()

	if help {
		printHelp()
		os.Exit(0)
	}
	if version {
		printVersion()
		os.Exit(0)
	}

	fmt.Printf("Connecting to %s\n", conf.Machine)
	client := etcd.NewClient([]string{conf.Machine})

	controller := handlers.NewController(conf, client, os.Stdout, os.Stderr, os.Stdin)
	controller.Add(handlers.NewLsHandler(controller))
	controller.Add(handlers.NewSetHandler(controller))
	controller.Add(handlers.NewHelpHandler(controller))
	controller.Add(handlers.NewCdHandler(controller))
	controller.Add(handlers.NewGetHandler(controller))
	code := controller.Start()

	os.Exit(code)
}

// printHelp prints the command line help information.
func printHelp() {
	printVersion()
	fmt.Println("USAGE:")
	fmt.Println("\tetcdsh [OPTIONS]")

	fmt.Println("")
	fmt.Println("OPTIONS:")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("\t-%-10s%s\n", f.Name, f.Usage)
	})

	fmt.Println("")
	fmt.Println("EXAMPLES:")
	fmt.Println("\tetcdsh -machine='http://192.168.1.23:4001'")

	fmt.Println("")
}

// printVersion prints the app version information.
func printVersion() {
	fmt.Printf("etcdsh %s - An interactive shell for the etcd server.\n\n", etcdsh.Version)
}
