package utils

import "strings"

func IsWindowMatch(processName string, includeTitle []string) bool {
	if processName == "repeat-what-shit.exe" {
		return false
	}

	if len(includeTitle) == 0 {
		return true
	}

	processName = strings.ToLower(processName)

	for _, title := range includeTitle {
		if strings.ToLower(title) == processName {
			return true
		}
	}

	return false
}
