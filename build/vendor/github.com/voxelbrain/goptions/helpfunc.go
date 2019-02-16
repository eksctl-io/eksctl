package goptions

import (
	"io"
	"sync"
	"text/tabwriter"
	"text/template"
)

// HelpFunc is the signature of a function responsible for printing the help.
type HelpFunc func(w io.Writer, fs *FlagSet)

// Generates a new HelpFunc taking a `text/template.Template`-formatted
// string as an argument. The resulting template will be executed with the FlagSet
// as its data.
func NewTemplatedHelpFunc(tpl string) HelpFunc {
	var once sync.Once
	var t *template.Template
	return func(w io.Writer, fs *FlagSet) {
		once.Do(func() {
			t = template.Must(template.New("helpTemplate").Parse(tpl))
		})
		err := t.Execute(w, fs)
		if err != nil {
			panic(err)
		}
	}
}

const (
	_DEFAULT_HELP = "\xffUsage: {{.Name}} [global options] {{with .Verbs}}<verb> [verb options]{{end}}\n" +
		"\n" +
		"Global options:\xff" +
		"{{range .Flags}}" +
		"\n\t" +
		"\t{{with .Short}}" + "-{{.}}," + "{{end}}" +
		"\t{{with .Long}}" + "--{{.}}" + "{{end}}" +
		"\t{{.Description}}" +
		"{{with .DefaultValue}}" +
		" (default: {{.}})" +
		"{{end}}" +
		"{{if .Obligatory}}" +
		" (*)" +
		"{{end}}" +
		"{{end}}" +
		"\xff\n\n{{with .Verbs}}Verbs:\xff" +
		"{{range .}}" +
		"\xff\n    {{.Name}}:\xff" +
		"{{range .Flags}}" +
		"\n\t" +
		"\t{{with .Short}}" + "-{{.}}," + "{{end}}" +
		"\t{{with .Long}}" + "--{{.}}" + "{{end}}" +
		"\t{{.Description}}" +
		"{{with .DefaultValue}}" +
		" (default: {{.}})" +
		"{{end}}" +
		"{{if .Obligatory}}" +
		" (*)" +
		"{{end}}" +
		"{{end}}" +
		"{{end}}" +
		"{{end}}" +
		"\n"
)

// DefaultHelpFunc is a HelpFunc which renders the default help template and pipes
// the output through a text/tabwriter.Writer before flushing it to the output.
func DefaultHelpFunc(w io.Writer, fs *FlagSet) {
	tw := tabwriter.NewWriter(w, 4, 4, 1, ' ', tabwriter.StripEscape|tabwriter.DiscardEmptyColumns)
	NewTemplatedHelpFunc(_DEFAULT_HELP)(tw, fs)
	tw.Flush()
}
