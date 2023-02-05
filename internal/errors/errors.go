package errors

import (
	"strings"
)

type ErrorType string

const (
	ErrorTypeJSONMarshal   ErrorType = "JSONMarshal"
	ErrorTypeJSONUnmarshal ErrorType = "JSONUnmarshal"
	ErrorTypeJWTExpired    ErrorType = "JWTExpired"
	ErrorTypeTokenRefresh  ErrorType = "TokenRefresh"
)

type ErrorSource string

const (
	ErrorSourceGraphQL ErrorSource = "graphql"
	ErrorSourceNhost   ErrorSource = "nhost"
	ErrorSourceGo      ErrorSource = "go"
)

type Error struct {
	Error   error
	Message string
	Type    ErrorType
	Source  ErrorSource
}

func New(err error, message string, typ ErrorType, source ErrorSource) *Error {
	return &Error{
		Error:   err,
		Message: message,
		Type:    typ,
		Source:  source,
	}
}

func Parse(err error) *Error {
	payload := strings.Split(err.Error(), ":")
	var response Error
	response.Error = err

	switch strings.TrimSpace(payload[0]) {
	case "graphql":
		response.Source = ErrorSourceGraphQL
	}

	response.Message = strings.TrimSpace(payload[1])

	if len(payload) > 2 {
		switch strings.TrimSpace(payload[2]) {
		case "JWTExpired":
			response.Type = ErrorTypeJWTExpired
		}
	}

	return &response
}

func (e *Error) IsType(errType ErrorType) bool {
	return e.Type == errType
}
