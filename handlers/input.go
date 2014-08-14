package handlers

// Represents user input.
type Input struct {
	Cmd   string
	Value string
}

// Creates a new Input type.
func New(cmd, value string) *Input {
	i := Input{cmd, value}
	return &i
}

func NewFromArray(parts []string) *Input {
	p_len := len(parts)
	cmd, value := "", ""

	if p_len > 0 {
		cmd = parts[0]
	}
	if p_len > 1 {
		value = parts[1]
	}

	return New(cmd, value)
}

func (i *Input) Reset() {
	i.Cmd, i.Value = "", ""
}
