package goptions

// Help Defines the common help flag. It is handled separately as it will cause
// Parse() to return ErrHelpRequest.
type Help bool

// Verbs marks the point in the struct where the verbs start. Its value will be
// the name of the selected verb.
type Verbs string

// A remainder catches all excessive arguments. If both a verb and
// the containing options struct have a remainder field, only the latter one
// will be used.
type Remainder []string
