package gass

import (
	"sort"
	"strings"
)

const (
	l_element = iota
	l_property
	l_variable
	l_comment
	l_media
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
	if strings.HasPrefix(s, "//") || strings.HasPrefix(s, "/*") {
		return l_comment
	} else if strings.HasPrefix(s, "$") && strings.Contains(s, ":") {
		return l_variable
	} else if !strings.HasPrefix(s, "&") && !strings.HasPrefix(s, "@") && strings.Contains(s, ":") {
		return l_property
	} else {
		return l_element
	}
}

func resolveAmpersand(prefix, previous, currname string) (pref, name string) {
	pref = prefix
	name = strings.Replace(currname, "&", previous, -1)
	if name == currname {
		if len(strings.TrimSpace(pref)) > 0 {
			pref += " "
		}
		pref += previous
	}
	return
}

func interpolateVariables(e *element, str string) (res string) {
	if strings.Contains(str, "$") {
		vars := strings.Split(str, "$")
		for i, v := range vars {
			if i > 0 {
				n := strings.IndexAny(v, " @#&(){}[];:,./")
				if n >= 0 {
					res += e.getVariable(v[:n])
					res += v[n:]
				} else {
					res += e.getVariable(v)
				}
			} else {
				res += v
			}
		}
		return
	}
	return str
}

func sortedRange(m map[string]string) (sorted []string) {
	for k, _ := range m {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	return
}

func stripLineComments(s string) string {
	if pos := strings.Index(s, "//"); pos >= 0 {
		return s[:pos]
	}
	return s
}

func isGassStr(str string) bool {
	if strings.HasPrefix(str, `"`) && strings.HasSuffix(str, `"`) {
		if strings.Count(str, `"`) == 2 {
			return true
		}
	}

	return false
}
