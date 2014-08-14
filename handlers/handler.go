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
	"flag"
	"fmt"
)

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

// printCommandHelp is used by handlers to display command help.
func printCommandHelp(handler Handler, flags *flag.FlagSet) {
	fmt.Println("SYNTAX")
	fmt.Println("\t" + handler.Syntax())
	fmt.Println("")
	fmt.Println("OPTIONS:")
	flags.VisitAll(func(f *flag.Flag) {
		fmt.Printf("\t-%-10s%s\n", f.Name, f.Usage)
	})
}
