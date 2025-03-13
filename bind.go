package lit

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/viebiz/lit/i18n"
)

// Bind binds the incoming request body and URI parameters to the provided object.
// Return error if the got error when binding and validating the object.
// For more about validation tags refer to: `https://pkg.go.dev/github.com/go-playground/validator/v10#hdr-Baked_In_Validators_and_Tags`
func (c litContext) Bind(obj interface{}) error {
	// Read more at `https://gin-gonic.com/docs/examples/binding-and-validation`
	if err := c.Context.ShouldBind(obj); err != nil {
		return convertValidationErr(c, err)
	}

	// Default binding does not include URI parameters.
	if err := c.Context.ShouldBindUri(obj); err != nil {
		return err
	}

	return nil
}

type ValidationError map[string]string

func (v ValidationError) Error() string {
	errStr := ""
	for field, reason := range v {
		errStr += field + ": " + reason + "\n"
	}

	return errStr
}

func (v ValidationError) StatusCode() int {
	return http.StatusBadRequest
}

func convertValidationErr(ctx Context, err error) error {
	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		return err
	}

	localize := i18n.FromContext(ctx)
	errs := make(ValidationError, len(validationErrs))
	for _, e := range validationErrs {
		errs[e.Field()] = localize.Localize(e.Tag(), map[string]interface{}{
			"Field":     e.Field(),
			"Value":     e.Value(),
			"Condition": e.Param(),
		})
	}
	return errs
}
