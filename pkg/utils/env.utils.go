package utils

import (
	"os"
)

func GetEnvVar(vars string, defaultVal string) string {
	tmp := os.Getenv(vars)

	if len(tmp) == 0 {
		return defaultVal
	}
	return tmp
}

func Remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
