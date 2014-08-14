package handlers

// Handler types are called when a command is given by the user.
type Handler interface {
	Command() string
	Handle(*Input) (string, error)
	Validate(*Input) bool
	Syntax() string
	Description() string
}

// Represents a map of Handler instances
type HandlerMap map[string]Handler
