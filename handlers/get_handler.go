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

import "fmt"

// GetHandler handles the "exit" command.
type GetHandler struct {
	controller *Controller
}

// NewGetHandler returns a new ExitHandler instance.
func NewGetHandler(controller *Controller) *GetHandler {
	return &GetHandler{
		controller: controller,
	}
}

// Command returns the string typed by the user that triggers to handler.
func (h *GetHandler) Command() string {
	return "get"
}

// Validate returns whether the user input is valid.
func (h *GetHandler) Validate(i *Input) bool {
	return len(i.Args) > 0
}

// Syntax returns a string that demonstrates how to use the command.
func (h *GetHandler) Syntax() string {
	return "get <key>"
}

// Description returns a string that describes the command.
func (h *GetHandler) Description() string {
	return "Displays the value of the given key"
}

// Handles the "get" command.
func (h *GetHandler) Handle(i *Input) (string, error) {
	dir := h.controller.WorkingDir(i.Args[0])
	resp, err := h.controller.Client().Get(dir, false, false)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n", resp.Node.Value), nil
}
