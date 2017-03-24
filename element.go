package gass

import (
	"errors"
	"fmt"
	"strings"
)

type element struct {
	indent     int
	level      int
	parent     *element
	children   []*element
	names      []string
	properties map[string]string
}

func (e *element) addName(name string) error {
	for _, n := range e.names {
		if n == name {
			return errors.New("name already exists:" + name)
		}
	}
	e.names = append(e.names, name)
	return nil
}

func (e *element) addProperty(name, value string) error {
	n := strings.TrimSpace(name)
	v := strings.TrimSpace(value)

	if e.properties == nil {
		e.properties = make(map[string]string)
	} else {
		if _, exists := e.properties[n]; exists {
			return errors.New("property already exists:" + n)
		}
	}

	e.properties[n] = v
	return nil
}

func (e *element) cssProperties(prefix, previous string) (res string) {
	res = ""

	if len(e.properties) <= 0 {
		return
	}

	for i, n := range e.names {

		cn := strings.Replace(n, "&", previous, -1)

		res += prefix

		if len(prefix) > 0 {
			res += " "
		}

		if n == cn {
			res += previous
			res += " "
		}

		res += cn

		if i < len(e.names)-1 {
			res += ",\n"
		}
	}

	res += " {\n"
	for n, v := range e.properties {
		res += fmt.Sprintf("  %s: %s;\n", n, v)
	}
	res += "}\n"

	return
}

func (e *element) cssChildren(prefix, previous string) (res string) {
	res = ""

	if len(e.names) > 0 {
		for _, n := range e.names {
			for _, c := range e.children {
				cn := strings.Replace(n, "&", previous, -1)
				p := ""

				if cn == n {
					p += prefix
					if len(prefix) > 0 {
						p += " "
					}
					p += previous
				}

				res += c.css(p, cn)
			}
		}
	} else {
		for _, c := range e.children {
			res += c.css("", "")
		}
	}

	return
}

func (e *element) css(prefix, previous string) (res string) {
	res = e.cssProperties(prefix, previous)
	return res + e.cssChildren(prefix, previous)
}

func (e *element) gass() (res string) {
	if len(e.names) > 0 {
		pre := strings.Repeat("\t", e.level-1)
		for i := 0; i < len(e.names); i++ {
			res += pre + e.names[i]
			if i < len(e.names)-1 {
				res += ",\n"
			}
		}
		res += "\n"
	}
	if e.properties != nil {
		for n, v := range e.properties {
			res += fmt.Sprintf("%s%s:%s\n", strings.Repeat("\t", e.level), n, v)
		}
		res += "\n"
	}
	if len(e.children) > 0 {
		for _, c := range e.children {
			res += c.gass()
		}
	}
	return
}

func newElement(indent int, parent *element) *element {
	e := &element{indent: indent, parent: parent}
	if parent != nil {
		parent.children = append(parent.children, e)
		e.level = parent.level + 1
	}
	return e
}
