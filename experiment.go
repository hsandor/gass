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
	Variables  map[string]string
}

func CalcType(s string) int {
	ss := strings.TrimSpace(s)
	hasAmp := strings.Contains(ss, "&")

	if strings.HasPrefix(ss, "//") || strings.HasPrefix(ss, "/*") {
		return COMMENT
	} else if strings.HasPrefix(ss, "$") && strings.Contains(ss, ":") && !hasAmp {
		return VARIABLE
	} else if strings.Contains(ss, ":") && !hasAmp {
		return PROPERTY
	} else {
		return ELEMENT
	}
}

func NewLine(data string) *Line {
	return &Line{
		Data:      strings.TrimSpace(data),
		Level:     calcIndentLevel(data),
		Type:      CalcType(data),
		Variables: make(map[string]string),
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
		if len(l.Properties) > 0 {
			res += fmt.Sprintf("%s {\n", FormatChain(l))
			for _, p := range l.Properties {
				res += "  " + p.Data + ";\n"
			}
			res += "}\n"
		}

		if len(l.Children) > 0 {
			for _, c := range l.Children {
				res += FormatOutput(c, p+"  ")
			}
		}
	}
	return res
}

func InterpolateVariables(l *Line, s string) string {
	vars := strings.Split(s, "$")

	if len(vars) <= 1 {
		return s
	}

	res := ""

	for _, v := range vars {
		i := strings.IndexAny(v, " /,:.#")
		if i >= 0 {
			res += FindVariable(l, v[:i])
			res += v[i:]
		} else {
			res += FindVariable(l, v)
		}
	}

	return res
}

func ParseVariable(l *Line, s string) {
	s = s[strings.Index(s, "$")+1:]
	n := strings.TrimSpace(s[:strings.Index(s, ":")])
	v := strings.TrimSpace(s[strings.Index(s, ":")+1:])
	l.Variables[n] = InterpolateVariables(l, v)
}

func FindVariable(l *Line, name string) string {
	for l != nil {
		if value, ok := l.Variables[name]; ok {
			return value
		}
		l = l.Parent
	}
	return name
}

func SetProp(p *Line, l *Line) {
	p.Properties = append(p.Properties, l)

	for i := len(p.Parent.Children) - 2; i >= 0 && strings.HasSuffix(p.Parent.Children[i].Data, ","); i-- {
		p.Parent.Children[i].Properties = append(p.Parent.Children[i].Properties, l)
	}
}

func ReplaceAmpersand(lin *Line) {
	if lin.Parent != nil {
		nd := strings.Replace(lin.Data, "&", lin.Parent.Data, -1)
		if nd != lin.Data {
			lin.Parent = lin.Parent.Parent
			lin.Data = nd
		}
	} else if strings.Contains(lin.Data, "&") {
		fmt.Println("ampersand for top level element error:", lin.Data)
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

				lin.Parent = parent
				ReplaceAmpersand(lin)
				lin.Parent.Children = append(lin.Parent.Children, lin)
				last = lin
			} else if lin.Type == VARIABLE {
				ParseVariable(last, lin.Data)
			} else {
				if last != root {
					lin.Data = InterpolateVariables(last, lin.Data)
					SetProp(last, lin)
				} else {
					fmt.Println("error: top level property '" + lin.Data + "'")
				}
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
