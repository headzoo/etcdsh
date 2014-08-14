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
