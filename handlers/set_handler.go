package handlers

import (
	"flag"
	"fmt"
)

// Command line options for the ls command.
type SetOptions struct {
	PrintHelp  bool
	TTL uint64
}

// SetHandler handles the "ls" command.
type SetHandler struct {
	controller *Controller
}

// NewSetHandler returns a new SetHandler instance.
func NewSetHandler(controller *Controller) *SetHandler {
	return &SetHandler{
		controller: controller,
	}
}

// Command returns the string typed by the user that triggers to handler.
func (h *SetHandler) Command() string {
	return "set"
}

// Validate returns whether the user input is valid.
func (h *SetHandler) Validate(i *Input) bool {
	return true
}

// Syntax returns a string that demonstrates how to use the command.
func (h *SetHandler) Syntax() string {
	return "set [options] <path> <value>"
}

// Description returns a string that describes the command.
func (h *SetHandler) Description() string {
	return "Sets the value of an object in the working directory"
}

// Handles the "ls" command.
func (h *SetHandler) Handle(i *Input) (string, error) {
	opts, args, err := h.setupOptions(i.Args)
	if opts == nil || err != nil {
		return "", err
	}
	
	resp, err := h.controller.Client().Set(args[0], args[1], opts.TTL)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n", resp.Node.Value), nil
}

// setupOptions builds a FlagSet and parses the args passed to the command.
func (h *SetHandler) setupOptions(args []string) (*SetOptions, []string, error) {
	opts := &SetOptions{}
	flags := flag.NewFlagSet("set_flags", flag.ContinueOnError)
	flags.BoolVar(&opts.PrintHelp, "h", false, "Show the command help")
	flags.Uint64Var(&opts.TTL, "t", 0, "Sets the TTL")

	err := flags.Parse(args)
	if err != nil {
		return nil, nil, err
	}

	args = flags.Args()
	if opts.PrintHelp || len(args) < 2 {
		printCommandHelp(h, flags)
		return nil, nil, nil
	}

	return opts, args, nil
}
