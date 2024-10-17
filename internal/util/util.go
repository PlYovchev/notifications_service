package util

import (
	"fmt"
	"strings"
)

// IsDevMode - Checks if the given string denotes any of the development environment.
func IsDevMode(s string) bool {
	return strings.Contains(s, "local") || strings.Contains(s, "dev")
}

func ArrayToString(a []int, delim string) string {
	return strings.Trim(strings.ReplaceAll(fmt.Sprint(a), " ", delim), "[]")
}

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}
