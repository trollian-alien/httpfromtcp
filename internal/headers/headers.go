package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
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
	key, value, found := strings.Cut(keyAndValue, ":")
	if found == false {
		return 0, false, fmt.Errorf("invalid header syntax: colon missing")
	} else if strings.HasSuffix(key, " ") {
		return 0, false, fmt.Errorf("invalid header field-name %v", key)
	}
	h[strings.TrimSpace(key)] = strings.TrimSpace(value)
	return index + 2, false, nil
}