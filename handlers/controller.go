package handlers

import (
	"bufio"
	"fmt"
	"os"
	"io"
	"path"
	"strings"

	"github.com/coreos/go-etcd/etcd"
	"github.com/headzoo/etcdsh/config"
	eio "github.com/headzoo/etcdsh/io"
	"github.com/bobappleyard/readline"
	"bytes"
)

const (
	Version    = "0.2"
	ExitSignal = "__exit"
)

// Represents a map of Handler instances
type HandlerMap map[string]Handler

// Handler types are called when a command is given by the user.
type Handler interface {
	Command() string
	Handle(*eio.Input) (string, error)
	Validate(*eio.Input) bool
	Syntax() string
	Description() string
}

// Controller stores handlers and calls them.
type Controller struct {
	handlers              HandlerMap
	scanner               *bufio.Scanner
	config                *config.Config
	client                *etcd.Client
	stdout, stderr, stdin *os.File
	wdir                  string
}

// Create a new Controller.
func NewController(config *config.Config, client *etcd.Client, stdout, stderr, stdin *os.File) *Controller {
	c := new(Controller)
	c.handlers = make(HandlerMap)
	c.scanner = bufio.NewScanner(stdin)
	c.config = config
	c.client = client
	c.stdout, c.stderr, c.stdin = stdout, stderr, stdin
	c.wdir = "/"

	return c
}

// Starts the controller.
func (c *Controller) Start() int {
	c.welcome()
	
	prompt := ""
	buffer := bytes.NewBufferString("")
	
	for {
		if buffer.Len() == 0 {
			prompt = c.ps1()
		} else {
			prompt = c.ps2()
		}
		
		line, err := readline.String(prompt)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if strings.HasSuffix(line, "\\") {
			buffer.WriteString(strings.TrimSuffix(line, "\\") + "\n")
		} else {
			if buffer.Len() > 0 {
				buffer.WriteString(line)
				line = buffer.String()
				buffer.Reset()
			}

			parts := strings.SplitN(line, " ", 3)
			if parts[0] != "" {
				readline.AddHistory(line)
				should_exit := c.handleInput(eio.NewFromArray(parts))
				if should_exit {
					break
				}
			}
		}
	}

	return 0
}

// Client returns the etcd client
func (c *Controller) Client() *etcd.Client {
	return c.client
}

// Config returns the app configuration.
func (c *Controller) Config() *config.Config {
	return c.config
}

// Add appends a handler to the map.
func (c *Controller) Add(h Handler) {
	c.handlers[h.Command()] = h
}

// Handlers returns the complete list of added handlers.
func (c *Controller) Handlers() HandlerMap {
	return c.handlers
}

// WorkingDir returns the working directory. The value of a is appended to the value.
func (c *Controller) WorkingDir(a string) string {
	a = strings.Replace(a, "\n", "", -1)
	wdir := c.wdir + a
	return path.Clean(wdir)
}

// ChangeWorkingDir changes the current working directory.
func (c *Controller) ChangeWorkingDir(wdir string) string {
	if strings.HasPrefix(wdir, "/") {
		c.wdir = wdir
	} else {
		c.wdir = c.WorkingDir("/" + wdir)
	}

	return c.wdir
}

// ps1 returns the first type of prompt.
func (c *Controller) ps1() string {
	return fmt.Sprintf("%s@etcd:%s$ ", os.Getenv("USER"), c.wdir)
}

// ps2 returns the second type of prompt.
func (c *Controller) ps2() string {
	return "> "
}

// hasHandler returns whether a command handler has been added with the given id.
func (c *Controller) hasHandler(id string) bool {
	_, ok := c.handlers[id]
	return ok
}

// Handles the user input.
// Returns a boolean indicating whether the skip should exit. True to exit, false otherwise.
func (c *Controller) handleInput(i *eio.Input) bool {
	should_exit := false
	handler, ok := c.handlers[i.Cmd]
	if !ok {
		fmt.Fprintln(c.stderr, fmt.Sprintf("The command %s does not exist.", i.Cmd))
	} else if !handler.Validate(i) {
		fmt.Fprintln(c.stderr, fmt.Sprintf("Invalid use of command, use: %s", handler.Syntax()))
	} else {
		output, err := handler.Handle(i)
		if output == ExitSignal {
			should_exit = true
		} else if err == nil {
			fmt.Fprint(c.stdout, output)
		} else {
			fmt.Fprintln(c.stderr, err)
		}
	}

	return should_exit
}

// Welcome displays a welcome message.
func (c *Controller) welcome() {
	fmt.Fprintln(c.stdout, "Interactive etcd shell started.")
	if c.hasHandler("help") {
		fmt.Fprintln(c.stdout, "Type 'help' for a list of commands.")
	}
	if c.hasHandler("q") {
		fmt.Fprintln(c.stdout, "Type 'q' to quit.")
	}
	fmt.Fprintln(c.stdout, "")
}
