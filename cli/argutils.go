package cli

import (
	"regexp"
	"strings"
)

func ArgSplit(str string) (result []string) {
	r := regexp.MustCompile(`[^\s"]+|"([^"]*)"`)
	match := r.FindAllString(str, -1)
	for _, s := range match {
		result = append(result, strings.Trim(strings.TrimSpace(s), `" `))
	}
	return
}
