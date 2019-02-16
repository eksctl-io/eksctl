package goptions

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"
)

// A FlagSet represents one set of flags which belong to one particular program.
// A FlagSet is also used to represent a subset of flags belonging to one verb.
type FlagSet struct {
	// This HelpFunc will be called when PrintHelp() is called.
	HelpFunc
	// Name of the program. Might be used by HelpFunc.
	Name          string
	helpFlag      *Flag
	remainderFlag *Flag
	shortMap      map[string]*Flag
	longMap       map[string]*Flag
	verbFlag      *Flag
	// Global option flags
	Flags []*Flag
	// Verbs and corresponding FlagSets
	Verbs  map[string]*FlagSet
	parent *FlagSet
}

// NewFlagSet returns a new FlagSet containing all the flags which result from
// parsing the tags of the struct. Said struct as to be passed to the function
// as a pointer.
// If a tag line is erroneous, NewFlagSet() panics as this is considered a
// compile time error rather than a runtme error.
func NewFlagSet(name string, v interface{}) *FlagSet {
	structValue := reflect.ValueOf(v)
	if structValue.Kind() != reflect.Ptr {
		panic("Value type is not a pointer to a struct")
	}
	structValue = structValue.Elem()
	if structValue.Kind() != reflect.Struct {
		panic("Value type is not a pointer to a struct")
	}
	return newFlagset(name, structValue, nil)
}

// Internal version which skips type checking and takes the "parent"'s
// remainder flag as a parameter.
func newFlagset(name string, structValue reflect.Value, parent *FlagSet) *FlagSet {
	var once sync.Once
	r := &FlagSet{
		Name:     name,
		Flags:    make([]*Flag, 0),
		HelpFunc: DefaultHelpFunc,
		parent:   parent,
	}

	if parent != nil && parent.remainderFlag != nil {
		r.remainderFlag = parent.remainderFlag
	}

	var i int
	// Parse Option fields
	for i = 0; i < structValue.Type().NumField(); i++ {
		// Skip unexported fields
		if StartsWithLowercase(structValue.Type().Field(i).Name) {
			continue
		}

		fieldValue := structValue.Field(i)
		tag := structValue.Type().Field(i).Tag.Get("goptions")
		flag, err := parseStructField(fieldValue, tag)

		if err != nil {
			panic(fmt.Sprintf("Invalid struct field: %s", err))
		}
		if fieldValue.Type().Name() == "Verbs" {
			r.verbFlag = flag
			break
		}
		if fieldValue.Type().Name() == "Help" {
			r.helpFlag = flag
		}
		if fieldValue.Type().Name() == "Remainder" && r.remainderFlag == nil {
			r.remainderFlag = flag
		}

		if len(tag) != 0 {
			r.Flags = append(r.Flags, flag)
		}
	}

	// Parse verb fields
	for i++; i < structValue.Type().NumField(); i++ {
		once.Do(func() {
			r.Verbs = make(map[string]*FlagSet)
		})
		fieldValue := structValue.Field(i)
		tag := structValue.Type().Field(i).Tag.Get("goptions")
		r.Verbs[tag] = newFlagset(tag, fieldValue, r)
	}
	r.createMaps()
	return r
}

var (
	ErrHelpRequest = errors.New("Request for Help")
)

// Parse takes the command line arguments and sets the corresponding values
// in the FlagSet's struct.
func (fs *FlagSet) Parse(args []string) (err error) {
	// Parse global flags
	for len(args) > 0 {
		if !((isLong(args[0]) && fs.hasLongFlag(args[0][2:])) ||
			(isShort(args[0]) && fs.hasShortFlag(args[0][1:2]))) {
			break
		}
		f := fs.FlagByName(args[0])
		args, err = f.Parse(args)
		if err != nil {
			return
		}
		if f == fs.helpFlag && f.WasSpecified {
			return ErrHelpRequest
		}
	}

	// Process verb
	if len(args) > 0 {
		if verb, ok := fs.Verbs[args[0]]; ok {
			fs.verbFlag.value.Set(reflect.ValueOf(Verbs(args[0])))
			err := verb.Parse(args[1:])
			if err != nil {
				return err
			}
			args = args[0:0]
		}
	}

	// Process remainder
	if len(args) > 0 {
		if fs.remainderFlag == nil {
			return fmt.Errorf("Invalid trailing arguments: %v", args)
		}
		remainder := reflect.MakeSlice(fs.remainderFlag.value.Type(), len(args), len(args))
		reflect.Copy(remainder, reflect.ValueOf(args))
		fs.remainderFlag.value.Set(remainder)
	}

	// Check for unset, obligatory, single Flags
	for _, f := range fs.Flags {
		if f.Obligatory && !f.WasSpecified && len(f.MutexGroups) == 0 {
			return fmt.Errorf("%s must be specified", f.Name())
		}
	}

	// Check for multiple set Flags in one mutex group
	// Check also for unset, obligatory mutex groups
	mgs := fs.MutexGroups()
	for _, mg := range mgs {
		if !mg.IsValid() {
			return fmt.Errorf("Exactly one of %s must be specified", strings.Join(mg.Names(), ", "))
		}
	}
	return nil
}

func (fs *FlagSet) createMaps() {
	fs.longMap = make(map[string]*Flag)
	fs.shortMap = make(map[string]*Flag)
	for _, flag := range fs.Flags {
		fs.longMap[flag.Long] = flag
		fs.shortMap[flag.Short] = flag
	}
}

func (fs *FlagSet) hasLongFlag(fname string) bool {
	_, ok := fs.longMap[fname]
	return ok
}

func (fs *FlagSet) hasShortFlag(fname string) bool {
	_, ok := fs.shortMap[fname]
	return ok
}

func (fs *FlagSet) FlagByName(fname string) *Flag {
	if isShort(fname) && fs.hasShortFlag(fname[1:2]) {
		return fs.shortMap[fname[1:2]]
	} else if isLong(fname) && fs.hasLongFlag(fname[2:]) {
		return fs.longMap[fname[2:]]
	}
	return nil
}

// MutexGroups returns a map of Flag lists which contain mutually
// exclusive flags.
func (fs *FlagSet) MutexGroups() map[string]MutexGroup {
	r := make(map[string]MutexGroup)
	for _, f := range fs.Flags {
		for _, mg := range f.MutexGroups {
			if len(mg) == 0 {
				continue
			}
			if _, ok := r[mg]; !ok {
				r[mg] = make(MutexGroup, 0)
			}
			r[mg] = append(r[mg], f)
		}
	}
	return r
}

// Prints the FlagSet's help to the given writer.
func (fs *FlagSet) PrintHelp(w io.Writer) {
	fs.HelpFunc(w, fs)
}

func (fs *FlagSet) ParseAndFail(w io.Writer, args []string) {
	err := fs.Parse(args)
	if err != nil {
		errCode := 0
		if err != ErrHelpRequest {
			errCode = 1
			fmt.Fprintf(w, "Error: %s\n", err)
		}
		fs.PrintHelp(w)
		os.Exit(errCode)
	}
}

func StartsWithLowercase(s string) bool {
	if len(s) <= 0 {
		return false
	}
	return strings.ToLower(s)[0] == s[0]
}
