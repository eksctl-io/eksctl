package slice

// Contains checks if a slice contains an element
func Contains(list []string, e string) bool {
	for _, x := range list {
		if x == e {
			return true
		}
	}

	return false
}
