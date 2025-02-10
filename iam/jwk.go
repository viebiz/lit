package iam

// JWK represents a JSON Web Key
// Refer https://datatracker.ietf.org/doc/html/rfc7517
type JWK struct {
	KID string   `json:"kid"`
	Kty string   `json:"kty"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
	X5t string   `json:"x5t"`
	Alg string   `json:"alg"`
}

// JWKSet set of keys containing the public keys used to verify any JSON Web Token (JWT)
// issued by the Authorization Server and signed using the RS256 signing algorithm
type JWKSet struct {
	Keys []JWK `json:"keys"`
}
