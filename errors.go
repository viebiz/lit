package lightning

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

type HttpError struct {
	Status      int    `json:"-"`
	Code        string `json:"error"`
	Description string `json:"error_description"`
}

func (e HttpError) Error() string {
	return fmt.Sprintf(`{"error":"%s","error_description":"%s"}`, e.Code, e.Description)
}
