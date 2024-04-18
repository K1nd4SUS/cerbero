package logs

import (
	"cerbero/configuration"
	"cerbero/logs/colors"
	"cerbero/logs/logfile"
	"fmt"
	"os"
	"time"
)

var (
	configured     bool
	verboseEnabled bool
	logFilePath    string
	logFile        *os.File
	failedLogs     string
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

	configured = true
	return nil
}

func print(message string, cf coloringFunction, ignoreErrors bool) {
	line := fmt.Sprintf("%v %v", time.Now().UTC(), message)

	// check if there are any failed logs; if yes,
	// try to write them again and if no errors occur,
	// clear them
	if failedLogs != "" {
		_, err := logFile.WriteString(failedLogs)
		if err == nil {
			failedLogs = ""
			// TODO: this always prints when starting the
			// firewall because when it's not configured
			// it collects some errors and writes the logs
			// to the logfile at a later time
			PrintInfo("Failed logs written successfully.")
		}
	}
	// write to logfile first and then to stdout
	// TODO: when the file is deleted, this does not return
	// any errors, not even when the file is moved
	_, err := logFile.WriteString(fmt.Sprintf("%v\n", line))
	if err != nil && !ignoreErrors {
		// only print the warning if the logs have
		// already been configured; if they haven't,
		// the logfile isn't loaded yet so it will always
		// return errors. save the unwritten logs anyways
		// so they can be written at a later time
		if configured {
			printSpecialWarning(fmt.Sprintf("Cerbero was not able to write the upcoming log line to the logfile. Fix the issue as soon as possible, another try to write will be made on the next log. Here is the error: %v.", err.Error()))
		}
		failedLogs += fmt.Sprintf("%v\n", line)
	}
	fmt.Printf("%v\n", cf(line))
}

func PrintDebug(message string) {
	if verboseEnabled {
		print(fmt.Sprintf("[DEBUG] %v", message), colors.GetDebugColored, false)
	}
}

func PrintInfo(message string) {
	print(fmt.Sprintf("[INFO] %v", message), colors.GetInfoColored, false)
}

func PrintWarning(message string) {
	print(fmt.Sprintf("[WARNING] %v", message), colors.GetWarningColored, false)
}

func PrintError(message string) {
	print(fmt.Sprintf("[ERROR] %v", message), colors.GetErrorColored, false)
}

func PrintCritical(message string) {
	print(fmt.Sprintf("[CRITICAL] %v", message), colors.GetCriticalColored, false)
}

// a "special warning" is a warning that ignores
// errors in writing to the logfile
// TODO: try to format this document better in order
// not to repeat any functions (this and PrintWarning)
func printSpecialWarning(message string) {
	print(fmt.Sprintf("[WARNING] %v", message), colors.GetWarningColored, true)
}
