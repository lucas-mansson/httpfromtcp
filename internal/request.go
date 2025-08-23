package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"slices"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	Method string
	RequestTarget string
	HttpVersion string
}

func RequestFromReader(reader io.Reader)(*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	requestLine, err := parseRequestLine(string(data))
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	request := Request { RequestLine: requestLine }

	return &request, nil
}

func parseRequestLine(request string)(RequestLine, error) {
	// split request
	unparsedLine := strings.Trim(strings.Split(request, "\r\n")[0], " ")
	
	// [method, target, version]
	splitLine := strings.Split(unparsedLine, " ")
	if len(splitLine) != 3 {
		return RequestLine{}, errors.New("Error parsing request line.")
	}

	// method
	method, err := checkMethod(splitLine[0])
	if err != nil {
		return RequestLine{}, err
	}

	// correct target
	target, err := checkTarget(splitLine[1])
	if err != nil {
		return RequestLine{}, err
	}

	// correct version
	version, err := checkVersion(splitLine[2])
	if err != nil {
		return RequestLine{}, err
	}

	return RequestLine { HttpVersion: version, RequestTarget: target, Method: method }, nil
}

func checkMethod(method string)(string, error) {
	allowedMethods := []string{"GET", "POST"}
	if !slices.Contains(allowedMethods, method) {
		return "", errors.New("Error parsing HTTP method.")
	}
	return method, nil
}

func checkTarget(target string)(string, error) {
	if target[0] != '/' {
		return "", errors.New("HTTP target path doesn't start with '/'!")
	}
	return target, nil
}

func checkVersion(version string)(string, error) {
	splitVersion := strings.Split(version, "/")
	if splitVersion[0] != "HTTP" {
		return "", errors.New("Error parsing HTTP version")
	}
	httpVersion := splitVersion[1]
	if httpVersion != "1.1" {
		return "", errors.New("Unsupported HTTP version")
	}
	
	return httpVersion, nil
}
