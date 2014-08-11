package handlers

import "github.com/headzoo/etcdsh/io"

// ExitHandler handles the "exit" command.
type ExitHandler struct {
	controller *Controller
}

// NewExitHandler returns a new ExitHandler instance.
func NewExitHandler(controller *Controller) *ExitHandler {
	h := new(ExitHandler)
	h.controller = controller

	return h
}

// Command returns the string typed by the user that triggers to handler.
func (h *ExitHandler) Command() string {
	return "q"
}

// Validate returns whether the user input is valid.
func (h *ExitHandler) Validate(i *io.Input) bool {
	return true
}

// Syntax returns a string that demonstrates how to use the command.
func (h *ExitHandler) Syntax() string {
	return "q"
}

// Description returns a string that describes the command.
func (h *ExitHandler) Description() string {
	return "Exits the application"
}

// Handles the "ls" command.
func (h *ExitHandler) Handle(i *io.Input) (string, error) {
	return ExitSignal, nil
}
