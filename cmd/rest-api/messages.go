package main

import (
	"strconv"
	"strings"
	"time"

	"go-fiddle/cmd/internal/regexputil"
)

// HTTPHeader represents a http header
type HTTPHeader struct {
	Name  string `bson:"name" json:"name"`
	Value string `bson:"value" json:"value"`
}

// HTTPRequest represents http request
type HTTPRequest struct {
	Method    string        `bson:"method" json:"method"`
	URI       string        `bson:"uri" json:"uri"`
	Version   string        `bson:"version" json:"version"`
	Timestamp *time.Time    `bson:"timestamp" json:"timestamp"`
	Headers   *[]HTTPHeader `bson:"headers" json:"headers"`
	Body      string        `bson:"body" json:"body"`
}

// HTTPResponse represents http response
type HTTPResponse struct {
	StatusCode int           `bson:"statuscode" json:"statuscode"`
	StatusText string        `bson:"statustext" json:"statustext"`
	Version    string        `bson:"version" json:"version"`
	Timestamp  *time.Time    `bson:"timestamp" json:"timestamp"`
	Headers    *[]HTTPHeader `bson:"headers" json:"headers"`
	Body       string        `bson:"body" json:"body"`
}

// HTTPMessage represents a message including the request and response
type HTTPMessage struct {
	ID       string        `bson:"_id" json:"id"`
	Request  *HTTPRequest  `bson:"request" json:"request"`
	Response *HTTPResponse `bson:"response" json:"response"`
}

// HTTPMessageSummary represents minimal set of message attributes
type HTTPMessageSummary struct {
	ID         string `bson:"_id" json:"id"`
	Method     string `bson:"method" json:"method"`
	URI        string `bson:"uri" json:"uri"`
	StatusCode int    `bson:"statuscode" json:"statuscode"`
}

// UnmarshalHTTPRequest deserializes bytes to HTTPRequest
func UnmarshalHTTPRequest(data []byte) (id string, request *HTTPRequest) {
	lines := strings.Split(string(data), "\r\n")
	requestLines := lines[2:]
	requestIDMatch := regexputil.RegexMapString("^request-id: (?P<requestid>.+)$", lines[0])

	if requestIDMatch != nil {
		id = (*requestIDMatch)["requestid"]
	}

	var timestamp time.Time
	timestampMatch := regexputil.RegexMapString("^timestamp: (?P<timestamp>.+)$", lines[1])
	if timestampMatch != nil {
		timestamp, _ = time.Parse(time.RFC3339, (*timestampMatch)["timestamp"])
	}

	match := regexputil.RegexMapString("^(?P<method>[^ ]+) (?P<uri>[^ ]+) (?P<version>.+)$", requestLines[0])
	if match != nil {
		result := HTTPRequest{}
		result.Method = (*match)["method"]
		result.URI = (*match)["uri"]
		result.Version = (*match)["version"]
		result.Timestamp = &timestamp

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
	responseLines := lines[2:]
	requestIDMatch := regexputil.RegexMapString("^request-id: (?P<requestid>.+)$", lines[0])

	if requestIDMatch != nil {
		id = (*requestIDMatch)["requestid"]
	}

	var timestamp time.Time
	timestampMatch := regexputil.RegexMapString("^timestamp: (?P<timestamp>.+)$", lines[1])
	if timestampMatch != nil {
		timestamp, _ = time.Parse(time.RFC3339, (*timestampMatch)["timestamp"])
	}

	match := regexputil.RegexMapString("^(?P<version>[^ ]+) (?P<statuscode>[^ ]+) (?P<status>.+)$", responseLines[0])
	if match != nil {
		result := HTTPResponse{}
		statusCode, _ := strconv.ParseInt((*match)["statuscode"], 10, 32)
		result.StatusCode = int(statusCode)
		result.StatusText = (*match)["status"]
		result.Version = (*match)["version"]
		result.Timestamp = &timestamp

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
