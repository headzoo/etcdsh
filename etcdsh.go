package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/coreos/go-etcd/etcd"
	"github.com/headzoo/etcdsh/config"
	"github.com/headzoo/etcdsh/handlers"
	"github.com/headzoo/etcdsh/io"
)

// Main method.
func main() {
	config := config.New()

	help := flag.Bool("help", false, "Prints command line options and exit.")
	version := flag.Bool("version", false, "Prints the etcdsh version and exit.")
	flag.StringVar(&config.Machine, "machine", config.Machine, "Connect to this etcd server.")
	flag.BoolVar(&config.Colors, "colors", config.Colors, "Use colors in display.")
	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}
	if *version {
		printVersion()
		os.Exit(0)
	}

	fmt.Printf("Connecting to %s\n", config.Machine)
	client := etcd.NewClient([]string{config.Machine})

	controller := handlers.NewController(config, client, io.Stdout, io.Stderr, io.Stdin)
	controller.Add(handlers.NewLsHandler(controller))
	controller.Add(handlers.NewSetHandler(controller))
	controller.Add(handlers.NewExitHandler(controller))
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
	fmt.Printf("etcdsh %s - An interactive shell for the etcd server.\n\n", handlers.Version)
}
