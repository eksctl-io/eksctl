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

package namer

import (
	"fmt"

	"github.com/kubicorn/kubicorn/pkg/rand"
)

func RandomName() string {
	adjl := len(adjectives)
	nounl := len(nouns)
	adj := adjectives[rand.GenerateRandomInt(0, adjl)]
	noun := nouns[rand.GenerateRandomInt(0, nounl)]

	return fmt.Sprintf("%s-%s", adj, noun)
}

var adjectives = []string{
	"unique",
	"beautiful",
	"attractive",
	"wonderful",
	"fabulous",
	"extravagant",
	"ridiculous",
	"exciting",
	"amazing",
	"adorable",
	"ferocious",
	"hilarious",
	"scrumptious",
	"floral",
}

var nouns = []string{
	"unicorn",
	"rainbow",
	"painting",
	"badger",
	"mongoose",
	"sculpture",
	"creature",
	"gopher",
	"wardrobe",
	"mushroom",
	"hideout",
	"party",
	"monster",
	"sheepdog",
	"outfit",
}
