package handlers

import (
	"fmt"
	"path"
	"strconv"

	"github.com/coreos/go-etcd/etcd"
	"github.com/headzoo/etcdsh/io"
)

const (
	// Represents a node is a directory in the output.
	typeKey = "k"

	// Represents a node is a file in the output.
	typeObject = "o"
)

// Column widths to use for the "ls" output.
type ColumnWidths struct {
	CreatedIndex  int
	ModifiedIndex int
	TTL           int
}

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

	widths := columnWidths(resp.Node)
	node := etcd.Node{
		Dir:           true,
		Key:           ".",
		CreatedIndex:  0,
		ModifiedIndex: 0,
	}
	output := formatNode(&node, widths)
	node.Key = ".."
	output += formatNode(&node, widths)

	total := 2
	for _, node := range resp.Node.Nodes {
		output += formatNode(node, widths)
		total++
		for _, n := range node.Nodes {
			output += formatNode(n, widths)
			total++
		}
	}

	return fmt.Sprintf("total %d\n%s", total, output), nil
}

// formatNode formats the node as a string for output to the console.
func formatNode(n *etcd.Node, w ColumnWidths) string {
	typeValue := typeKey
	if !n.Dir {
		typeValue = typeObject
	}

	return fmt.Sprintf(
		"%*d %*d %*d %s: %s\n",
		w.CreatedIndex,
		n.CreatedIndex,
		w.ModifiedIndex,
		n.ModifiedIndex,
		w.TTL,
		n.TTL,
		typeValue,
		path.Base(n.Key),
	)
}

// columnWidths returns the widths for each column in the "ls" output.
func columnWidths(resp_node *etcd.Node) ColumnWidths {
	widths := ColumnWidths{
		CreatedIndex:  len(strconv.FormatUint(resp_node.CreatedIndex, 10)),
		ModifiedIndex: len(strconv.FormatUint(resp_node.ModifiedIndex, 10)),
		TTL:           len(strconv.FormatInt(resp_node.TTL, 10)),
	}
	cw := 0

	for _, node := range resp_node.Nodes {
		cw = len(strconv.FormatUint(node.CreatedIndex, 10))
		if cw > widths.CreatedIndex {
			widths.CreatedIndex = cw
		}
		cw = len(strconv.FormatUint(node.ModifiedIndex, 10))
		if cw > widths.ModifiedIndex {
			widths.ModifiedIndex = cw
		}
		cw = len(strconv.FormatInt(node.TTL, 10))
		if cw > widths.TTL {
			widths.TTL = cw
		}
		for _, n := range node.Nodes {
			cw = len(strconv.FormatUint(n.CreatedIndex, 10))
			if cw > widths.CreatedIndex {
				widths.CreatedIndex = cw
			}
			cw = len(strconv.FormatUint(n.ModifiedIndex, 10))
			if cw > widths.ModifiedIndex {
				widths.ModifiedIndex = cw
			}
			cw = len(strconv.FormatInt(n.TTL, 10))
			if cw > widths.TTL {
				widths.TTL = cw
			}
		}
	}

	return widths
}
