package handlers

// Represents user input.
type Input struct {
	Cmd   string
	Key   string
	Value string
}

// Creates a new Input type.
func New(cmd, key, value string) *Input {
	i := Input{cmd, key, value}
	return &i
}

func NewFromArray(parts []string) *Input {
	p_len := len(parts)
	cmd, key, value := "", "", ""

	if p_len > 0 {
		cmd = parts[0]
	}
	if p_len > 1 {
		key = parts[1]
	}
	if p_len > 2 {
		value = parts[2]
	}

	return New(cmd, key, value)
}

func (i *Input) Reset() {
	i.Cmd, i.Key, i.Value = "", "", ""
}
