package hw02unpackstring

import (
	"errors"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	sr := []rune(s)
	var isSlash bool
	res := []rune{}

	if len(s) == 0 {
		return "", nil
	}
	if s == `\` || unicode.IsDigit(sr[0]) {
		return "", ErrInvalidString
	}

	for i, item := range sr {
		if len(sr) == i+1 && item == '\\' && sr[i-1] != '\\' ||
			unicode.IsDigit(item) && unicode.IsDigit(sr[i-1]) && sr[i-2] != '\\' {
			return "", ErrInvalidString
		}

		if item == '\\' && !isSlash {
			isSlash = true
			continue
		}
		if isSlash && unicode.IsLetter(item) {
			return "", ErrInvalidString
		}
		if isSlash {
			res = append(res, item)
			isSlash = false
			continue
		}
		if unicode.IsDigit(item) {
			n, _ := strconv.Atoi(string(item))
			if n == 0 {
				res = res[:len(res)-1]
				continue
			}
			res = MulRune(res, sr[i-1], n)
			continue
		}
		res = append(res, item)
	}

	return string(res), nil
}

func MulRune(res []rune, l rune, n int) []rune {
	for i := 1; i < n; i++ {
		res = append(res, l)
	}

	return res
}
