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

func (e *element) css(prefix []string) string {
	res := ""

	if len(e.properties) > 0 {
		for i, n := range e.names {
			cn := resolveAmpersand(prefix, n)

			if cn == n {
				res += strings.Join(prefix, " ")
				if len(prefix) > 0 {
					res += " "
				}
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
	}

	if len(e.names) > 0 {
		for _, n := range e.names {
			for _, c := range e.children {
				res += c.css(append(prefix, resolveAmpersand(prefix, n)))
			}
		}
	} else {
		for _, c := range e.children {
			res += c.css(prefix)
		}
	}

	return res
}

func newElement(indent int, parent *element) *element {
	e := &element{indent: indent, parent: parent}
	if parent != nil {
		parent.children = append(parent.children, e)
		e.level = parent.level + 1
	}
	return e
}
