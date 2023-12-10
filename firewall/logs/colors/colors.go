package colors

import (
	"cerbero3/configuration"
	"fmt"
)

var colorsEnabled bool

func Configure(c configuration.Configuration) {
	colorsEnabled = c.ColoredLogs
}

func getColored(color, text string) string {
	if !colorsEnabled {
		return text
	}
	return fmt.Sprintf(color, text)
}

func GetDebugColored(text string) string {
	return getColored("\x1b[38;5;240m%v\x1b[0m", text)
}

func GetInfoColored(text string) string {
	return text
}

func GetWarningColored(text string) string {
	return getColored("\x1b[33m%v\x1b[0m", text)
}

func GetErrorColored(text string) string {
	return getColored("\x1b[31m%v\x1b[0m", text)
}

func GetCriticalColored(text string) string {
	return getColored("\x1b[39;41m%v\x1b[0m", text)
}
