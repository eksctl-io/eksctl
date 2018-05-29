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

import (
	"bytes"
	"os"
	"strings"
)

// SanitizeString iterated a strings and removes any characters not in the allow list
func SanitizeString(s string) string {
	var out bytes.Buffer
	allowed := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"
	for _, c := range s {
		if strings.IndexRune(allowed, c) != -1 {
			out.WriteRune(c)
		} else {
			out.WriteRune('_')
		}
	}

	return string(out.Bytes())
}

// ExpandPath replaces common path aliases: ~ -> $HOME
func ExpandPath(p string) string {
	if strings.HasPrefix(p, "~/") {
		p = os.Getenv("HOME") + p[1:]
	}

	return p
}
