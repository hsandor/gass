package gass

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

var funcs []string = []string{
	"str-length",
	"to-upper-case",
	"to-lower-case",
	"random",
}

var rgx = regexp.MustCompile(`\((.*?)\)`)

func getVariables(str string) string {
	rs := rgx.FindStringSubmatch(str)

	if len(rs[1]) <= 0 {
		/*return str*/
		/* throw error */
	}

	return rs[1]
}

func callFuncByName(funcName, args string) string {
	switch funcName {
	/* STRING */
	case "to-upper-case":
		return toUpperCase(args)
	case "to-lower-case":
		return toLowerCase(args)
	case "str-length":
		return strLength(args)
	/* NUMBER */
	case "random":
		return random(args)
	}

	return ""
}

func callFunctions(e *element, str string) string {
	for _, funcName := range funcs {
		if !strings.Contains(str, funcName) {
			continue
		}

		// func called!
		vars := getVariables(str)
		replacable := funcName + "(" + vars + ")"
		res := callFuncByName(funcName, vars)
		str = strings.Replace(str, replacable, res, 1)

		return str
	}

	return str
}

/* STRING */

// http://sass-lang.com/documentation/Sass/Script/Functions.html#str_length-instance_method
func strLength(str string) string {
	if !isGassStr(str) {
		// error
	}

	return fmt.Sprintf("%v", len(str))
}

// http://sass-lang.com/documentation/Sass/Script/Functions.html#to_upper_case-instance_method
func toUpperCase(str string) string {
	return strings.ToUpper(str)
}

// http://sass-lang.com/documentation/Sass/Script/Functions.html#to_lower_case-instance_method
func toLowerCase(str string) string {
	return strings.ToLower(str)
}

/* NUMBER */

// http://sass-lang.com/documentation/Sass/Script/Functions.html#random-instance_method
func random(str string) string {
	return fmt.Sprintf("%v", rand.Intn(2))
}
