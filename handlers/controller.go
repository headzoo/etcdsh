package handlers

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/coreos/go-etcd/etcd"
	"github.com/headzoo/etcdsh/config"
	"github.com/headzoo/etcdsh/io"
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
	Handle(*io.Input) (string, error)
	Validate(*io.Input) bool
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

	for c.prompt() && c.scanner.Scan() {
		parts := strings.SplitN(c.scanner.Text(), " ", 3)
		if parts[0] != "" {
			should_exit := c.handleInput(io.NewFromArray(parts))
			if should_exit {
				break
			}
		}

	}
	if err := c.scanner.Err(); err != nil {
		fmt.Fprintln(c.stderr, "Reading stdin: %s", err)
		return 1
	}

	return 0
}

// Client returns the etcd client
func (c *Controller) Client() *etcd.Client {
	return c.client
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

func (c *Controller) prompt() bool {
	fmt.Fprintf(c.stdout, "%s@etcd:%s$ ", os.ExpandEnv("$USER"), c.wdir)
	return true
}

// hasHandler returns whether a command handler has been added with the given id.
func (c *Controller) hasHandler(id string) bool {
	_, ok := c.handlers[id]
	return ok
}

// Handles the user input.
// Returns a boolean indicating whether the skip should exit. True to exit, false otherwise.
func (c *Controller) handleInput(i *io.Input) bool {
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
