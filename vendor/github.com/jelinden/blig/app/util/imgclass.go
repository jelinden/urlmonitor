package util

import "strings"

func ImgClass(html string) string {
	stripped := stripExtraLines(html)
	count := strings.Count(stripped, "img src=")
	for i := 0; i < count; i++ {
		index := strings.Index(stripped, "img src=") + 4
		stripped = stripped[:index] + " class=\"pure-img\"" + stripped[index:]
	}
	return stripped
}

func stripExtraLines(html string) string {
	return strings.Replace(html, "\n\n", "\n", -1)
}
