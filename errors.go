package lit

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// HTTP errors
	ErrInternalServerError = HttpError{Status: http.StatusInternalServerError, Code: "internal_server_error", Description: "internal server error"}

	// gRPC errors
	ErrGRPCInternalServerError = status.Errorf(codes.Internal, "internal server error")
)

// ExpectedError represents a known error that should be returned to the client.
// It helps the framework distinguish expected errors from unexpected ones for monitoring purposes.
type ExpectedError interface {
	error

	StatusCode() int

	ErrorCode() string
}

// HttpError represents an expected error from HTTP request
type HttpError struct {
	Status      int    `json:"-"`
	Code        string `json:"error"`
	Description string `json:"error_description"`
}

func (e HttpError) Error() string {
	return fmt.Sprintf(`{"error":"%s","error_description":"%s"}`, e.Code, e.Description)
}

func (e HttpError) StatusCode() int {
	return e.Status
}

func (e HttpError) ErrorCode() string {
	return e.Code
}
