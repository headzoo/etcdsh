package handlers

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

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

// The color codes to use when outputing.
type OutputColors struct {
	Key    string
	Object string
}

// LsHandler handles the "ls" command.
type LsHandler struct {
	controller *Controller
	colors     OutputColors
	use_colors bool
}

// NewLsHandler creates a new LsHandler instance.
func NewLsHandler(controller *Controller) *LsHandler {
	h := new(LsHandler)
	h.controller = controller
	h.setupColors()

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

	//return h.respToShortOutput(resp), nil
	return h.respToLongOutput(resp), nil
}

// respToLongOuput formats an etcd response for output in the long format.
func (h *LsHandler) respToLongOutput(resp *etcd.Response) string {
	output := bytes.NewBufferString("")
	widths := columnWidths(resp.Node)
	node := etcd.Node{
		Dir:           true,
		Key:           ".",
		CreatedIndex:  0,
		ModifiedIndex: 0,
	}
	output.WriteString(h.formatNode(&node, widths))
	node.Key = ".."
	output.WriteString(h.formatNode(&node, widths))

	total := 2
	for _, node := range resp.Node.Nodes {
		output.WriteString(h.formatNode(node, widths))
		total++
		for _, n := range node.Nodes {
			output.WriteString(h.formatNode(n, widths))
			total++
		}
	}

	return fmt.Sprintf("total %d\n%s", total, output.String())
}

// respToShortOutput formats an etcd response for output in the short format.
func (h *LsHandler) respToShortOutput(resp *etcd.Response) string {
	output := bytes.NewBufferString("")
	for _, node := range resp.Node.Nodes {
		output.WriteString(path.Base(node.Key))
		output.WriteString(" ")

		for _, n := range node.Nodes {
			output.WriteString(path.Base(n.Key))
			output.WriteString(" ")
		}
	}

	output.WriteString("\n")
	return output.String()
}

// formatNode formats the node as a string for output to the console.
func (h *LsHandler) formatNode(n *etcd.Node, w ColumnWidths) string {
	typeValue := typeKey
	if !n.Dir {
		typeValue = typeObject
	}

	prefix := ""
	postfix := ""
	if h.use_colors {
		if n.Dir {
			prefix = "\x1b[" + h.colors.Key + ";1m"
		} else {
			prefix = "\x1b[" + h.colors.Object + ";1m"
		}
		postfix = "\x1b[0m"
	}

	return fmt.Sprintf(
		"%*d %*d %*d %s: %s%s%s\n",
		w.CreatedIndex,
		n.CreatedIndex,
		w.ModifiedIndex,
		n.ModifiedIndex,
		w.TTL,
		n.TTL,
		typeValue,
		prefix,
		path.Base(n.Key),
		postfix,
	)
}

// setupColors sets the value of LsHandler.colors.
func (h *LsHandler) setupColors() {
	h.colors = OutputColors{}
	h.use_colors = false

	if h.controller.Config().Colors && runtime.GOOS == "linux" {
		h.colors = OutputColors{
			Key:    "34",
			Object: "0",
		}
		h.use_colors = true

		ls_colors := os.Getenv("LS_COLORS")
		if ls_colors != "" {
			colors := strings.Split(ls_colors, ":")
			for _, color := range colors {
				if strings.HasPrefix(color, "di=") {
					p := strings.Split(color, "=")
					if len(p) > 1 {
						p = strings.Split(p[1], ";")
						if len(p) > 1 {
							h.colors.Key = p[1]
						}
					}
				} else if strings.HasPrefix(color, "fi=") {
					p := strings.Split(color, "=")
					if len(p) > 1 {
						p = strings.Split(p[1], ";")
						if len(p) > 1 {
							h.colors.Object = p[1]
						}
					}
				}
			}
		}
	}
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
