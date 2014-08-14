package handlers

// CdHandler handles the "exit" command.
type CdHandler struct {
	CommandHandler
}

// NewCdHandler returns a new ExitHandler instance.
func NewCdHandler(controller *Controller) *CdHandler {
	h := new(CdHandler)
	h.controller = controller

	return h
}

// Command returns the string typed by the user that triggers to handler.
func (h *CdHandler) Command() string {
	return "cd"
}

// Validate returns whether the user input is valid.
func (h *CdHandler) Validate(i *Input) bool {
	return i.Key != ""
}

// Syntax returns a string that demonstrates how to use the command.
func (h *CdHandler) Syntax() string {
	return "cd <directory>"
}

// Description returns a string that describes the command.
func (h *CdHandler) Description() string {
	return "Changes the working directory"
}

// Handles the "cd" command.
func (h *CdHandler) Handle(i *Input) (string, error) {
	h.controller.ChangeWorkingDir(i.Key)
	return "", nil
}