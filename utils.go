package gass

import (
	"strings"
)

const (
	l_element = iota
	l_property
	l_variable
)

func calcIndentLevel(s string) int {
	i := 0
	for _, c := range s {
		if c == 32 {
			i += 1
		} else if c == 9 {
			i = ((i + 8) / 8) * 8
		} else if c != 13 {
			break
		}
	}
	return i
}

func decideLineType(s string) int {
	if !strings.HasPrefix(s, "&") && strings.Contains(s, ":") {
		return l_property
	} else {
		return l_element
	}
}
