/*
package goptions implements a flexible parser for command line options.

Key targets were the support for both long and short flag versions, mutually
exclusive flags, and verbs. Flags and their corresponding variables are defined
by the tags in a (possibly anonymous) struct.

    var options struct {
    	Name string `goptions:"-n, --name"`
    	Force bool `goptions:"-f, --force"`
    	Verbosity int `goptions:"-v, --verbose"`
    }

Short flags can be combined (e.g. `-nfv`). Long flags take their value after a
separating space. The equals notation (`--long-flag=value`) is NOT supported
right now.

Every member of the struct which is supposed to catch a command line value
has to have a "goptions" tag. The contains the short and long flag names for this
member but can additionally specify any of these options below.

    obligatory        - Flag must be specified. Otherwise an error will be returned
                        when Parse() is called.
    description='...' - Set the description for this particular flag. Will be
                        used by the HelpFunc.
    mutexgroup='...'  - Add this flag to a MutexGroup. Only one flag of the
                        ones sharing a MutexGroup can be set. Otherwise an error
                        will be returned when Parse() is called. If one flag in a
                        MutexGroup is `obligatory` one flag of the group must be
                        specified. A flag can be in multiple MutexGroups at once.

Depending on the type of the struct member, additional options might become available:

    Type: *os.File
        The given string is interpreted as a path to a file. If the string is "-"
        os.Stdin or os.Stdout will be used. os.Stdin will be returned, if the
        `rdonly` flag was set. os.Stdout will be returned, if `wronly` was set.
    Available options:
        Any combination of create, append, rdonly, wronly, rdwr,
        excl, sync, trunc and perm can be specified and correspond directly with
        the combination of the homonymous flags in the os package.

    Type: *net.TCPAddr
        The given string is interpreted as a tcp address. It is passed to
        net.ResolvTCPAddr() with "tcp" as the network type identifier.

    Type: *net/url.URL
        The given string is parsed by net/url.Parse()

    Type: time.Duration
        The given string is parsed by time.ParseDuration()

If a member is a slice type, multiple definitions of the flags are possible. For each
specification the underlying type will be used.

    var options struct {
        Servers []string `goptions:"-s, --server, description='Servers to connect to'"`
    }{}

goptions also has support for verbs. Each verb accepts its own set of flags which
take exactly the same tag format as global options. For an usage example of verbs
see the PrintHelp() example.
*/
package goptions

import (
	"os"
	"path/filepath"
)

const (
	VERSION = "2.5.11"
)

var (
	globalFlagSet *FlagSet
)

// ParseAndFail is a convenience function to parse os.Args[1:] and print
// the help if an error occurs. This should cover 90% of this library's
// applications.
func ParseAndFail(v interface{}) {
	globalFlagSet = NewFlagSet(filepath.Base(os.Args[0]), v)
	globalFlagSet.ParseAndFail(os.Stderr, os.Args[1:])
}

// Parse parses the command-line flags from os.Args[1:].
func Parse(v interface{}) error {
	globalFlagSet = NewFlagSet(filepath.Base(os.Args[0]), v)
	return globalFlagSet.Parse(os.Args[1:])
}

// PrintHelp renders the default help to os.Stderr.
func PrintHelp() {
	if globalFlagSet == nil {
		panic("Must call Parse() before PrintHelp()")
	}
	globalFlagSet.PrintHelp(os.Stderr)
}
