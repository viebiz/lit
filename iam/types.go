package iam

import (
	"net/http"
)

// HTTPClient defines an interface for making HTTP requests,
// allowing for easier unit testing by mocking http.Client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ExpressionFunction defines a custom function for Casbin models, compatible with govaluate.ExpressionFunction.
type ExpressionFunction func(arguments ...interface{}) (interface{}, error)
