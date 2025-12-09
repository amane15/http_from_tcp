package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	httpMessage, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	requestLine, err := parseRequestLine(httpMessage)
	if err != nil {
		return nil, err
	}

	request := Request{
		RequestLine: *requestLine,
	}

	return &request, nil
}

func parseRequestLine(httpMessage []byte) (*RequestLine, error) {
	line := strings.Split(string(httpMessage), "\r\n")[0]

	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		return nil, errors.New("invalid request line")
	}

	httpMethod := parts[0]
	httpTarget := parts[1]
	httpVersion := strings.Split(parts[2], "/")[1]

	if httpMethod != strings.ToUpper(httpMethod) {
		return nil, fmt.Errorf("invalid method name: %s", httpMethod)
	}
	if httpVersion != "1.1" {
		return nil, fmt.Errorf("version %s is not supported", httpVersion)
	}

	requestLine := RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: httpTarget,
		Method:        httpMethod,
	}

	return &requestLine, nil
}
