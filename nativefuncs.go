package gass

import (
	"fmt"
	"math/rand"
	"strings"
)

var gassFuncs []string = []string{
	"str-length",
	"to-upper-case",
	"to-lower-case",
	"random",
}

var cssFuncs []string = []string{
	"attr",
	"url",
	"calc",
	"linear-gradient",
	"radial-gradient",
	"repeating-linear-gradient",
	"repeating-radial-gradient",
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

func callFunctions(str string) (string, error) {
	openerPos := strings.Index(str, "(")
	result := str

	if openerPos > -1 && openerPos < len(str) {
		part := str[0:openerPos]

		hasGassNative, funcIndex := arrayOfStrContains(gassFuncs, part)

		// got a gass native function
		if hasGassNative {
			funcName := gassFuncs[funcIndex]
			funcPlaceIndex := strings.Index(part, funcName)

			result = str[0:funcPlaceIndex] // leave it

			closerPos := strings.Index(str, ")")

			if closerPos <= -1 {
				// throw error
				fmt.Println("Error: " + string(closerPos))
			}

			// collect the variables
			variables := str[openerPos+1 : closerPos]

			// get the result by function name and the variables
			res, err := callFuncByName(funcName, variables)

			if err == nil {
				result = result + res
			} else {
				fmt.Println("HIBA 2")
			}

			// content still remains
			if closerPos < len(str) {
				remains := str[closerPos+1 : len(str)]

				fmt.Println(remains)

				res, err := callFunctions(remains)

				if err == nil {
					result = result + res
				}
			}
		}
	}

	return result, nil
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
