package handlers

// SetHandler handles the "ls" command.
type SetHandler struct {
	CommandHandler
}

// NewSetHandler returns a new SetHandler instance.
func NewSetHandler(controller *Controller) *SetHandler {
	h := new(SetHandler)
	h.controller = controller

	return h
}

// Command returns the string typed by the user that triggers to handler.
func (h *SetHandler) Command() string {
	return "set"
}

// Validate returns whether the user input is valid.
func (h *SetHandler) Validate(i *Input) bool {
	return i.Key != "" && i.Value != ""
}

// Syntax returns a string that demonstrates how to use the command.
func (h *SetHandler) Syntax() string {
	return "set <path> <value>"
}

// Description returns a string that describes the command.
func (h *SetHandler) Description() string {
	return "Sets the value of an object in the working directory"
}

// Handles the "ls" command.
func (h *SetHandler) Handle(i *Input) (string, error) {
	resp, err := h.controller.Client().Set(i.Key, i.Value, 0)
	if err != nil {
		return "", err
	}

	return resp.Node.Value + "\n", nil
}