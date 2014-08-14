package handlers

import (
	"bytes"
	"flag"
	"fmt"
	"path"
	"runtime"
	"strconv"

	"github.com/coreos/go-etcd/etcd"
	"github.com/headzoo/etcdsh/env"
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
type LsColumnWidths struct {
	CreatedIndex  int
	ModifiedIndex int
	TTL           int
	Keys          int
}

// The color codes to use when outputting.
type LsOutputColors struct {
	Key    string
	Object string
}

// Command line options for the ls command.
type LsOptions struct {
	PrintHelp  bool
	LongFormat bool
	Sorted     bool
}

// LsHandler handles the "ls" command.
type LsHandler struct {
	controller *Controller
	colors     LsOutputColors
	use_colors bool
}

// NewLsHandler creates a new LsHandler instance.
func NewLsHandler(controller *Controller) *LsHandler {
	h := &LsHandler{
		controller: controller,
	}
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
	opts, args, err := h.setupOptions(i.Args)
	if opts == nil || err != nil {
		return "", err
	}
	
	dir := h.controller.WorkingDir(args[0])
	resp, err := h.controller.Client().Get(dir, opts.Sorted, false)
	if err != nil {
		return "", err
	}

	if opts.LongFormat {
		return h.respToLongOutput(resp), nil
	}
	return h.respToShortOutput(resp), nil
}

// setupOptions builds a FlagSet and parses the args passed to the command.
func (h *LsHandler) setupOptions(args []string) (*LsOptions, []string, error) {
	opts := &LsOptions{}
	flags := flag.NewFlagSet("ls_flags", flag.ContinueOnError)
	flags.BoolVar(&opts.PrintHelp, "h", false, "Show the command help")
	flags.BoolVar(&opts.LongFormat, "l", false, "Use long list format")
	flags.BoolVar(&opts.Sorted, "s", false, "Sort the results")

	err := flags.Parse(args)
	if err != nil {
		return nil, nil, err
	}
	if opts.PrintHelp {
		printCommandHelp(h, flags)
		return nil, nil, nil
	}

	args = flags.Args()
	if len(args) == 0 {
		args = []string{"/"}
	}

	return opts, args, nil
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
	output.WriteString(h.formatNodeLong(&node, widths))
	node.Key = ".."
	output.WriteString(h.formatNodeLong(&node, widths))

	total := 2
	for _, node := range resp.Node.Nodes {
		output.WriteString(h.formatNodeLong(node, widths))
		total++
		for _, n := range node.Nodes {
			output.WriteString(h.formatNodeLong(n, widths))
			total++
		}
	}

	return fmt.Sprintf("total %d\n%s", total, output.String())
}

// respToShortOutput formats an etcd response for output in the short format.
func (h *LsHandler) respToShortOutput(resp *etcd.Response) string {
	output := bytes.NewBufferString("")
	widths := columnWidths(resp.Node)
	for _, node := range resp.Node.Nodes {
		output.WriteString(h.formatNodeShort(node, widths))
		//output.WriteString(" ")

		for _, n := range node.Nodes {
			output.WriteString(h.formatNodeShort(n, widths))
			//output.WriteString(" ")
		}
	}

	output.WriteString("\n")
	return output.String()
}

func (h *LsHandler) formatNodeShort(n *etcd.Node, w LsColumnWidths) string {
	prefix, postfix := "", ""
	if h.use_colors {
		if n.Dir {
			prefix = env.ColorPrefixCode(h.colors.Key)
		} else {
			prefix = env.ColorPrefixCode(h.colors.Object)
		}
		postfix = env.ColorPostfixCode()
	}

	return fmt.Sprintf(
		"%s%-*s%s",
		prefix,
		w.Keys,
		path.Base(n.Key),
		postfix,
	)
}

// formatNodeLong formats the node as a string for output to the console.
func (h *LsHandler) formatNodeLong(n *etcd.Node, w LsColumnWidths) string {
	typeValue := SymbolTypeKeys
	if !n.Dir {
		typeValue = SymbolTypeObjects
	}

	prefix := ""
	postfix := ""
	if h.use_colors {
		if n.Dir {
			prefix = env.ColorPrefixCode(h.colors.Key)
		} else {
			fmt.Println(h.colors.Object)
			prefix = env.ColorPrefixCode(h.colors.Object)
		}
		postfix = env.ColorPostfixCode()
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
	h.colors = LsOutputColors{}
	h.use_colors = false

	if h.controller.Config().Colors && runtime.GOOS == "linux" {
		envColors := env.NewColors()
		di, _ := envColors.GetLSDefault("di", DefaultColorKeys)
		fi, _ := envColors.GetLSDefault("fi", DefaultColorObjects)
		h.colors = LsOutputColors{
			Key:    di,
			Object: fi,
		}
		h.use_colors = true
	}
}

// columnWidths returns the widths for each column in the "ls" output.
func columnWidths(resp_node *etcd.Node) LsColumnWidths {
	widths := LsColumnWidths{
		CreatedIndex:  len(strconv.FormatUint(resp_node.CreatedIndex, 10)),
		ModifiedIndex: len(strconv.FormatUint(resp_node.ModifiedIndex, 10)),
		TTL:           len(strconv.FormatInt(resp_node.TTL, 10)),
		Keys:          len(resp_node.Key),
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
		cw = len(node.Key)
		if cw > widths.Keys {
			widths.Keys = cw
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
			cw = len(n.Key)
			if cw > widths.Keys {
				widths.Keys = cw
			}
		}
	}

	return widths
}
