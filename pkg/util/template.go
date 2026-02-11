package util

import (
	"strconv"
	"strings"
)

func FormatSEK(n int) string {
	s := strconv.Itoa(n)
	var out []string
	for len(s) > 3 {
		out = append([]string{s[len(s)-3:]}, out...)
		s = s[:len(s)-3]
	}
	if len(s) > 0 {
		out = append([]string{s}, out...)
	}
	return strings.Join(out, " ")
}