package logs

import (
	"cerbero3/configuration"
	"cerbero3/logs/colors"
	"cerbero3/logs/logfile"
	"fmt"
	"os"
	"time"
)

var (
	verboseEnabled bool
	logFilePath    string
	logFile        *os.File
)

type coloringFunction func(string) string

func Configure(c configuration.Configuration) error {
	verboseEnabled = c.Verbose
	logFilePath = c.LogFile

	// the explicit declaration is necessary so Go
	// can understand that logFile is not a new variable
	var err error
	logFile, err = logfile.Create(logFilePath)
	if err != nil {
		return err
	}

	colors.Configure(c)
	return nil
}

func print(message string, cf coloringFunction) {
	line := fmt.Sprintf("%v %v", time.Now().UTC(), message)

	// write to logfile first and then to stdout
	// TODO: handle errors in case it is not able to write to file
	logFile.WriteString(fmt.Sprintf("%v\n", line))
	fmt.Printf("%v\n", cf(line))
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
