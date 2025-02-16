package lit

import (
	"fmt"
	"net/http"
)

// Error represents standard error of lit framework
type Error interface {
	error

	StatusCode() int // Suppose return status code

	ErrorCode() string // Suppose support internalization
}

var (
	ErrDefaultInternal = HttpError{Status: http.StatusInternalServerError, Code: "internal_server_error", Desc: "Something went wrong"}
)

// HttpError represents an expected error from HTTP request
type HttpError struct {
	Status int    `json:"-"`
	Code   string `json:"error"`
	Desc   string `json:"error_description"`
}

func (e HttpError) StatusCode() int {
	return e.Status
}

func (e HttpError) ErrorCode() string {
	return e.Code
}

func (e HttpError) Error() string {
	return fmt.Sprintf("Status: [%d], Code: [%s], Desc: [%s]", e.Status, e.Code, e.Desc)
}
