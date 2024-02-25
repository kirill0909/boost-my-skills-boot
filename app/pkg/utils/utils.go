package utils

import (
	"strings"
)

func FormatSnipets(message string) string {
	message = strings.Replace(message, SnipetBegin, "```go\n", -1)
	message = strings.Replace(message, SnipetEnd, "```\n", -1)

	return message
}
