package util

import (
	"regexp"
	"strings"

	"golang.org/x/text/unicode/norm"
)

var nonLatin = regexp.MustCompile("[^\\w-]")
var whitespaceAndComma = regexp.MustCompile("[\\s+,/]")
var multipleDashes = regexp.MustCompile("-+")
var startEndDashes = regexp.MustCompile("^-|-$")

func Slugify(title string) string {
	title = whitespaceAndComma.ReplaceAllString(title, "-")
	title = norm.NFD.String(title)
	title = nonLatin.ReplaceAllString(title, "")
	title = whitespaceAndComma.ReplaceAllString(title, "-")
	title = multipleDashes.ReplaceAllString(title, "-")
	title = startEndDashes.ReplaceAllString(title, "")
	return strings.ToLower(title)
}
