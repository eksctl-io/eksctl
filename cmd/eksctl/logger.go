package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/kris-nova/logger"
	lol "github.com/kris-nova/lolgopher"
)

func initLogger(level int, colorValue string) {
	logger.Layout = "2006-01-02 15:04:05"

	var bitwiseLevel int
	switch level {
	case 4:
		bitwiseLevel = logger.LogDeprecated | logger.LogAlways | logger.LogSuccess | logger.LogCritical | logger.LogWarning | logger.LogInfo | logger.LogDebug
	case 3:
		bitwiseLevel = logger.LogDeprecated | logger.LogAlways | logger.LogSuccess | logger.LogCritical | logger.LogWarning | logger.LogInfo
	case 2:
		bitwiseLevel = logger.LogDeprecated | logger.LogAlways | logger.LogSuccess | logger.LogCritical | logger.LogWarning
	case 1:
		bitwiseLevel = logger.LogDeprecated | logger.LogAlways | logger.LogSuccess | logger.LogCritical
	case 0:
		bitwiseLevel = logger.LogDeprecated | logger.LogAlways | logger.LogSuccess
	default:
		bitwiseLevel = logger.LogDeprecated | logger.LogEverything
	}
	logger.BitwiseLevel = bitwiseLevel

	switch colorValue {
	case "fabulous":
		logger.Writer = lol.NewLolWriter()
	case "true":
		logger.Writer = color.Output
	}

	logger.Line = func(prefix, format string, a ...interface{}) string {
		if !strings.Contains(format, "\n") {
			format = fmt.Sprintf("%s%s", format, "\n")
		}
		now := time.Now()
		fNow := now.Format(logger.Layout)
		var colorize func(format string, a ...interface{}) string
		var icon string
		switch prefix {
		case logger.PreAlways:
			icon = "✿"
			colorize = color.GreenString
		case logger.PreCritical:
			icon = "✖"
			colorize = color.RedString
		case logger.PreInfo:
			icon = "ℹ"
			colorize = color.CyanString
		case logger.PreDebug:
			icon = "▶"
			colorize = color.GreenString
		case logger.PreSuccess:
			icon = "✔"
			colorize = color.CyanString
		case logger.PreWarning:
			icon = "!"
			colorize = color.GreenString
		default:
			icon = "ℹ"
			colorize = color.CyanString
		}

		out := fmt.Sprintf(format, a...)
		out = fmt.Sprintf("%s [%s]  %s", fNow, icon, out)
		if colorValue == "true" {
			out = colorize(out)
		}

		return out
	}
}
