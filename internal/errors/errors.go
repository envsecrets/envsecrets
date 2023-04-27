package errors

import (
	"net/http"
	"strings"

	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

type ErrorType string

const (
	ErrorTypeBase64Encode ErrorType = "Base64Encode"
	ErrorTypeBase64Decode ErrorType = "Base64Decode"

	ErrorTypeJSONMarshal      ErrorType = "JSONMarshal"
	ErrorTypeJSONUnmarshal    ErrorType = "JSONUnmarshal"
	ErrorTypeJWTExpired       ErrorType = "JWTExpired"
	ErrorTypeUnauthorized     ErrorType = "Unauthorized"
	ErrorTypePermissionDenied ErrorType = "PermissionDenied"
	ErrorTypeTokenRefresh     ErrorType = "TokenRefresh"

	ErrorTypeInvalidResponse ErrorType = "InvalidResponse"
	ErrorTypeBadResponse     ErrorType = "BadResponse"
	ErrorTypeBadRequest      ErrorType = "BadRequest"
	ErrorTypeBadGateway      ErrorType = "BadGateway"
	ErrorTypeRequestFailed   ErrorType = "RequestFailed"

	ErrorTypeDoesNotExist ErrorType = "DoesNotExist"
	ErrorTypeInvalidKey   ErrorType = "InvalidKey"
	ErrorTypeKeyNotFound  ErrorType = "KeyNotFound"

	ErrorTypeInvalidToken ErrorType = "InvalidToken"

	ErrorTypeInvalidAccountConfiguration ErrorType = "InvalidAccountConfiguration"
	ErrorTypeInvalidProjectConfiguration ErrorType = "InvalidProjectConfiguration"

	ErrorTypeEmailFailed ErrorType = "EmailFailed"
)

var ResponseCodeMap = map[ErrorType]int{
	ErrorTypeJSONMarshal:      http.StatusBadRequest,
	ErrorTypeJSONUnmarshal:    http.StatusBadRequest,
	ErrorTypeJWTExpired:       http.StatusUnauthorized,
	ErrorTypeUnauthorized:     http.StatusUnauthorized,
	ErrorTypePermissionDenied: http.StatusForbidden,
	ErrorTypeTokenRefresh:     http.StatusConflict,

	ErrorTypeInvalidResponse: http.StatusInternalServerError,
	ErrorTypeBadResponse:     http.StatusInternalServerError,
	ErrorTypeBadRequest:      http.StatusBadRequest,
	ErrorTypeBadGateway:      http.StatusBadGateway,
	ErrorTypeRequestFailed:   http.StatusBadRequest,

	ErrorTypeDoesNotExist:                http.StatusNotFound,
	ErrorTypeInvalidKey:                  http.StatusBadRequest,
	ErrorTypeKeyNotFound:                 http.StatusNotFound,
	ErrorTypeInvalidAccountConfiguration: http.StatusBadRequest,
	ErrorTypeInvalidProjectConfiguration: http.StatusBadRequest,

	ErrorTypeInvalidToken: http.StatusBadRequest,

	ErrorTypeEmailFailed: http.StatusInternalServerError,
}

func (e *ErrorType) GetStatusCode() int {
	for key, value := range ResponseCodeMap {
		if &key == e {
			return value
		}
	}

	return http.StatusBadRequest
}

type ErrorSource string

const (
	ErrorSourceHTTP    ErrorSource = "http"
	ErrorSourceGraphQL ErrorSource = "graphql"
	ErrorSourceHermes  ErrorSource = "hermes"
	ErrorSourceMailer  ErrorSource = "mailer"

	ErrorSourceVault ErrorSource = "vault"
	ErrorSourceNhost ErrorSource = "nhost"
	ErrorSourceGo    ErrorSource = "go"

	ErrorSourceGithub ErrorSource = "github"
	ErrorSourceVercel ErrorSource = "vercel"
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

	if strings.Contains(response.Message, "permission has failed") {
		response.Type = ErrorTypePermissionDenied
		response.Message = response.GenerateMessage("Permission Denied")
	}

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

func (e *Error) GenerateMessage(defaultMessage string) string {

	switch e.Type {
	case ErrorTypePermissionDenied:
		return "You do not have permissions to perform this action"
	default:
		return defaultMessage
	}
}

func (e *Error) Log(logger *logrus.Logger, defaultMessage string) {

	level := logger.GetLevel()

	if level == logrus.DebugLevel {
		log.Debug(e.Error)
	}

	switch e.Type {
	case ErrorTypeUnauthorized:
		logger.Error("you are not authorized to perform this action")
	}

	if defaultMessage != "" {
		if level == logrus.ErrorLevel {
			logger.Error(defaultMessage)
		} else {
			logger.Fatal(defaultMessage)
		}
	}
}
