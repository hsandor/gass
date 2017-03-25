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

func GetVariables(str string) string {
	rs := rgx.FindStringSubmatch(str)

	if len(rs[1]) <= 0 {
		/*return str*/
		/* throw error */
	}

	return rs[1]
}

func CallFuncByName(funcName, args string) string {
	switch funcName {
	/* STRING */
	case "to-upper-case":
		return ToUpperCase(args)
	case "to-lower-case":
		return ToLowerCase(args)
	case "str-length":
		return StrLength(args)
	/* NUMBER */
	case "random":
		return Random(args)
	}

	return ""
}

func callFunctions(e *element, str string) string {
	for _, funcName := range funcs {
		if !strings.Contains(str, funcName) {
			continue
		}

		// func called!
		vars := GetVariables(str)
		replacable := funcName + "(" + vars + ")"
		res := CallFuncByName(funcName, vars)
		str = strings.Replace(str, replacable, res, 1)

		return str
	}

	return str
}

/* STRING */

// http://sass-lang.com/documentation/Sass/Script/Functions.html#str_length-instance_method
func StrLength(str string) string {
	return fmt.Sprintf("%v", len(str))
}

// http://sass-lang.com/documentation/Sass/Script/Functions.html#to_upper_case-instance_method
func ToUpperCase(str string) string {
	return strings.ToUpper(str)
}

// http://sass-lang.com/documentation/Sass/Script/Functions.html#to_lower_case-instance_method
func ToLowerCase(str string) string {
	return strings.ToLower(str)
}

/* NUMBER */

// http://sass-lang.com/documentation/Sass/Script/Functions.html#random-instance_method
func Random(str string) string {
	return fmt.Sprintf("%v", rand.Intn(2))
}
