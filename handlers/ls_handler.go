package handlers

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"strconv"

	"github.com/coreos/go-etcd/etcd"
	"github.com/headzoo/etcdsh/env"
	"flag"
)

const (
	// Represents a node is a directory in the output.
	SymbolTypeKeys = "k"

	// Represents a node is a file in the output.
	SymbolTypeObjects = "o"

	// Default color for keys.
	DefaultColorKeys = "34"

	// Default color for objects.
	DefaultColorObjects = "0"
)

// Column widths to use for the "ls" output.
type ColumnWidths struct {
	CreatedIndex  int
	ModifiedIndex int
	TTL           int
}

// The color codes to use when outputting.
type OutputColors struct {
	Key    string
	Object string
}

// LsHandler handles the "ls" command.
type LsHandler struct {
	CommandHandler
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
func (h *LsHandler) Validate(i *Input) bool {
	return true
}

// Description returns a string that describes the command.
func (h *LsHandler) Description() string {
	return "Displays a listing of the current working directory"
}

// Handles the "ls" command.
func (h *LsHandler) Handle(i *Input) (string, error) {
	flags := flag.NewFlagSet("ls_flagset", flag.ContinueOnError)
	help := flags.Bool("h", false, "Show command help")
	sort := flags.Bool("s", false, "Sort the results")
	long := flags.Bool("l", false, "Use long list format")
	if flags.Parse(i.Args) != nil {
		return "", nil
	}
	if *help || len(flags.Args()) == 0 {
		printCommandHelp(h, flags)
		return "", nil
	}
	
	args := flags.Args()
	dir := h.controller.WorkingDir(args[0])
	resp, err := h.controller.Client().Get(dir, *sort, false)
	if err != nil {
		return "", err
	}

	if *long {
		return h.respToLongOutput(resp), nil
	}
	return h.respToShortOutput(resp), nil
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
	typeValue := SymbolTypeKeys
	if !n.Dir {
		typeValue = SymbolTypeObjects
	}

	prefix := ""
	postfix := ""
	if h.use_colors {
		if n.Dir {
			prefix = "\x1b["+h.colors.Key+";1m"
		} else {
			prefix = "\x1b["+h.colors.Object+";1m"
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
		envColors := env.NewColors()
		di, _ := envColors.GetLSDefault("di", DefaultColorKeys)
		fi, _ := envColors.GetLSDefault("fi", DefaultColorObjects)
		h.colors = OutputColors{
			Key:    di,
			Object: fi,
		}
		h.use_colors = true
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

