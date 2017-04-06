// line
package line

import (
	"strings"
)

type Line struct {
	Indent int
	Text   string
}

func (l *Line) Parse(raw string) {
	l.Text = strings.TrimSpace(raw)
	l.calcIndentLevel(raw)
}

func (l *Line) calcIndentLevel(s string) {
	i := 0

	if len(s) > 0 {
		for _, c := range s {
			if c == 32 {
				i += 1
			} else if c == 9 {
				i = ((i + 8) / 8) * 8
			} else if c != 13 {
				break
			}
		}
	}

	l.Indent = i
}
