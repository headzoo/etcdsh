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
	help := flag.Bool("help", false, "Show command line options and exit.")
	machines := flag.String("machines", "http://127.0.0.1:4001", "Set the etcd machine to connect to.")
	version := flag.Bool("version", false, "Show the version and exit.")
	flag.Parse()

	if *help {
		showHelp()
		os.Exit(0)
	}
	if *version {
		fmt.Printf("etcdsh %s - An interactive shell for the etcd server.\n\n", handlers.Version)
		os.Exit(0)
	}

	client := etcd.NewClient([]string{*machines})
	//client.Set("/go/etcdsh", "You caught me!", 0)

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

// showHelp displays the flag defaults. Later on this should be change to make the output a little more pretty.
func showHelp() {
	flag.PrintDefaults()
}
