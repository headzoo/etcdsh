package handlers

import (
	"flag"
	"fmt"
)

// Handler types are called when a command is given by the user.
type Handler interface {
	Command() string
	Handle(*Input) (string, error)
	Validate(*Input) bool
	Syntax() string
	Description() string
}

// Represents a map of Handler instances
type HandlerMap map[string]Handler

// Base struct for other handlers.
type CommandHandler struct {
	Handler
	controller *Controller
	flags *flag.FlagSet
}

// printCommandHelp is used by handlers to display command help.
func printCommandHelp(handler Handler, flags *flag.FlagSet) {
	fmt.Println("SYNTAX")
	fmt.Println("\t" + handler.Syntax())
	fmt.Println("")
	fmt.Println("OPTIONS:")
	flags.VisitAll(func(f *flag.Flag) {
		fmt.Printf("\t-%-10s%s\n", f.Name, f.Usage)
	})
}
