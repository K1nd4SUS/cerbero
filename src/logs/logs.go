package logs

import (
	"cerbero3/configuration"
	"cerbero3/logs/colors"
	"fmt"
	"time"
)

var (
	verboseEnabled bool
)

type coloringFunction func(string) string

func Configure(c configuration.Configuration) {
	verboseEnabled = c.Verbose
}

func print(message string, cf coloringFunction) {
	fmt.Printf("%v\n", cf(fmt.Sprintf("%v %v", time.Now().UTC(), message)))
}

func PrintDebug(message string) {
	if verboseEnabled {
		print(fmt.Sprintf("[DEBUG] %v", message), colors.GetDebugColored)
	}
}

func PrintInfo(message string) {
	print(fmt.Sprintf("[INFO] %v", message), colors.GetInfoColored)
}

func PrintWarning(message string) {
	print(fmt.Sprintf("[WARNING] %v", message), colors.GetWarningColored)
}

func PrintError(message string) {
	print(fmt.Sprintf("[ERROR] %v", message), colors.GetErrorColored)
}

func PrintCritical(message string) {
	print(fmt.Sprintf("[CRITICAL] %v", message), colors.GetCriticalColored)
}
