package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"slices"
)

// GET /path HTTP/1.1 
type RequestLine struct {
	Method string
	RequestTarget string
	HttpVersion string
}

type Request struct {
	RequestLine RequestLine
    // Headers     map[string]string
   //  Body        []byte
}

func RequestFromReader(reader io.Reader)(*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("Error reading with io.ReadAll"), err)
	}

	requestLine, err := parseRequestLine(string(data))
	if err != nil {
		return nil, err
	}
	return &Request { 
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(request string)(*RequestLine, error) {
	// split request
	unparsedLine := strings.Trim(strings.Split(request, "\r\n")[0], " ")
	
	// Get array [method, target, version]
	splitLine := strings.Split(unparsedLine, " ")
	filterEmptyStrings(&splitLine); //Remove any extra whitespace
	if len(splitLine) != 3 {
		return &RequestLine{}, fmt.Errorf("Malformed request line")
	}

	// method
	method, err := checkMethod(splitLine[0])
	if err != nil {
		return &RequestLine{}, err
	}

	// correct target
	target, err := checkTarget(splitLine[1])
	if err != nil {
		return &RequestLine{}, err
	}

	// correct version
	version, err := checkVersion(splitLine[2])
	if err != nil {
		return &RequestLine{}, err
	}

	return &RequestLine { 
		HttpVersion: version, 
		RequestTarget: target, 
		Method: method,
	}, nil
}

func filterEmptyStrings(arr *[]string) {
    filtered := make([]string, 0, len(*arr))
    for _, val := range *arr {
        if val != "" {
            filtered = append(filtered, val)
        }
    }
    *arr = filtered
}

func checkMethod(method string)(string, error) {
	allowedMethods := []string{"GET", "POST"}
	if !slices.Contains(allowedMethods, method) {
		return "", fmt.Errorf("Unallowed HTTP method.")
	}
	return method, nil
}

func checkTarget(target string)(string, error) {
	if target[0] != '/' {
		return "", fmt.Errorf("HTTP target path doesn't start with '/'!")
	}
	return target, nil
}

func checkVersion(version string)(string, error) {
	splitVersion := strings.Split(version, "/")
	if len(splitVersion) != 2 || splitVersion[0] != "HTTP" {
		return "", fmt.Errorf("Error parsing HTTP version")
	}
	httpVersion := splitVersion[1]
	if httpVersion != "1.1" {
		return "", fmt.Errorf("Unsupported HTTP version")
	}
	
	return httpVersion, nil
}
