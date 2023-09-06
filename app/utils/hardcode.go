package utils

import "regexp"

// Regexp utilities

func ReGetFirst(re *regexp.Regexp, text string) string {

	if result := re.FindStringSubmatch(text); result != nil {

		// e.g. result => ["matched string", "substring"]

		return result[1]
	}

	return ""
}

// convert string to boolean

func Str2Bool(valBool string) bool {
	return valBool == "true"
}
