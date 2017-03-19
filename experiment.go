package gass

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	ELEMENT = iota
	PROPERTY
	COMMENT
	VARIABLE
)

type Line struct {
	Type       int
	Parent     *Line
	Data       string
	Level      int
	Properties []*Line
	Children   []*Line
}

var variables map[string]string

func CalcIndent(s string) int {
	i := 0
	for _, c := range s {
		if c == 32 {
			i += 1
		} else if c == 9 {
			i += 1
			for i%8 != 0 {
				i += 1
			}
		} else if c != 13 {
			break
		}
	}
	return i
}

func CalcType(s string) int {
	ss := strings.TrimSpace(s)

	if strings.HasPrefix(ss, "//") || strings.HasPrefix(ss, "/*") {
		return COMMENT
	} else if strings.HasPrefix(ss, "$") && strings.Contains(ss, ":") {
		return VARIABLE
	} else if strings.Contains(ss, ":") {
		return PROPERTY
	} else {
		return ELEMENT
	}
}

func NewLine(data string) *Line {
	return &Line{
		Data:  strings.TrimSpace(data),
		Level: CalcIndent(data),
		Type:  CalcType(data),
	}
}

func FormatChain(l *Line) string {
	if l != nil {
		s := FormatChain(l.Parent)
		if len(s) > 0 {
			return s + " " + strings.TrimRight(l.Data, ",")
		}
		return strings.TrimRight(l.Data, ",")
	}
	return ""
}

func FormatOutput(l *Line, p string) string {
	res := ""
	if l != nil {
		if len(l.Data) > 0 && len(l.Properties) > 0 {
			res += fmt.Sprintf("%s {\n", FormatChain(l))
			for _, p := range l.Properties {
				res += "  " + p.Data + ";\n"
			}
			res += "}\n"
		} else {
			for _, c := range l.Children {
				res += FormatOutput(c, p+"  ")
			}
		}
	}
	return res
}

func ParseVariable(s string) {

}

func SetProp(p *Line, l *Line) {
	p.Properties = append(p.Properties, l)

	for i := len(p.Parent.Children) - 2; i >= 0 && strings.HasSuffix(p.Parent.Children[i].Data, ","); i -= 1 {
		p.Parent.Children[i].Properties = append(p.Parent.Children[i].Properties, l)
	}
}

func CompileString(src string) string {
	root := NewLine("")
	last := root
	parent := root
	var comment *Line

	for _, line := range strings.Split(src, "\n") {
		if len(strings.TrimSpace(line)) > 0 {
			lin := NewLine(line)

			if lin.Type == COMMENT {
				comment = lin
			} else if comment != nil {
				if lin.Level <= comment.Level {
					comment = nil
				}
			}

			if comment != nil {
				continue
			}

			if lin.Type == ELEMENT {
				if lin.Level > last.Level {
					parent = last
				} else if lin.Level < last.Level {
					for parent.Level >= lin.Level && parent != root {
						parent = parent.Parent
					}
				}

				parent.Children = append(parent.Children, lin)
				lin.Parent = parent
				last = lin
			} else {
				SetProp(last, lin)
			}
		}
	}

	return FormatOutput(root, "")
}

func CompileFile(name string) {
	src, err := ioutil.ReadFile(name)

	if err != nil {
		fmt.Println(err)
		return
	}

	css := CompileString(string(src))
	nfn := strings.TrimSuffix(name, filepath.Ext(name)) + ".css"

	ioutil.WriteFile(nfn, []byte(css), 0)
}
