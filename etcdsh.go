package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/coreos/go-etcd/etcd"
	"github.com/headzoo/etcdsh/handlers"
	"github.com/headzoo/etcdsh/io"
)

// Main method.
func main() {
	help := flag.Bool("help", false, "Prints command line options and exit.")
	machines := flag.String("machines", "http://127.0.0.1:4001", "Connect to these etcd servers. Defaults to 'http://127.0.0.1:4001'.")
	version := flag.Bool("version", false, "Prints the etcdsh version and exit.")
	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}
	if *version {
		printVersion()
		os.Exit(0)
	}

	client := etcd.NewClient([]string{*machines})
	controller := handlers.NewController(client, io.Stdout, io.Stderr, io.Stdin)
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
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("-%-10s%s\n", f.Name, f.Usage)
	})
}

// printVersion prints the app version information.
func printVersion() {
	fmt.Printf("etcdsh %s - An interactive shell for the etcd server.\n\n", handlers.Version)
}
