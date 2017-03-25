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
	variables  map[string]string
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
	v = callFunctions(e, v)
	e.properties[n] = interpolateVariables(e, v)
	return nil
}

func (e *element) addVariable(name, value string) {
	n := strings.TrimSpace(name)
	v := strings.TrimSpace(value)
	if e.variables == nil {
		e.variables = make(map[string]string)
	}
	e.variables[n] = interpolateVariables(e, v)
	fmt.Println("new variable:", n, ":", e.variables[n])
}

func (e *element) getVariable(name string) (value string) {
	if e.variables != nil {
		if value, ok := e.variables[name]; ok {
			return value
		}
	}
	if e.parent != nil {
		return e.parent.getVariable(name)
	}
	return ""
}

func (e *element) cssProperties(prefix, previous string) (res string) {
	res = ""
	if len(e.properties) > 0 {
		for i, n := range e.names {
			pref, name := resolveAmpersand(prefix, previous, n)
			res += pref
			if len(res) > 0 {
				res += " "
			}
			res += name
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
	return
}

func (e *element) cssChildren(prefix, previous string) (res string) {
	res = ""
	if len(e.names) > 0 {
		for _, n := range e.names {
			for _, c := range e.children {
				pref, name := resolveAmpersand(prefix, previous, n)
				res += c.css(pref, name)
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
	return e.cssProperties(prefix, previous) + e.cssChildren(prefix, previous)
}

func (e *element) gass() (res string) {
	if len(e.names) > 0 {
		t := strings.Repeat("\t", e.level-1)
		for i := 0; i < len(e.names); i++ {
			res += t + e.names[i]
			if i < len(e.names)-1 {
				res += ",\n"
			}
		}
		res += "\n"
	}
	if e.properties != nil {
		t := strings.Repeat("\t", e.level)
		for n, v := range e.properties {
			res += fmt.Sprintf("%s%s:%s\n", t, n, v)
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
