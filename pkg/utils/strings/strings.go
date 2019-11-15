package strings

// Pointer returns a pointer to the provided string.
func Pointer(s string) *string {
	return &s
}

// NilIfEmpty returns nil if the provided string is empty, or else a pointer to
// the provided string.
func NilIfEmpty(s string) *string {
	if len(s) == 0 {
		return nil
	}
	return &s
}

// EmptyIfNil returns an empty string if the provided *string is nil, or else
// the value the provided pointer points to.
func EmptyIfNil(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// ToPointersMap converts the provided map[string]string into a
// map[string]*string.
func ToPointersMap(in map[string]string) map[string]*string {
	out := make(map[string]*string, len(in))
	for k := range in {
		v := in[k]
		out[k] = &v
	}
	return out
}

// NilPointersMapIfEmpty returns nil if the provided map[string]*string is
// empty, or the provided map otherwise.
func NilPointersMapIfEmpty(in map[string]*string) map[string]*string {
	if len(in) == 0 {
		return nil
	}
	return in
}

// ToValuesMap converts the provided map[string]*string into a
// map[string]string.
func ToValuesMap(in map[string]*string) map[string]string {
	out := make(map[string]string, len(in))
	for k := range in {
		v := in[k]
		out[k] = *v
	}
	return out
}

// ToPointersArray converts the provided []string into a []*string.
func ToPointersArray(in []string) []*string {
	out := make([]*string, len(in))
	for i := range in {
		out[i] = &in[i]
	}
	return out
}

// NilPointersArrayIfEmpty returns nil if the provided []*string array is
// empty, or the provided array otherwise.
func NilPointersArrayIfEmpty(in []*string) []*string {
	if len(in) == 0 {
		return nil
	}
	return in
}

// ToValuesArray converts the provided []*string into a []string.
func ToValuesArray(in []*string) []string {
	out := make([]string, len(in))
	for i := range in {
		out[i] = *in[i]
	}
	return out
}
