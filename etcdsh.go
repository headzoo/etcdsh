package main

import (
	"os"

	"github.com/coreos/go-etcd/etcd"
	"github.com/headzoo/etcdsh/handlers"
	"github.com/headzoo/etcdsh/io"
)

// Main method.
func main() {
	machines := []string{"http://127.0.0.1:4001"}
	client := etcd.NewClient(machines)
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
