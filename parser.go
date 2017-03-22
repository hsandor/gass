package gass

import (
	"errors"
	"fmt"
	"strings"
)

type parser struct {
	errors []string
	root   *element
	parent *element
	last   *element
	list   bool
}

func (p *parser) addError(err error) {
	if err != nil {
		p.errors = append(p.errors, err.Error())
	}
}

func (p *parser) parseElement(indent int, str string) {
	if indent > p.last.indent {
		p.parent = p.last
	} else if indent < p.last.indent {
		for p.parent.indent >= indent && p.parent != p.root {
			p.parent = p.parent.parent
		}
	}
	if !p.list || indent != p.last.indent {
		p.last = newElement(indent, p.parent)
	}
	if p.list = strings.HasSuffix(str, ","); p.list {
		p.addError(p.last.addName(str[:len(str)-1]))
	} else {
		p.addError(p.last.addName(str))
	}
}

func (p *parser) parseProperty(str string) {
	prop := strings.SplitN(str, ":", 2)
	if len(prop) == 2 {
		p.addError(p.last.addProperty(prop[0], prop[1]))
	}
}

func (p *parser) parseLine(line string) error {
	l := strings.TrimSpace(line)
	if len(l) > 0 {
		indent := calcIndentLevel(line)
		ltype := decideLineType(l)

		if ltype == l_element {
			p.parseElement(indent, l)
		} else if ltype == l_property {
			if p.list {
				p.addError(errors.New("open list followed by property:" + line))
				p.list = false
			}
			p.parseProperty(l)
		}
	}
	return nil
}

func (p *parser) parseString(str string) (string, error) {
	for _, line := range strings.Split(str, "\n") {
		if err := p.parseLine(line); err != nil {
			return "", err
		}
	}

	if len(p.errors) > 0 {
		fmt.Println(strings.Join(p.errors, "\n"))
	}

	//dumpElementTree(p.root)

	return p.root.css(""), nil
}

func newParser() *parser {
	p := new(parser)
	p.root = newElement(-1, nil)
	p.parent = p.root
	p.last = p.root
	return p
}

func ParseString(s string) (string, error) {
	p := newParser()
	return p.parseString(s)
}
