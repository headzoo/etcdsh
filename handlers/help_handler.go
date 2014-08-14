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
	"bytes"
	"fmt"
)

// HelpHandler handles the "exit" command.
type HelpHandler struct {
	controller *Controller
}

// NewHelpHandler returns a new ExitHandler instance.
func NewHelpHandler(controller *Controller) *HelpHandler {
	return &HelpHandler{
		controller: controller,
	}
}

// Command returns the string typed by the user that triggers to handler.
func (h *HelpHandler) Command() string {
	return "help"
}

// Validate returns whether the user input is valid.
func (h *HelpHandler) Validate(i *Input) bool {
	return true
}

// Syntax returns a string that demonstrates how to use the command.
func (h *HelpHandler) Syntax() string {
	return "help"
}

// Description returns a string that describes the command.
func (h *HelpHandler) Description() string {
	return "Shows the help page"
}

// Handles the "ls" command.
func (h *HelpHandler) Handle(i *Input) (string, error) {
	handlers := h.controller.Handlers()
	buffer := bytes.NewBufferString("Etcdsh - An interactive shell for the etcd server.\n")

	buffer.WriteString("\nCOMMANDS\n")
	for key := range handlers {
		buffer.WriteString(fmt.Sprintf("\t%s - %s\n", key, handlers[key].Description()))
	}

	buffer.WriteString("\nSYNTAX:\n")
	for key := range handlers {
		buffer.WriteString(fmt.Sprintf("\t%s\n", handlers[key].Syntax()))
	}

	return buffer.String(), nil
}
