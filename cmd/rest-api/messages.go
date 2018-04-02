package main

import (
	"strconv"
	"strings"

	"go-fiddle/internal/regexputil"
)

// HTTPHeader represents a http header
type HTTPHeader struct {
	Name  string
	Value string
}

// HTTPRequest represents http request
type HTTPRequest struct {
	Method  string
	URI     string
	Version string
	Headers *[]HTTPHeader
	Body    string
}

// HTTPResponse represents http response
type HTTPResponse struct {
	StatusCode int
	Version    string
	Headers    *[]HTTPHeader
	Body       string
}

// HTTPMessage represents a message including the request and response
type HTTPMessage struct {
	ID       string `bson:"_id" json:"_id"`
	Request  *HTTPRequest
	Response *HTTPResponse
}

// UnmarshalHTTPRequest deserializes bytes to HTTPRequest
func UnmarshalHTTPRequest(data []byte) (id string, request *HTTPRequest) {
	lines := strings.Split(string(data), "\r\n")
	requestLines := lines[1:]
	match := regexputil.RegexMapString("^request-id: (?P<requestid>.+)$", lines[0])

	if match != nil {
		id = (*match)["requestid"]
	}

	match = regexputil.RegexMapString("^(?P<method>[^ ]+) (?P<uri>[^ ]+) (?P<version>.+)$", requestLines[0])
	if match != nil {
		result := HTTPRequest{}
		result.Method = (*match)["method"]
		result.URI = (*match)["uri"]
		result.Version = (*match)["version"]

		headers := []HTTPHeader{}

		for i, line := range requestLines[1:] {
			if line == "" {
				result.Body = strings.Join(requestLines[i+2:], "\r\n")
				break
			}
			match = regexputil.RegexMapString("^(?P<name>[^:]+): (?P<value>.+)$", line)
			if match != nil {
				headers = append(headers, HTTPHeader{(*match)["name"], (*match)["value"]})
			}
		}

		result.Headers = &headers
		request = &result
	}

	return
}

// UnmarshalHTTPResponse deserializes bytes to HTTPRequest
func UnmarshalHTTPResponse(data []byte) (id string, response *HTTPResponse) {
	lines := strings.Split(string(data), "\r\n")
	responseLines := lines[1:]
	match := regexputil.RegexMapString("^request-id: (?P<requestid>.+)$", lines[0])

	if match != nil {
		id = (*match)["requestid"]
	}

	match = regexputil.RegexMapString("^(?P<version>[^ ]+) (?P<statuscode>[^ ]+) (?P<status>.+)$", responseLines[0])
	if match != nil {
		result := HTTPResponse{}
		statusCode, _ := strconv.ParseInt((*match)["statuscode"], 10, 32)
		result.StatusCode = int(statusCode)
		result.Version = (*match)["version"]

		headers := []HTTPHeader{}

		for i, line := range responseLines[1:] {
			if line == "" {
				result.Body = strings.Join(responseLines[i+2:], "\r\n")
				break
			}
			match = regexputil.RegexMapString("^(?P<name>[^:]+): (?P<value>.+)$", line)
			if match != nil {
				headers = append(headers, HTTPHeader{(*match)["name"], (*match)["value"]})
			}
		}

		result.Headers = &headers
		response = &result
	}

	return
}
