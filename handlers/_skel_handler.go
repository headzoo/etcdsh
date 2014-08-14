package handlers

// SkelHandler handles the "exit" command.
type SkelHandler struct {
	controller *Controller
}

// NewSkelHandler returns a new ExitHandler instance.
func NewSkelHandler(controller *Controller) *SkelHandler {
	h := &SkelHandler{
		controller: controller,
	}

	return h
}

// Command returns the string typed by the user that triggers to handler.
func (h *SkelHandler) Command() string {
	return "skel"
}

// Validate returns whether the user input is valid.
func (h *SkelHandler) Validate(i *Input) bool {
	return true
}

// Syntax returns a string that demonstrates how to use the command.
func (h *SkelHandler) Syntax() string {
	return "skel"
}

// Description returns a string that describes the command.
func (h *SkelHandler) Description() string {
	return "Skel description"
}

// Handles the "skel" command.
func (h *SkelHandler) Handle(i *Input) (string, error) {
	return "", nil
}
