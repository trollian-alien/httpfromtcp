package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func isTokenChar(r rune) bool {
    if r >= 'a' && r <= 'z' { return true }
    if r >= '0' && r <= '9' { return true }
    switch r {
    case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
        return true
    }
    return false
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	str := string(data)
	index := strings.Index(str, "\r\n")
	switch index {
	case -1:
		return 0, false, nil
	case 0:
		return 2, true, nil
	}
	
	keyAndValue := str[:index]
	//extracting the key and value and checking syntax
	key, value, found := strings.Cut(keyAndValue, ":")
	if found == false {
		return 0, false, fmt.Errorf("invalid header syntax: colon missing")
	} else if strings.HasSuffix(key, " ") {
		return 0, false, fmt.Errorf("invalid header field-name %v", key)
	}
	key = strings.ToLower(strings.TrimSpace(key))
	for _, r := range key {
		if !isTokenChar(r) {
			return 0, false, fmt.Errorf("invalid header key: %v", key)
		}
	}
	value = strings.TrimSpace(value) 

	//adding the value to the key map
	oldVal, ok := h[key]
	if !ok {
		h[key] = strings.TrimSpace(value)
	} else {
		h[key] = oldVal + ", " + value
	}

	return index + 2, false, nil
}