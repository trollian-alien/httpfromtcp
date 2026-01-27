package response

import (
	"fmt"

	"github.com/trollian-alien/httpfromtcp/internal/headers"
)

//Me: Go, I want enums
//Go: we already gave enums at home
//Enums at home:
type StatusCode int

const (
	StatusCodeSuccess             StatusCode = 200
	StatusCodeBadRequest          StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

func getStatusLine (statusCode StatusCode) []byte {
	reasonPhrase := "" //we set the text in the switch statement below
	switch statusCode {
	case 200:
		reasonPhrase = "OK"
	case 400:
		reasonPhrase = "Bad Request"
	case 500:
		reasonPhrase = "Internal Server Error"
	}
	return []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase))
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}