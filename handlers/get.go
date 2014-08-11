package handlers

import "github.com/headzoo/etcdsh/io"

// GetHandler handles the "exit" command.
type GetHandler struct {
	controller *Controller
}

// NewGetHandler returns a new ExitHandler instance.
func NewGetHandler(controller *Controller) *GetHandler {
	h := new(GetHandler)
	h.controller = controller

	return h
}

// Command returns the string typed by the user that triggers to handler.
func (h *GetHandler) Command() string {
	return "get"
}

// Validate returns whether the user input is valid.
func (h *GetHandler) Validate(i *io.Input) bool {
	return i.Key != ""
}

// Syntax returns a string that demonstrates how to use the command.
func (h *GetHandler) Syntax() string {
	return "get <key>"
}

// Description returns a string that describes the command.
func (h *GetHandler) Description() string {
	return "Displays the value of the given key"
}

// Handles the "skel" command.
func (h *GetHandler) Handle(i *io.Input) (string, error) {
	dir := h.controller.WorkingDir(i.Key)
	resp, err := h.controller.Client().Get(dir, false, false)
	if err != nil {
		return "", err
	}

	return resp.Node.Value + "\n", nil
}
