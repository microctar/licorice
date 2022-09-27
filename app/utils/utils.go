package utils

import (
	"regexp"
)

// Regexp utilities

var GetProt = Get

func Get(pattern, text string) string {
	r := regexp.MustCompile(pattern)

	if result := r.FindStringSubmatch(text); result != nil {

		// e.g. result => ["matched string", "substring"]

		return result[1]
	}

	return ""
}

func GetByParam(pattern, text string) string {
	r := regexp.MustCompile(pattern)

	if result := r.FindStringSubmatch(text); result != nil {

		// e.g. result => ["matched string", "substring", "", ...]

		for _, param := range result[1:] {
			if param != "" {
				return param
			}
		}
	}

	return ""
}

// convert string to boolean

func Str2Bool(boolean_val string) bool {
	if boolean_val == "true" {
		return true
	} else {
		return false
	}
}
