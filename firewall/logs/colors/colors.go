package colors

import "fmt"

func GetDebugColored(text string) string {
	return fmt.Sprintf("\x1b[38;5;240m%v\x1b[0m", text)
}

func GetInfoColored(text string) string {
	return text
}

func GetWarningColored(text string) string {
	return fmt.Sprintf("\x1b[33m%v\x1b[0m", text)
}

func GetErrorColored(text string) string {
	return fmt.Sprintf("\x1b[31m%v\x1b[0m", text)
}

func GetCriticalColored(text string) string {
	return fmt.Sprintf("\x1b[39;41m%v\x1b[0m", text)
}
