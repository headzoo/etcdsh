/**
The MIT License (MIT)

Copyright (c) 2014 Sean Hickey <sean@dulotech.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/bobappleyard/readline"
	"github.com/coreos/go-etcd/etcd"
	"github.com/flynn/go-shlex"
	"github.com/headzoo/etcdsh/config"
	"github.com/headzoo/etcdsh/etcdsh"
	"github.com/headzoo/etcdsh/parser"
)

// Controller stores handlers and calls them.
type Controller struct {
	wdir                  string
	wdirKeys              []string
	handlers              HandlerMap
	config                *config.Config
	client                *etcd.Client
	stdout, stderr, stdin *os.File
	prompter              *parser.Prompt
}

// Create a new Controller.
func NewController(conf *config.Config, client *etcd.Client, stdout, stderr, stdin *os.File) *Controller {
	c := &Controller{
		config:   conf,
		client:   client,
		stdout:   stdout,
		stdin:    stdin,
		stderr:   stderr,
		wdir:     "/",
		handlers: make(HandlerMap),
		prompter: parser.NewPrompt(),
	}

	c.prompter.AddFormatter('w', func() string {
		return c.wdir
	})
	c.prompter.AddFormatter('W', func() string {
		return path.Base(c.wdir)
	})
	c.prompter.AddFormatter('v', func() string {
		return etcdsh.Version
	})
	c.prompter.AddFormatter('m', func() string {
		u, err := url.Parse(conf.Machine)
		if err == nil {
			return u.Host
		} else {
			return conf.Machine
		}
	})

	return c
}

// Starts the controller.
func (c *Controller) Start() int {
	c.welcome()
	c.ChangeWorkingDir("/")

	readline.Completer = c.filenameCompleter
	buffer := bytes.NewBufferString("")
	prompt := ""

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

		line = strings.TrimSpace(line)
		if strings.ToLower(line) == "q" || strings.ToLower(line) == "exit" {
			return 0
		}
		if strings.HasSuffix(line, "\\") {
			buffer.WriteString(strings.TrimSuffix(line, "\\") + "\n")
		} else {
			if buffer.Len() > 0 {
				buffer.WriteString(line)
				line = buffer.String()
				buffer.Reset()
			}

			parts, err := shlex.Split(line)
			if err != nil {
				panic(err)
			}
			readline.AddHistory(line)

			in := NewInput(parts[0])
			if len(parts) > 1 {
				in.Args = parts[1:]
			}
			c.handleInput(in)
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
	wdir := c.wdir + "/" + a
	return path.Clean(wdir)
}

// ChangeWorkingDir changes the current working directory.
func (c *Controller) ChangeWorkingDir(wdir string) string {
	if strings.HasPrefix(wdir, "/") {
		c.wdir = wdir
	} else {
		c.wdir = c.WorkingDir("/" + wdir)
	}

	resp, err := c.client.Get(c.wdir, true, true)
	if err != nil {
		panic(err)
	}

	count := c.getNodeCount(resp.Node, 0)
	c.wdirKeys = make([]string, count)
	c.addNodeToWorkingDirKeys(resp.Node, 0)

	return c.wdir
}

// addNodeToWDir adds the keys from all child nodes to the working dir keys.
func (c *Controller) addNodeToWorkingDirKeys(node *etcd.Node, index int) int {
	for _, n := range node.Nodes {
		c.wdirKeys[index] = n.Key
		index++
		index = c.addNodeToWorkingDirKeys(n, index)
	}

	return index
}

// getNodeCount recursively counts the child nodes in node.
func (c *Controller) getNodeCount(node *etcd.Node, count int) int {
	for _, n := range node.Nodes {
		count++
		count = c.getNodeCount(n, count)
	}

	return count
}

// ps1 returns the first type of prompt.
func (c *Controller) ps1() string {
	prompt, _ := c.prompter.Parse(c.config.PS1)
	return prompt
}

// ps2 returns the second type of prompt.
func (c *Controller) ps2() string {
	prompt, _ := c.prompter.Parse(c.config.PS2)
	return prompt
}

// hasHandler returns whether a command handler has been added with the given id.
func (c *Controller) hasHandler(id string) bool {
	_, ok := c.handlers[id]
	return ok
}

// Handles the user input.
func (c *Controller) handleInput(i *Input) {
	handler, ok := c.handlers[i.Cmd]
	if !ok {
		fmt.Fprintln(c.stderr, fmt.Sprintf("The command %s does not exist.", i.Cmd))
	} else if !handler.Validate(i) {
		fmt.Fprintln(c.stderr, fmt.Sprintf("Invalid use of command, use: %s", handler.Syntax()))
	} else {
		output, err := handler.Handle(i)
		if err == nil {
			fmt.Fprint(c.stdout, output)
		} else {
			fmt.Fprintln(c.stderr, err)
		}
	}
}

// filenameCompleter is a callback function for the readline.Completer variable.
func (c *Controller) filenameCompleter(query, ctx string) []string {
	var keys []string
	for _, key := range c.wdirKeys {
		base := strings.TrimPrefix(key, "/")
		if strings.HasPrefix(base, query) {
			keys = append(keys, base)
		}
	}

	return keys
}

// Welcome displays a welcome message.
func (c *Controller) welcome() {
	fmt.Fprintln(c.stdout, "Interactive etcd shell started.")
	if c.hasHandler("help") {
		fmt.Fprintln(c.stdout, "Type 'help' for a list of commands.")
	}
	fmt.Fprintln(c.stdout, "Type 'exit' or 'q' to quit.")
	fmt.Fprintln(c.stdout, "")
}
