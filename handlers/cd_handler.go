package handlers

// CdHandler handles the "exit" command.
type CdHandler struct {
	controller *Controller
}

// NewCdHandler returns a new ExitHandler instance.
func NewCdHandler(controller *Controller) *CdHandler {
	h := &CdHandler{
		controller: controller,
	}

	return h
}

// Command returns the string typed by the user that triggers to handler.
func (h *CdHandler) Command() string {
	return "cd"
}

// Validate returns whether the user input is valid.
func (h *CdHandler) Validate(i *Input) bool {
	return len(i.Args) > 0
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
	h.controller.ChangeWorkingDir(i.Args[0])
	return "", nil
}
