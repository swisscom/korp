package string_utils

import (
	"strings"
)

// TrimQuotes - Remove single or double quotes from a string
func TrimQuotes(str string) string {

	return strings.Trim(str, "'\"")
}
