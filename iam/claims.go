package iam

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/viebiz/lit/jwt"
)

type Claims struct {
	RegisteredClaims jwt.RegisteredClaims
	ExtraClaims      map[string]interface{}
}

func (c Claims) Valid() error {
	return nil
}

func (c Claims) MarshalJSON() ([]byte, error) {
	// Initial claim map, that contains all claims in plain
	claimMap := make(map[string]interface{})

	// Add registered claims to claim map
	claimMap["iss"] = c.RegisteredClaims.Issuer

	claimMap["sub"] = c.RegisteredClaims.Subject

	if len(c.RegisteredClaims.Audience) > 0 {
		claimMap["aud"] = c.RegisteredClaims.Audience
	}

	if c.RegisteredClaims.IssuedAt != nil {
		claimMap["iat"] = c.RegisteredClaims.IssuedAt
	}

	if c.RegisteredClaims.ExpiresAt != nil {
		claimMap["exp"] = c.RegisteredClaims.ExpiresAt
	}

	if c.RegisteredClaims.NotBefore != nil {
		claimMap["nbf"] = c.RegisteredClaims.NotBefore
	}

	claimMap["jti"] = c.RegisteredClaims.JTI

	claimMap["client_id"] = c.RegisteredClaims.ClientID

	// Add extra claims to claim map
	for k, v := range c.ExtraClaims {
		claimMap[k] = v
	}

	return json.Marshal(claimMap)
}

func (c *Claims) UnmarshalJSON(data []byte) error {
	var rawMap map[string]interface{}
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	aud, err := unmarshalAud(rawMap)
	if err != nil {
		return err
	}

	iat, err := unmarshalTimestamp(rawMap, "iat")
	if err != nil {
		return err
	}

	exp, err := unmarshalTimestamp(rawMap, "exp")
	if err != nil {
		return err
	}

	nbf, err := unmarshalTimestamp(rawMap, "nbf")
	if err != nil {
		return err
	}

	c.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    unmarshalString(rawMap, "iss"),
		Subject:   unmarshalString(rawMap, "sub"),
		Audience:  aud,
		IssuedAt:  iat,
		ExpiresAt: exp,
		NotBefore: nbf,
		JTI:       unmarshalString(rawMap, "jti"),
		ClientID:  unmarshalString(rawMap, "client_id"),
	}

	// Remove registered claims from raw map
	delete(rawMap, "iss")
	delete(rawMap, "sub")
	delete(rawMap, "aud")
	delete(rawMap, "iat")
	delete(rawMap, "exp")
	delete(rawMap, "nbf")
	delete(rawMap, "jti")
	delete(rawMap, "client_id")

	if len(rawMap) > 0 {
		c.ExtraClaims = make(map[string]interface{})
	}
	for k, v := range rawMap {
		c.ExtraClaims[k] = v
	}

	return nil
}

func unmarshalString(rawMap map[string]interface{}, key string) string {
	v, exists := rawMap[key]
	if !exists {
		return ""
	}

	rs, ok := v.(string)
	if !ok {
		return ""
	}

	return rs
}

func unmarshalAud(rawMap map[string]interface{}) ([]string, error) {
	aud, exists := rawMap["aud"]
	if !exists {
		return nil, nil
	}

	switch value := aud.(type) {
	case string:
		return []string{value}, nil
	case []interface{}:
		var result []string
		for _, v := range value {
			result = append(result, v.(string))
		}

		return result, nil
	case []string:
		var result []string
		for _, v := range value {
			result = append(result, v)
		}

		return result, nil
	default:
		return nil, fmt.Errorf("unknown audience type")
	}
}

func unmarshalTimestamp(rawMap map[string]interface{}, key string) (*int64, error) {
	v, exists := rawMap[key]
	if !exists {
		return nil, nil
	}

	var result int64
	var err error
	switch value := v.(type) {
	case string:
		result, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}

	case json.Number:
		result, err = value.Int64()
		if err != nil {
			return nil, err
		}

	case float64:
		result = int64(value)

	default:
		return nil, fmt.Errorf("unknown timestamp type")
	}

	return &result, nil
}
