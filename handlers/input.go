package handlers

// Represents user input.
type Input struct {
	Cmd   string
	Args  []string
}

// Creates a new Input type.
func NewInput(cmd string) *Input {
	i := Input{
		Cmd: cmd,
		Args: []string{},
	}
	
	return &i
}

func (i *Input) Reset() {
	i.Cmd = ""
	i.Args = []string{}
}
