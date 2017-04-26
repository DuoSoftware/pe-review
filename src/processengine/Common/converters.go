package Common

import (
	"bytes"
	"regexp"
	"strings"
)

// initializing the Regex statement
var camelingRegex = regexp.MustCompile("[0-9A-Za-z]+")

// the method which converts the string to camelcasing
func CamelCase(src string) string {
	// convert the string to bytes
	byteSrc := []byte(src)
	// convert it to camelcasing
	chunks := camelingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		if idx > 0 {
			chunks[idx] = bytes.Title(val)
		}
	}
	// return after string is totally converted
	return string(bytes.Join(chunks, nil))
}

func MakeFirstLowerCase(s string) string {

	if len(s) < 2 {
		return strings.ToLower(s)
	}

	bts := []byte(s)

	lc := bytes.ToLower([]byte{bts[0]})
	rest := bts[1:]

	return string(bytes.Join([][]byte{lc, rest}, nil))
}
