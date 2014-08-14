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
	h := &HelpHandler{
		controller: controller,
	}

	return h
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
