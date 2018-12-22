package utils

import (
	"math/rand"
	"time"
	"fmt"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

const (
	randNodegroupNameLength     = 8
	randNodeGroupNameComponents = "abcdef0123456789"
)

// NodegroupName generates a random hex string with the fixed length of 8
func NodegroupName() string {
	b := make([]byte, randNodegroupNameLength)
	for i := 0; i < randNodegroupNameLength; i++ {
		b[i] = randNodeGroupNameComponents[r.Intn(len(randNodeGroupNameComponents))]
	}
	return fmt.Sprintf("ng-%s", string(b))
}
