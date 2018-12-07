package cmdutils

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Grouping holds a superset of all flagsets for all commands
type Grouping struct {
	groups map[*cobra.Command]*NamedFlagSetGroup
}

type namedFlagSet struct {
	name string
	fs   *pflag.FlagSet
}

// NamedFlagSetGroup holds a single group of flagsets
type NamedFlagSetGroup struct {
	list []namedFlagSet
}

// NewGrouping creates an instance of Grouping
func NewGrouping() *Grouping {
	return &Grouping{
		make(map[*cobra.Command]*NamedFlagSetGroup),
	}
}

// New creates a new group of flagsets for use with a subcommand
func (g *Grouping) New(cmd *cobra.Command) *NamedFlagSetGroup {
	n := &NamedFlagSetGroup{}
	g.groups[cmd] = n
	return n
}

// InFlagSet returs new or existing named GlagSet in a group
func (n *NamedFlagSetGroup) InFlagSet(name string, cb func(*pflag.FlagSet)) {
	for _, nfs := range n.list {
		if nfs.name == name {
			cb(nfs.fs)
			return
		}
	}

	nfs := namedFlagSet{
		name: name,
		fs:   &pflag.FlagSet{},
	}
	cb(nfs.fs)
	n.list = append(n.list, nfs)
}

// AddTo mixes all flagsets in the given group into another flagset
func (n *NamedFlagSetGroup) AddTo(cmd *cobra.Command) {
	for _, nfs := range n.list {
		cmd.Flags().AddFlagSet(nfs.fs)
	}
}

// Usage is for use with (*cobra.Command).SetUsageFunc
func (g *Grouping) Usage(cmd *cobra.Command) error {
	if cmd == nil {
		return fmt.Errorf("nil command")
	}

	group := g.groups[cmd]

	usage := []string{fmt.Sprintf("Usage: %s", cmd.UseLine())}

	if cmd.HasAvailableSubCommands() {
		usage = append(usage, "\nCommands:")
		for _, subCommand := range cmd.Commands() {
			usage = append(usage, fmt.Sprintf("  %s %-10s  %s", cmd.CommandPath(), subCommand.Name(), subCommand.Short))
		}
	}

	if len(cmd.Aliases) > 0 {
		usage = append(usage, "\nAliases: "+cmd.NameAndAliases())
	}

	if group != nil {
		for _, nfs := range group.list {
			usage = append(usage, fmt.Sprintf("\n%s flags:", nfs.name))
			usage = append(usage, strings.TrimRightFunc(nfs.fs.FlagUsages(), unicode.IsSpace))
		}
	}

	usage = append(usage, "\nCommon flags:")
	if len(cmd.PersistentFlags().FlagUsages()) != 0 {
		usage = append(usage, strings.TrimRightFunc(cmd.PersistentFlags().FlagUsages(), unicode.IsSpace))
	}
	if len(cmd.InheritedFlags().FlagUsages()) != 0 {
		usage = append(usage, strings.TrimRightFunc(cmd.InheritedFlags().FlagUsages(), unicode.IsSpace))
	}

	usage = append(usage, fmt.Sprintf("\nUse '%s [command] --help' for more information about a command.\n", cmd.CommandPath()))

	cmd.Println(strings.Join(usage, "\n"))

	return nil
}
