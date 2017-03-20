package gass

func calcIndentLevel(s string) int {
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
