package handlers

import (
	"fmt"
	"path"

	"github.com/coreos/go-etcd/etcd"
	"github.com/headzoo/etcdsh/io"
)

const (
	// Represents a node is a directory in the output.
	typeKey = "k"

	// Represents a node is a file in the output.
	typeObject = "o"
)

// LsHandler handles the "ls" command.
type LsHandler struct {
	controller *Controller
}

// NewLsHandler creates a new LsHandler instance.
func NewLsHandler(controller *Controller) *LsHandler {
	h := new(LsHandler)
	h.controller = controller

	return h
}

// Command returns the string typed by the user that triggers to handler.
func (h *LsHandler) Command() string {
	return "ls"
}

// Syntax returns a string that demonstrates how to use the command.
func (h *LsHandler) Syntax() string {
	return "ls <path>"
}

// Validate returns whether the user input is valid for this handler.
func (h *LsHandler) Validate(i *io.Input) bool {
	return true
}

// Description returns a string that describes the command.
func (h *LsHandler) Description() string {
	return "Displays a listing of the current working directory"
}

// Handles the "ls" command.
func (h *LsHandler) Handle(i *io.Input) (string, error) {
	dir := h.controller.WorkingDir(i.Key)
	resp, err := h.controller.Client().Get(dir, false, false)
	if err != nil {
		return "", err
	}

	total := 2
	node := etcd.Node{
		Dir:           true,
		Key:           ".",
		CreatedIndex:  0,
		ModifiedIndex: 0,
	}
	output := formatNode(&node)
	node.Key = ".."
	output += formatNode(&node)

	for _, node := range resp.Node.Nodes {
		output += formatNode(node)
		total++
		for _, n := range node.Nodes {
			output += formatNode(n)
			total++
		}
	}

	output = fmt.Sprintf("total %d\n%s", total, output)
	return output, nil
}

// formatNode formats the node as a string for output to the console.
func formatNode(n *etcd.Node) string {
	t := typeKey
	if !n.Dir {
		t = typeObject
	}
	return fmt.Sprintf("%-3d%-3d%-3d%s %s\n", n.CreatedIndex, n.ModifiedIndex, n.TTL, t, path.Base(n.Key))
}
