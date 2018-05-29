// Copyright © 2017 The Kubicorn Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rand

import (
	"crypto/rand"
	"math/big"
)

// MustGenerateRandomBytes generates random bytes or panics if it can't
func MustGenerateRandomBytes(length int) []byte {
	res, err := GenerateRandomBytes(length)

	if err != nil {
		panic("Could not generate random bytes")
	}

	return res
}

// GenerateRandomBytes ...
func GenerateRandomBytes(length int) ([]byte, error) {
	b := make([]byte, length)

	_, err := rand.Read(b)

	return b, err
}

// Generate random number in [0,n)
func GenerateRandomInt(min, max int) int {
	return int(GenerateRandomInt64(int64(min), int64(max)))
}

// Generate random number in [0,n)
func GenerateRandomInt64(min, max int64) int64 {
	upper := max - min

	nBig, err := rand.Int(rand.Reader, big.NewInt(upper))

	if err != nil {
		panic(err)
	}

	return nBig.Int64() + min
}
