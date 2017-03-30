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

func callFuncByName(funcName, args string) (res string, err error) {
	res = ""
	err = nil

	switch funcName {
	/* STRING */
	case "to-upper-case":
		res, err = toUpperCase(args)
	case "to-lower-case":
		res, err = toLowerCase(args)
	case "str-length":
		res, err = strLength(args)
	/* NUMBER */
	case "random":
		res, err = random(args)
	}

	return res, err
}

func callFunctions(e *element, str string) (string, error) {
	for _, funcName := range funcs {
		if !strings.Contains(str, funcName) {
			continue
		}

		// func called!
		vars := getVariables(str)
		replacable := funcName + "(" + vars + ")"
		res, err := callFuncByName(funcName, vars)

		if err != nil {
			return str, err
		}

		str = strings.Replace(str, replacable, res, 1)

		return str, nil
	}

	return str, nil
}

/* STRING */

// http://sass-lang.com/documentation/Sass/Script/Functions.html#str_length-instance_method
func strLength(str string) (string, error) {
	if _, err := isGassStr(str); err != nil {
		return str, err
	}

	return fmt.Sprintf("%v", len(str)), nil
}

// http://sass-lang.com/documentation/Sass/Script/Functions.html#to_upper_case-instance_method
func toUpperCase(str string) (string, error) {
	if _, err := isGassStr(str); err != nil {
		return str, err
	}

	return strings.ToUpper(str), nil
}

// http://sass-lang.com/documentation/Sass/Script/Functions.html#to_lower_case-instance_method
func toLowerCase(str string) (string, error) {
	if _, err := isGassStr(str); err != nil {
		return str, err
	}

	return strings.ToLower(str), nil
}

/* NUMBER */

// http://sass-lang.com/documentation/Sass/Script/Functions.html#random-instance_method
func random(str string) (string, error) {
	return fmt.Sprintf("%v", rand.Intn(2)), nil
}
