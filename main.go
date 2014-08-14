/**
 * Requirements:
 * 	libreadline-dev
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

	help := flag.Bool("help", false, "Prints command line options and exit.")
	version := flag.Bool("version", false, "Prints the etcdsh version and exit.")
	flag.StringVar(&conf.Machine, "machine", conf.Machine, "Connect to this etcd server.")
	flag.StringVar(&conf.PS1, "ps1", conf.PS1, "First prompt format")
	flag.StringVar(&conf.PS2, "ps2", conf.PS2, "Second prompt format")
	flag.BoolVar(&conf.Colors, "colors", conf.Colors, "Use colors in display.")
	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}
	if *version {
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
