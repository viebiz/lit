package iam

import (
	"strings"
)

const (
	hasPermissionKeyMatch = "hasPermission"
)

func hasPermission(args ...interface{}) (interface{}, error) {
	// keyMatchFunc return false if provided action not contains required action
	keyMatchFunc := func(required, provided string) bool {
		for _, ch := range required {
			if !strings.ContainsRune(provided, ch) {
				return false
			}
		}

		return true
	}

	provided := args[0].(string)
	required := args[1].(string)

	return keyMatchFunc(provided, required), nil
}
