package optstr

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlphabet(c rune) bool {
	return (c >= 'a' && c <= 'z') || c >= 'A' && c <= 'Z'
}

func isValidKeyStart(c rune) bool {
	return isAlphabet(c) || c == '_'
}

func isValidKeyPart(c rune) bool {
	return c == '-' || isValidKeyStart(c) || isDigit(c)
}

func isValidRawStart(c rune) bool {
	return c > '\x20' && c != '\x7f' && c != '"' && c != ' '
}

func isValidRawPart(c rune) bool {
	return isValidRawStart(c) || c == '"'
}

//go:generate stringer -linecomment -type parseError -output parse_error_string.go
type parseError int

const (
	_                  parseError = iota
	errBadUtf8Encoding            // bad utf-8 encoding
	errExpectKey                  // expect key
	errBadKey                     // bad key
	errBadRawValue                // bad raw value
	errExpectValue                // expect value
	errExpectEQ                   // expect '='
)

func (i parseError) Error() string {
	return i.String()
}

func readUtf8(s string) (ch rune, n int, e error) {
	ch, n = utf8.DecodeRuneInString(s)
	if n == 1 && ch == utf8.RuneError {
		e = errBadUtf8Encoding
	}
	return
}

func readKey(input string) (key string, n int, e error) {
	ch, cn, e := readUtf8(input)
	if e != nil {
		return
	}
	if cn == 0 || !isValidKeyStart(ch) {
		e = errExpectKey
		return
	}
	i := cn
	for i < len(input) {
		var ch rune
		var cn int
		ch, cn, e = readUtf8(input[i:])
		if e != nil {
			return
		}
		if ch == '=' {
			break
		}
		i += cn
		if cn == 0 || !isValidKeyPart(ch) {
			e = errBadKey
			return
		}
	}
	return input[:i], i, nil
}

func readValue(input string) (value string, n int, e error) {
	var ch rune
	ch, _, e = readUtf8(input)
	if e != nil {
		return
	}
	if ch == '"' {
		return readQuotedValue(input)
	}
	return readRawValue(input)
}

func readRawValue(input string) (value string, n int, e error) {
	i := 0
	for i < len(input) {
		var ch rune
		var cn int
		ch, cn, e = readUtf8(input[i:])
		if e != nil {
			return
		}
		if ch == ' ' {
			break
		}
		if cn == 0 || !isValidRawPart(ch) {
			e = errBadRawValue
			return
		}
		i += cn
	}
	return input[:i], i, nil
}

func readQuotedValue(input string) (value string, n int, e error) {
	var qs string
	qs, e = strconv.QuotedPrefix(input)
	if e != nil {
		return
	}
	value, e = strconv.Unquote(qs)
	if e != nil {
		return
	}
	n = len(qs)
	return
}

type Option struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func ParseString(input string) (rs []Option, e error) {
	input = strings.TrimSpace(input)
	rs = make([]Option, 0)
	var i = 0
	defer func() {
		if e != nil {
			e = fmt.Errorf("parse options failed: %w", e)
		}
	}()

	for i < len(input) {
		var (
			key   string
			value string
			ch    rune
			n     int
		)
		ch, n, e = readUtf8(input[i:])
		if e != nil {
			return
		}
		if ch == ' ' {
			i += n
			continue
		}
		key, n, e = readKey(input[i:])
		if e != nil {
			return
		}
		i += n
		ch, n, e = readUtf8(input[i:])
		if e != nil {
			return
		}
		if ch != '=' {
			e = errExpectEQ
			return
		}
		i += n
		value, n, e = readValue(input[i:])
		if e != nil {
			return
		}
		i += n
		rs = append(rs, Option{Key: key, Value: value})
	}
	return
}
