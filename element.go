package gass

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type property struct {
	name  string
	value string
}

type element struct {
	indent     int
	level      int
	parent     *element
	children   []*element
	names      []string
	properties []property
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

func (e *element) addProperty(name, value string) {
	n := strings.TrimSpace(name)
	v := strings.TrimSpace(value)
	v = callFunctions(e, interpolateVariables(e, v))
	e.properties = append(e.properties, property{n, v})
}

func (e *element) addVariable(name, value string) {
	n := strings.TrimSpace(name)
	v := strings.TrimSpace(value)
	if e.variables == nil {
		e.variables = make(map[string]string)
	}
	e.variables[n] = interpolateVariables(e, v)
}

func (e *element) findVariable(name string) (string, bool) {
	if e.variables != nil {
		if value, ok := e.variables[name]; ok {
			return value, true
		}
	}
	if e.parent != nil {
		return e.parent.findVariable(name)
	}
	return "", false
}

func (e *element) getVariable(name string) string {
	value, ok := e.findVariable(name)
	if !ok {
		fmt.Println("variable doesn't exists:", name)
	}
	return value
}

func (e *element) css(w io.Writer, prefix, previous string) {
	if len(e.properties) > 0 {
		for i, n := range e.names {
			pref, name := resolveAmpersand(prefix, previous, n)
			fmt.Fprint(w, pref)
			if len(strings.TrimSpace(pref)) > 0 {
				fmt.Fprint(w, " ")
			}
			fmt.Fprint(w, name)
			if i < len(e.names)-1 {
				fmt.Fprintln(w, ",")
			}
		}
		fmt.Fprintln(w, " {")
		t := ""
		if prefix == "\t" {
			t = "\t"
		}
		for _, p := range e.properties {
			fmt.Fprintf(w, "\t%s%s: %s;\n", t, p.name, p.value)
		}
		fmt.Fprintln(w, t+"}")
	}
	if len(e.names) > 0 {
		for _, n := range e.names {
			for _, c := range e.children {
				if strings.HasPrefix(n, "@media") {
					fmt.Fprintln(w, n, " {")
					c.css(w, "\t", "")
					fmt.Fprintln(w, "}")
				} else {
					pref, name := resolveAmpersand(prefix, previous, n)
					c.css(w, pref, name)
				}
			}
		}
	} else {
		for _, c := range e.children {
			c.css(w, "", "")
		}
	}
	return
}

func (e *element) gass(w io.Writer) {
	if len(e.names) > 0 {
		t := strings.Repeat("\t", e.level-1)
		for i := 0; i < len(e.names); i++ {
			fmt.Fprint(w, t, e.names[i])
			if i < len(e.names)-1 {
				fmt.Fprintln(w, ",")
			}
		}
		fmt.Fprintln(w, "")
	}
	if e.properties != nil {
		t := strings.Repeat("\t", e.level)
		for n, v := range e.properties {
			fmt.Fprintf(w, "%s%s:%s\n", t, n, v)
		}
		fmt.Println(w, "")
	}
	if len(e.children) > 0 {
		for _, c := range e.children {
			c.gass(w)
		}
	}
}

func newElement(indent int, parent *element) *element {
	e := &element{indent: indent, parent: parent}
	if parent != nil {
		parent.children = append(parent.children, e)
		e.level = parent.level + 1
	}
	return e
}
