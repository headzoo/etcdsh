/**
 * Parses command prompt definition strings.
 * See http://www.cyberciti.biz/tips/howto-linux-unix-bash-shell-setup-prompt.html
 */
package parser

import (
	"os/user"
	"strings"
	"bytes"
	"os"
	"path"
	"time"
)

// Used when the real host name cannot be determined.
const DefaultHostname = "etcd"

// Formatter returns a string associated with an escape sequence.
type Formatter func() string

// FormatterMap maps runes to Formatter functions.
type FormatterMap map[rune]Formatter

// Prompt is used to parse a prompt definition string.
type Prompt struct {
	formatters FormatterMap
}

// Creates and returns a new prompt parser instance.
func NewPrompt() *Prompt {
	p := new(Prompt)
	p.formatters = make(FormatterMap)
	p.formatters['\\'] = formatSlash
	p.formatters['a'] = formatBell
	p.formatters['d'] = formatTime
	p.formatters['e'] = formatEscape
	p.formatters['h'] = formatHostnameShort
	p.formatters['H'] = formatHostnameLong
	p.formatters['n'] = formatNewline
	p.formatters['r'] = formatCarriageReturn
	p.formatters['s'] = formatShell
	p.formatters['t'] = formatTime24
	p.formatters['T'] = formatTime12
	p.formatters['@'] = formatTimeAmPm
	p.formatters['A'] = formatTime24Long
	p.formatters['u'] = formatUser
	p.formatters['$'] = formatUid

	return p
}

// AddFormatter adds a custom formatter to the prompt parser.
func (p *Prompt) AddFormatter(key rune, f Formatter) {
	p.formatters[key] = f
}

// Parse parses the given prompt definition.
func (p *Prompt) Parse(s string) (string, error) {
	buffer := bytes.Buffer{}
	escaped := false

	for _, ch := range s {
		if ch == '\\' && !escaped {
			escaped = true
		} else if escaped {
			escaped = false
			found, ok := p.formatters[ch]
			if ok {
				buffer.WriteString(found())
			} else {
				buffer.WriteRune(ch)
			}
		} else {
			buffer.WriteRune(ch)
		}
	}

	return buffer.String(), nil
}

// formatEscape handles the \e escape sequence.
// Returns an ASCII escape character.
func formatEscape() string {
	return "\x1B"
}

// formatNewline handles the \n escape sequence.
// Returns a newline.
func formatNewline() string {
	return "\n"
}

// formatCarriageReturn handles the \r escape sequence.
// Returns a carriage return.
func formatCarriageReturn() string {
	return "\r"
}

// formatSlash handles the \\ escape sequence.
// Returns a backslash.
func formatSlash() string {
	return "\\"
}

// formatBell handles the \a escape sequence.
// Returns an ASCII bell character (07).
func formatBell() string {
	return "\x07"
}

// formatShell handles the \s escape sequence.
// Returns the name of the shell, the basename of $0 (the portion following the final slash).
func formatShell() string {
	return path.Base(os.Args[0])
}

// formatHostnameShort handles the \h escape sequence.
// Returns the hostname up to the first '.'.
func formatHostnameShort() string {
	host := formatHostnameLong()
	parts := strings.Split(host, ".")
	return parts[0]
}

// formatHostnameLong handles the \H escape sequence.
// Returns the hostname.
func formatHostnameLong() string {
	host, err := os.Hostname()
	if err == nil {
		return host
	} else {
		return DefaultHostname
	}
}

// formatUser handles the \u escape sequence.
// Returns the username of the current user.
func formatUser() string {
	usr, err := user.Current()
	if err == nil {
		return strings.ToLower(usr.Name)
	} else {
		return DefaultHostname
	}
}

// formatUid handles the \$ escape sequence.
// Returns if the effective UID is 0, a #, otherwise a $.
func formatUid() string {
	usr, err := user.Current()
	if err == nil {
		if usr.Uid == "0" {
			return "#"
		} else {
			return "$"
		}
	} else {
		return "$"
	}
}

// formatTime handles the \d escape sequence.
// Returns the date in "Weekday Month Date" format (e.g., "Tue May 26").
func formatTime() string {
	return time.Now().Format("Mon Jan 2")
}

// formatTime24 handles the \t escape sequence.
// Returns the current time in 24-hour HH:MM:SS format.
func formatTime24() string {
	return time.Now().Format("15:04:05")
}

// formatTime12 handles the \T escape sequence.
// Returns the current time in 12-hour HH:MM:SS format.
func formatTime12() string {
	return time.Now().Format("03:04:05")
}

// formatTimeAmPm handles the \@ escape sequence.
// Returns the current time in 12-hour am/pm format.
func formatTimeAmPm() string {
	return time.Now().Format("03:04:05pm")
}

// formatTime24Long handles the \A escape sequence.
// Returns the current time in 24-hour HH:MM format.
func formatTime24Long() string {
	return time.Now().Format("15:04")
}

