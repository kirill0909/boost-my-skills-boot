package utils

import (
	"strings"
)

func FormatSnipets(message string) string {
	switch {
	case strings.Contains(message, SnipetGoBegin):
		message = strings.Replace(message, SnipetGoBegin, "```go\n", -1)
		message = strings.Replace(message, SnipetEnd, "```\n", -1)
	case strings.Contains(message, SnipetBashBegin):
		message = strings.Replace(message, SnipetBashBegin, "```bash\n", -1)
		message = strings.Replace(message, SnipetEnd, "```\n", -1)
	case strings.Contains(message, SnipetSQLBegin):
		message = strings.Replace(message, SnipetSQLBegin, "```sql\n", -1)
		message = strings.Replace(message, SnipetEnd, "```\n", -1)
	case strings.Contains(message, SnipetRustBegin):
		message = strings.Replace(message, SnipetRustBegin, "```rust\n", -1)
		message = strings.Replace(message, SnipetEnd, "```\n", -1)
	}

	return message
}

func FormatBadCharacters(text string) string {
	text = strings.ReplaceAll(text, "{", "\\{")
	text = strings.ReplaceAll(text, "}", "\\}")
	text = strings.ReplaceAll(text, "(", "\\(")
	text = strings.ReplaceAll(text, ")", "\\)")
	text = strings.ReplaceAll(text, "[", "\\[")
	text = strings.ReplaceAll(text, "]", "\\]")
	text = strings.ReplaceAll(text, "*", "\\*")
	text = strings.ReplaceAll(text, "%", "\\%")
	text = strings.ReplaceAll(text, "$", "\\$")
	text = strings.ReplaceAll(text, "&", "\\&")
	text = strings.ReplaceAll(text, "#", "\\#")
	text = strings.ReplaceAll(text, ".", "\\.")
	text = strings.ReplaceAll(text, ">", "\\>")
	text = strings.ReplaceAll(text, "<", "\\<")
	text = strings.ReplaceAll(text, "?", "\\?")
	text = strings.ReplaceAll(text, "!", "\\!")
	text = strings.ReplaceAll(text, "+", "\\+")
	text = strings.ReplaceAll(text, "-", "\\-")
	text = strings.ReplaceAll(text, "_", "\\_")
	text = strings.ReplaceAll(text, ";", "\\;")
	text = strings.ReplaceAll(text, "/", "\\/")
	text = strings.ReplaceAll(text, "@", "\\@")
	text = strings.ReplaceAll(text, "=", "\\=")
	text = strings.ReplaceAll(text, "|", "\\|")

	return text
}
