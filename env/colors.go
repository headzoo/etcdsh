package env

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Holds information about system colors.
type Colors struct {
	ls_colors map[string]string
}

// Creates and returns a new Colors instance.
func NewColors() *Colors {
	c := new(Colors)

	return c
}

// GetLS returns a color value defined in the LS_COLORS environment variable.
// Use the key to get a specific color, eg "di", "fi", "ln", etc.
// See http://blog.twistedcode.org/2008/04/lscolors-explained.html
func (c *Colors) GetLS(key string) (string, error) {
	if c.ls_colors == nil {
		c.ls_colors = make(map[string]string)
		value := os.Getenv("LS_COLORS")
		if value != "" {
			parts := strings.Split(value, ":")
			for _, part := range parts {
				op := strings.SplitN(part, "=", 2)
				if len(op) == 2 {
					c.ls_colors[op[0]] = op[1]
				}
			}
		}
	}

	_, ok := c.ls_colors[key]
	if ok {
		return c.ls_colors[key], nil
	}
	return "", errors.New(fmt.Sprintf("No color value defined for key '%s'.", key))
}

// GetLSDefault works exactly like GetLS, but you may provide a default value if the color
// is not set.
func (c *Colors) GetLSDefault(key, def string) (string, error) {
	value, err := c.GetLS(key)
	if err != nil {
		value = def
	}

	return value, nil
}
