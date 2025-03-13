package i18n

import (
	"errors"
)

var (
	// ErrGivenLangNotSupported means the given lang is not supported
	ErrGivenLangNotSupported = errors.New("given lang not supported")

	// ErrBundleNotInitialized means the bundleMessage was not initialized before using
	ErrBundleNotInitialized = errors.New("bundleMessage not initialized")
)
