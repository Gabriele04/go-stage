package http

import (
	"mysql/app/apperr"
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	GenericErrorMessage = "An error occurred"
)

// codes represents an HTTP status code.
var codes = map[string]int{
	apperr.ECONFLICT:       http.StatusConflict,
	apperr.EFORBIDDEN:      http.StatusForbidden,
	apperr.EINVALID:        http.StatusBadRequest,
	apperr.ENOTFOUND:       http.StatusNotFound,
	apperr.ENOTIMPLEMENTED: http.StatusNotImplemented,
	apperr.EUNAUTHORIZED:   http.StatusUnauthorized,
	apperr.EINTERNAL:       http.StatusInternalServerError,
	apperr.EUNKNOWN:        http.StatusInternalServerError,
}

// MessageFromErr returns the message for the given app error.
// EINTERNAL & EUNKNOWN message is obscured by HTTP response.
func MessageFromErr(err error) string {

	appErrMessage := apperr.ErrorMessage(err)
	appErrCode := apperr.ErrorCode(err)

	if appErrMessage == "" || appErrCode == apperr.EINTERNAL || appErrCode == apperr.EUNKNOWN {
		return GenericErrorMessage
	}

	return appErrMessage
}

// StatusCodeFromErr returns the HTTP status code for the given app error.
func StatusCodeFromErr(err error) int {

	appErrCode := apperr.ErrorCode(err)

	code, ok := codes[appErrCode]
	if !ok {
		code = http.StatusInternalServerError
	}
	return code
}

// ErrorAPI represents an error returned by the API.
type ErrorAPI struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

// NewErrorAPI returns an ErrorAPI instance.
func NewErrorAPI(err error, details interface{}) *ErrorAPI {
	e := &ErrorAPI{
		Code:    apperr.ErrorCode(err),
		Message: MessageFromErr(err),
		Details: details,
	}
	return e
}

// ErrorResponseJSON returns an HTTP error response with JSON content.
func ErrorResponseJSON(c echo.Context, err error, details interface{}) error {
	return c.JSON(StatusCodeFromErr(err), NewErrorAPI(err, details))
}