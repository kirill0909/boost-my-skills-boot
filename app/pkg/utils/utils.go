package utils

import (
	"regexp"
)

func FormatSnipets(message string) string {
	reBegin := regexp.MustCompile("snipet (go|bash|sql|rust|proto|yaml|json) begin")
	reEnd := regexp.MustCompile("snipet end")

	message = reBegin.ReplaceAllString(message, "```$1")
	message = reEnd.ReplaceAllString(message, "```")

	return message
}

func FormatBadCharacters(text string) string {
	re := regexp.MustCompile(`([{}()\[\]*%$&#!.<>\?+\-_;/@=|])`)
	text = re.ReplaceAllString(text, "\\$1")

	return text
}
