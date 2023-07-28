package apperror

import (
	"fmt"
	"net/http"
)

const name = "PS"

var (
	ErrInternalSystem = NewAppError("00100", http.StatusInternalServerError, "internal system error")
	ErrBadRequest     = NewAppError("00101", http.StatusBadRequest, "bad request")
	ErrValidation     = NewAppError("00102", http.StatusBadRequest, "validation error")
	ErrNotFound       = NewAppError("00103", http.StatusNotFound, "not found")
	ErrUnauthorized   = NewAppError("00104", http.StatusUnauthorized, "Not Authorized")
	ErrForbidden      = NewAppError("00105", http.StatusForbidden, "access forbidden")
)

type AppError struct {
	Err           error  `json:"-"`
	Message       string `json:"message,omitempty"`
	Code          string `json:"-"`
	TransportCode int    `json:"code"`
}

func NewAppError(code string, transportCode int, message string) *AppError {
	return &AppError{
		Err:           fmt.Errorf(message),
		Code:          name + "-" + code,
		TransportCode: transportCode,
		Message:       message,
	}
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func ErrorWithMessage(err *AppError, message string) *AppError {
	if message == "" {
		return err
	} else {
		err.Message = message
		return err
	}
}
