package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 {
		// empty line
		// headers are done, consume the CRLF
		return 2, true, nil
	}

	parts := bytes.SplitN(data[:idx], []byte(":"), 2)
	key := strings.ToLower(string(parts[0]))

	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	value := string(bytes.TrimSpace(parts[1]))
	key = strings.TrimSpace(key)
	isValidKey := isValidHeaderKey(key)
	if !isValidKey {
		return 0, false, fmt.Errorf("invalid character in header key: %s", key)
	}

	if oldValue, found := h[key]; found {
		value = strings.Join([]string{oldValue, value}, ", ")
		fmt.Println("==========")
		fmt.Println(oldValue)
	}

	h.Set(strings.ToLower(key), value)
	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	h[strings.ToLower(key)] = value
}

func (h Headers) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	value, ok := h[key]
	return value, ok
}

func (h Headers) Override(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func (h Headers) Remove(key string) {
	key = strings.ToLower(key)
	delete(h, key)
}

func isValidHeaderKey(key string) bool {
	for _, c := range key {
		if c > 127 {
			return false
		}

		switch {
		case 'A' <= c && c <= 'Z':
		case 'a' <= c && c <= 'z':
		case '0' <= c && c <= '9':
		case c == '!' || c == '#' || c == '$' || c == '%' || c == '&' ||
			c == '\'' || c == '*' || c == '+' || c == '-' || c == '.' ||
			c == '^' || c == '_' || c == '`' || c == '|' || c == '~':
		default:
			return false
		}
	}

	return len(key) > 0
}
