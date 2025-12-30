package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
	"github.com/trollian-alien/httpfromtcp/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Headers headers.Headers
	state requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

//Me: Go, I want enums 
//Go: we already gave enums at home
//Enums at home:
type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
	requestStateDone
)

const readSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{
		Headers: headers.NewHeaders(),
		state: requestStateInitialized,
	}
	readToIndex := 0
	b := make([]byte, readSize)
	for request.state != requestStateDone {

		if readToIndex >= len(b) {
		newB := make([]byte, len(b) *2)
		copy(newB, b)
		b = newB
		}
		numBytesRead, err := reader.Read(b[readToIndex:])
		if err == io.EOF {
			// final parse for the remianing bytes, if there are any
			if readToIndex > 0 {
            _, perr := request.parse(b[:readToIndex])
				if perr != nil {
					return nil, perr
				}
        	}

			// incomplete request
			if request.state != requestStateDone {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", request.state, numBytesRead)
			}
			break
		} else if err != nil {
			return nil, fmt.Errorf("error reading: %v", err)
		}
		readToIndex += numBytesRead
		numBytesParsed, err := request.parse(b[:readToIndex])
		if err != nil {
			return nil, err
		}
		if numBytesParsed > readToIndex {
    		return nil, fmt.Errorf("parser bug: parsed %d bytes from buffer of size %d", numBytesParsed, readToIndex)
		}

		copy(b, b[numBytesParsed:])
		readToIndex -= numBytesParsed
	}

	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			// something went wrong
			return 0, err
		}
		if n == 0 {
			// need more data
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.state = requestStateParsingHeaders
		return n, nil
	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.state = requestStateDone
		}
		return n, nil
	case requestStateDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}

func parseRequestLine(bytes []byte) (*RequestLine, int, error) {
	requestLine := string(bytes) // we get rid of the extra stuff now
	index := strings.Index(requestLine, "\r\n")
	if index == -1 {
		// no CRLF, need more data
		return nil, 0, nil
	}

	requestLine = requestLine[:index]
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return nil, 0, fmt.Errorf("invalid request line")
	}

	if parts[2] != "HTTP/1.1" {
		return nil, 0, fmt.Errorf("http version %s not supported", parts[2])
	}
	method := parts[0]
	if method == "" {
		return nil, 0, fmt.Errorf("http method can't be empty!")
	}
	for _, m := range method {
		if !(unicode.IsLetter(m) && unicode.IsUpper(m)) {
			return nil, 0, fmt.Errorf("http methods can only have uppercase alphabetic characters")
		}
	}
	temp := strings.Split(parts[2], "/")
	version := temp[1]
	r := RequestLine{
		HttpVersion: version,
		RequestTarget: parts[1],
		Method: method,
	}

	return &r, index +2, nil
}