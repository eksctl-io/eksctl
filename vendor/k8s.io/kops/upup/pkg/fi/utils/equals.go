/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

func StringSlicesEqual(l, r []string) bool {
	if len(l) != len(r) {
		return false
	}
	for i, v := range l {
		if r[i] != v {
			return false
		}
	}
	return true
}

func StringSlicesEqualIgnoreOrder(l, r []string) bool {
	if len(l) != len(r) {
		return false
	}

	lMap := map[string]bool{}
	for _, lv := range l {
		lMap[lv] = true
	}
	for _, rv := range r {
		if !lMap[rv] {
			return false
		}
	}
	return true
}
