package utils

import (
	"fmt"
	"math/rand"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

const (
	randNodeGroupNameLength     = 8
	randNodeGroupNameComponents = "abcdef0123456789"
)

// NodeGroupName generates a random hex string with the fixed length of 8
func NodeGroupName() string {
	b := make([]byte, randNodeGroupNameLength)
	for i := 0; i < randNodeGroupNameLength; i++ {
		b[i] = randNodeGroupNameComponents[r.Intn(len(randNodeGroupNameComponents))]
	}
	return fmt.Sprintf("ng-%s", string(b))
}
