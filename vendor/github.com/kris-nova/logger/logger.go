// Copyright © 2017
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

package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	lol "github.com/kris-nova/lolgopher"
)

type Logger func(format string, a ...interface{})

const (
	AlwaysLabel   = "✿"
	CriticalLabel = "✖"
	DebugLabel    = "▶"
	InfoLabel     = "ℹ"
	SuccessLabel  = "✔"
	WarningLabel  = "!"
)

var (
	Level              = 2
	Color              = true
	Fabulous           = false
	FabulousWriter     = lol.NewLolWriter()
	FabulousTrueWriter = lol.NewTruecolorLolWriter()
	TestMode           = false
	Timestamps         = true
)

func Always(format string, a ...interface{}) {
	a, w := extractLoggerArgs(format, a...)
	s := fmt.Sprintf(label(format, AlwaysLabel), a...)

	if !TestMode {
		if Color {
			w = color.Output
			s = color.GreenString(s)
		} else if Fabulous {
			w = FabulousWriter
		}
	}

	fmt.Fprintf(w, s)
}

func Critical(format string, a ...interface{}) {
	if Level >= 1 {
		a, w := extractLoggerArgs(format, a...)
		s := fmt.Sprintf(label(format, CriticalLabel), a...)

		if !TestMode {
			if Color {
				w = color.Output
				s = color.RedString(s)
			} else if Fabulous {
				w = FabulousWriter
			}
		}

		fmt.Fprintf(w, s)
	}
}

func Info(format string, a ...interface{}) {
	if Level >= 3 {
		a, w := extractLoggerArgs(format, a...)
		s := fmt.Sprintf(label(format, InfoLabel), a...)

		if !TestMode {
			if Color {
				w = color.Output
				s = color.CyanString(s)
			} else if Fabulous {
				w = FabulousWriter
			}
		}

		fmt.Fprintf(w, s)
	}
}

func Success(format string, a ...interface{}) {
	if Level >= 3 {
		a, w := extractLoggerArgs(format, a...)
		s := fmt.Sprintf(label(format, SuccessLabel), a...)

		if !TestMode {
			if Color {
				w = color.Output
				s = color.CyanString(s)
			} else if Fabulous {
				w = FabulousWriter
			}
		}

		fmt.Fprintf(w, s)
	}
}

func Debug(format string, a ...interface{}) {
	if Level >= 4 {
		a, w := extractLoggerArgs(format, a...)
		s := fmt.Sprintf(label(format, DebugLabel), a...)

		if !TestMode {
			if Color {
				w = color.Output
				s = color.GreenString(s)
			} else if Fabulous {
				w = FabulousWriter
			}
		}

		fmt.Fprintf(w, s)

	}
}

func Warning(format string, a ...interface{}) {
	if Level >= 2 {
		a, w := extractLoggerArgs(format, a...)
		s := fmt.Sprintf(label(format, WarningLabel), a...)

		if !TestMode {
			if Color {
				w = color.Output
				s = color.GreenString(s)
			} else if Fabulous {
				w = FabulousWriter
			}
		}

		fmt.Fprintf(w, s)
	}
}

func extractLoggerArgs(format string, a ...interface{}) ([]interface{}, io.Writer) {
	var w io.Writer = os.Stdout

	if n := len(a); n > 0 {
		// extract an io.Writer at the end of a
		if value, ok := a[n-1].(io.Writer); ok {
			w = value
			a = a[0 : n-1]
		}
	}

	return a, w
}

func label(format, label string) string {
	if Timestamps {
		return labelWithTime(format, label)
	} else {
		return labelWithoutTime(format, label)
	}
}

func labelWithTime(format, label string) string {
	t := time.Now()
	rfct := t.Format(time.RFC3339)
	if !strings.Contains(format, "\n") {
		format = fmt.Sprintf("%s%s", format, "\n")
	}
	return fmt.Sprintf("%s [%s]  %s", rfct, label, format)
}

func labelWithoutTime(format, label string) string {
	if !strings.Contains(format, "\n") {
		format = fmt.Sprintf("%s%s", format, "\n")
	}
	return fmt.Sprintf("[%s]  %s", label, format)
}
