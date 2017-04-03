package gass

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

var gassFuncs []string = []string{
	"unquote",
	"str-length",
	"to-upper-case",
	"to-lower-case",
	"random",
}

// http://www.quackit.com/css/functions/
var cssFuncs []string = []string{
	"attr",
	"blur",
	"brightness",
	"calc",
	"circle",
	"contrast",
	"counter",
	"counters",
	"cubic-bezier",
	"drop-shadow",
	"ellipse",
	"format",
	"grayscale",
	"hsl",
	"hsla",
	"hue-rotate",
	"hwb",
	"image",
	"inset",
	"invert",
	"linear-gradient",
	"matrix",
	"matrix3d",
	"opacity",
	"perspective",
	"polygon",
	"radial-gradient",
	"repeating-linear-gradient",
	"repeating-radial-gradient",
	"rgb",
	"rgba",
	"rotate",
	"rotate3d",
	"rotateX",
	"rotateY",
	"rotateZ",
	"saturate",
	"sepia",
	"scale",
	"scale3d",
	"scaleX",
	"scaleY",
	"scaleZ",
	"skew",
	"skewX",
	"skewY",
	"symbols",
	"translate",
	"translate3d",
	"translateX",
	"translateY",
	"translateZ",
	"url",
}

func callFuncByName(funcName, args string) (res string, err error) {
	res = ""
	err = nil

	switch funcName {
	/* STRING */
	case "unquote":
		res, err = unquote(args)
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

		isGassNative, funcIndex := arrayOfStrContains(gassFuncs, part)

		// got a gass native function
		if isGassNative {
			funcName := gassFuncs[funcIndex]
			funcPlaceIndex := strings.Index(part, funcName)

			result = str[0:funcPlaceIndex] // leave it

			closerPos := strings.Index(str, ")")

			if closerPos <= -1 {
				// throw error
				fmt.Println("Error: " + string(closerPos))
			}

			// collect the arguments
			arguments := str[openerPos+1 : closerPos]

			if strings.Contains(arguments, "(") {
				vars := str[openerPos+1 : closerPos+1]

				res, err := callFunctions(vars)

				if err == nil {
					arguments = res
				}
			}

			// get the result by function name and the arguments
			res, err := callFuncByName(funcName, arguments)

			if err == nil {
				result = result + res
			} else {
				return result, err
			}

			// content still remains
			if closerPos < len(str) {
				remains := str[closerPos+1 : len(str)]

				if remains != ")" {
					res, err := callFunctions(remains)

					if err == nil {
						result = result + res
					}
				}
			}
		} else if isCssNative, _ := arrayOfStrContains(cssFuncs, part); isCssNative {
			nativeArgs := str[len(part)+1 : len(str)]

			if strings.Contains(nativeArgs, "(") {
				res, err := callFunctions(nativeArgs)

				if err == nil {
					result = part + "(" + res

					if !strings.HasSuffix(result, ")") {
						result = result + ")"
					}
				} else {
					return result, err
				}
			}
		}
	}

	return result, nil
}

/* STRING */

// http://sass-lang.com/documentation/Sass/Script/Functions.html#unquote-instance_method
func unquote(str string) (string, error) {
	if _, err := isGassStr(str); err != nil {
		return str, err
	}

	return strings.Replace(strings.Replace(str, `"`, ``, -1), `'`, ``, -1), nil
}

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
	num := 1

	if len(strings.TrimSpace(str)) > 0 {
		i, err := strconv.ParseInt(str, 10, 8)

		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%v", rand.Intn(int(i))), nil
	}

	return fmt.Sprintf("%v", rand.Intn(num)), nil
}
