package random

import (
	"math/rand"
	"time"
)

const charBytes = "abcdefghijklmnopqrstuvwxyz0123456789"

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// String returns a random alphanumerical string of length n.
func String(n int) string {
	if n < 0 {
		return ""
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = charBytes[rand.Intn(len(charBytes))]
	}
	return string(b)
}
