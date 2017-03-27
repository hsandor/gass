package gass

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

type parser struct {
	root    *element
	parent  *element
	last    *element
	list    bool // comma separated list of elements
	comment int  // indentation level of starting comment block
	linecnt int
}

func (p *parser) parseElement(indent int, str string) error {
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
		if err := p.last.addName(str[:len(str)-1]); err != nil {
			return err
		}
	} else {
		if err := p.last.addName(str); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseVariable(str string) {
	s := str[strings.Index(str, "$")+1:]
	n := strings.TrimSpace(s[:strings.Index(s, ":")])
	v := strings.TrimSpace(s[strings.Index(s, ":")+1:])
	p.last.addVariable(n, v)
}

func (p *parser) parseProperty(str string) error {
	prop := strings.SplitN(str, ":", 2)
	if len(prop) == 2 {
		if err := p.last.addProperty(prop[0], prop[1]); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseLine(line string) error {
	p.linecnt++
	l := strings.TrimSpace(line)
	if len(l) > 0 {
		indent := calcIndentLevel(line)
		ltype := decideLineType(l)
		if ltype == l_comment {
			p.comment = indent
		} else if p.comment > 0 {
			if indent <= p.comment {
				p.comment = 0
			}
		}

		if p.comment <= 0 {
			l = stripLineComments(l)
			if strings.HasSuffix(l, ";") {
				return errors.New("line shouldn't terminate with semicolon:" + strconv.Itoa(p.linecnt))
			}
			if ltype == l_element {
				if err := p.parseElement(indent, l); err != nil {
					return err
				}
			} else if ltype == l_variable {
				p.parseVariable(l)
			} else if ltype == l_property {
				if p.list {
					return errors.New("open list followed by property:" + strconv.Itoa(p.linecnt))
				}
				p.parseProperty(l)
			}
		}
	}
	return nil
}

func (p *parser) compile(w io.Writer, r io.Reader) error {
	s := bufio.NewScanner(r)
	for s.Scan() {
		if err := p.parseLine(s.Text()); err != nil {
			return err
		}
	}
	p.root.css(w, "", "")
	return nil
}

func newParser() *parser {
	p := new(parser)
	p.root = newElement(-1, nil)
	p.parent = p.root
	p.last = p.root
	return p
}
