package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var res strings.Builder
	var buf rune

	if str == "" {
		return "", nil
	}

	for _, char := range str {
		curIsDigit := unicode.IsDigit(char)

		if buf == 0 && curIsDigit {
			return "", ErrInvalidString
		}

		if curIsDigit {
			cnt, err := strconv.Atoi(string(char))
			if err != nil {
				return "", ErrInvalidString
			}
			if cnt > 0 {
				res.WriteString(strings.Repeat(string(buf), cnt-1))
			} else {
				splitArr := []rune(res.String())
				splitArr = splitArr[:len(splitArr)-1]
				res.Reset()
				res.WriteString(string(splitArr))
			}
			buf = 0
		} else {
			buf = char
			res.WriteRune(char)
		}
	}

	return res.String(), nil
}
