package gass

import (
	"fmt"
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

func dumpElementTree(e *element) {
	if len(e.names) > 0 {
		fmt.Print(strings.Repeat("\t", e.level-1))
		fmt.Println(strings.Join(e.names, ","))
	}
	if e.properties != nil {
		for n, v := range e.properties {
			fmt.Printf("%s%s:%s\n", strings.Repeat("\t", e.level), n, v)
		}
		fmt.Println("")
	}
	if len(e.children) > 0 {
		for _, c := range e.children {
			dumpElementTree(c)
		}
	}
}
