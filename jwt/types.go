package jwt

import (
	"crypto"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// Signer represents an interface for creating digital signatures
type Signer interface {
	Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error)
}

// VerifyKey represents a key for verify token
type VerifyKey interface{}

// ClaimStrings represents a claim's value as a slice of strings
type ClaimStrings []string

func (c *ClaimStrings) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(*c))
}

func (c *ClaimStrings) UnmarshalJSON(data []byte) error {
	var unmarshalled interface{}
	if err := json.Unmarshal(data, &unmarshalled); err != nil {
		return err
	}

	switch val := unmarshalled.(type) {
	case string:
		if val == "" {
			return nil
		}

		*c = []string{val}
	case []string:
		*c = val
	case []interface{}:
		rs := make([]string, len(val))
		for idx, el := range val {
			rs[idx] = fmt.Sprintf("%v", el)
		}

		*c = rs
	default:
		return errors.New("audience must be a string or []string")
	}

	return nil
}
