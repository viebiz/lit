package aes

import "errors"

var (
	ErrInvalidData         = errors.New("invalid data")
	ErrInvalidPaddingValue = errors.New("invalid padding value")
	ErrIncorrectPadding    = errors.New("padding data is incorrect")
)
