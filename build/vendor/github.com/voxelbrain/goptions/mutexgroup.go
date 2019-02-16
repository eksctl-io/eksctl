package goptions

// A MutexGroup holds a set of flags which are mutually exclusive and cannot
// be specified at the same time.
type MutexGroup []*Flag

// IsObligatory returns true if exactly one of the flags in the MutexGroup has
// to be specified
func (mg MutexGroup) IsObligatory() bool {
	for _, flag := range mg {
		if flag.Obligatory {
			return true
		}
	}
	return false
}

func (mg MutexGroup) WasSpecified() bool {
	for _, flag := range mg {
		if flag.WasSpecified {
			return true
		}
	}
	return false
}

// IsValid checks if the flags in the MutexGroup describe a valid state.
// I.e. At most one has been specified or – if it is an obligatory MutexGroup –
// exactly one has been specified.
func (mg MutexGroup) IsValid() bool {
	c := 0
	for _, flag := range mg {
		if flag.WasSpecified {
			c++
		}
	}
	return c <= 1 && (!mg.IsObligatory() || c == 1)
}

// Names is a convenience function to return the array of names of the flags
// in the MutexGroup.
func (mg MutexGroup) Names() []string {
	r := make([]string, len(mg))
	for i, flag := range mg {
		r[i] = flag.Name()
	}
	return r
}
