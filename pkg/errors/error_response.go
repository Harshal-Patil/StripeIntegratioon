package errors

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"
)

// ErrorResponse ...
type ErrorResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// Error ...
type Error struct {
	ErrorO ErrorResponse `json:"error"`
}

// Error ...
func (e Error) Error() string {
	return e.ErrorO.Message
}

// InternalServerError creates a new error response representing an internal server error (HTTP 500)
func InternalServerError(msg string) Error {
	if msg == "" {
		msg = "We encountered an error while processing your request."
	}
	return Error{
		ErrorO: ErrorResponse{
			Type:    "Internal Server Error",
			Message: msg,
		},
	}
}

// NotFound creates a new error response representing a resource-not-found error (HTTP 404)
func NotFound(msg string) Error {
	if msg == "" {
		msg = "The requested resource was not found."
	}
	return Error{
		ErrorO: ErrorResponse{
			Type:    "Not Found",
			Message: msg,
		},
	}
}

// BadRequest creates a new error response representing a bad request (HTTP 400)
func BadRequest(msg string) Error {
	if msg == "" {
		msg = "Your request is in a bad format."
	}
	return Error{
		ErrorO: ErrorResponse{
			Type:    "Bad Request",
			Message: msg,
		},
	}
}

// InvalidInput creates a new error response representing a data validation error (HTTP 400).
func InvalidInput(errs error) ErrorResponse {
	return ErrorResponse{
		Type:    "Bad Request",
		Message: "There is some problem with the data you submitted.",
	}
}

// ParseValidatorError ....
func ParseValidatorError(errorStr []string) Error {
	return Error{
		ErrorO: ErrorResponse{
			Type:    "Bad Request",
			Message: "There is some problem with the data you submitted.",
		},
	}
}

// Forbidden ....
func Forbidden(errorStr string) Error {
	return Error{
		ErrorO: ErrorResponse{
			Type:    "Unauthorized",
			Message: errorStr,
		},
	}
}

// ParseError ...
func ParseError(err error) (Error, int) {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return NotFound(""), http.StatusNotFound
	case errors.Is(err, context.DeadlineExceeded):
	case strings.Contains(err.Error(), "Error:Field validation"):
		splitedError := splitErrorMessage(err)
		return ParseValidatorError(splitedError), http.StatusBadRequest
	case strings.Contains(err.Error(), "cannot unmarshal"):
		return BadRequest(err.Error()), http.StatusBadRequest
	case strings.Contains(err.Error(), "Unmarshal"):
		return BadRequest(""), http.StatusBadRequest
	case strings.Contains(err.Error(), "invalid UUID length"):
		return BadRequest("invalid uuid length"), http.StatusBadRequest
	case strings.Contains(err.Error(), "strconv.ParseBool: parsing"):
		return BadRequest("invalid boolean flag"), http.StatusBadRequest
	case strings.Contains(err.Error(), "request Content-Type isn't multipart/form-data"):
		return BadRequest("Invalid file format"), http.StatusBadRequest
	case strings.Contains(err.Error(), "header should contains authorization field with valid jwt"):
		return Forbidden("Invalid file format"), http.StatusForbidden
	default:
		return InternalServerError(err.Error()), http.StatusInternalServerError
	}
	return InternalServerError(err.Error()), http.StatusInternalServerError
}

func splitErrorMessage(err error) []string {
	splitedError := make([]string, 0)
	for i := 0; i < len(err.Error()); i++ {
		if err.Error()[i] == 39 {
			var subStr string
			var j int
			for j = i + 1; err.Error()[j] != 39; j++ {
				subStr = subStr + string(err.Error()[j])
			}
			splitedError = append(splitedError, subStr)
			i = j + 1
		}
	}
	return splitedError
}
