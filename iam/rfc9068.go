package iam

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/viebiz/lit/jwt"
)

const (
	jwkKeyUseSig = "sig" // JWK property `use` determines the JWK is for signature verification

	jwkAlgRSA = "RSA" //

	timeout = 30 * time.Second
)

var (
	rfc9068AllowTokenTypes = []string{
		"at+jwt",
		"application/at+jwt",
	}
)

// rfc9068Validator represents validator for validate oauth2 JWT
// Refer https://datatracker.ietf.org/doc/rfc9068/
type rfc9068Validator struct {
	// jwksURI contains the URL to fetch JSON Web Key set
	jwksURI string

	// issuer contains the address token issuer
	issuer string

	// audience contains the address of current service
	audience string

	// httpClient is client to call external service
	httpClient HTTPClient

	// signingKeyMap stores crypto.PublicKey
	// it will be initial at the first time run
	signingKeyMap map[string]crypto.PublicKey

	tokenParser jwt.Parser[Claims]
}

// Validate validates given claims
func (v *rfc9068Validator) Validate(tokenString string) (jwt.Token[Claims], error) {
	// 1. Parse token from string
	tk, err := v.tokenParser.Parse(tokenString, v.getKeyFunc)
	if err != nil {
		return jwt.Token[Claims]{}, err
	}

	// 2. Verify the "typ" header value
	if typ := tk.Header["typ"]; !slices.Contains(rfc9068AllowTokenTypes, typ) {
		return jwt.Token[Claims]{}, ErrInvalidToken
	}

	// TODO: Support encrypted token

	// 3. Verify claims information
	if err := v.validateClaims(tk.Claims); err != nil {
		return jwt.Token[Claims]{}, err
	}

	return tk, nil
}

func (v *rfc9068Validator) validateClaims(c Claims) error {
	now := timeNowFunc()

	// Verify issuer claim
	if c.RegisteredClaims.Issuer == "" {
		return ErrInvalidToken
	}

	if strings.TrimRight(c.RegisteredClaims.Issuer, "/") != v.issuer {
		return ErrInvalidToken
	}

	// Verify expires_at claim
	if c.RegisteredClaims.ExpiresAt == nil {
		return ErrInvalidToken
	}

	if exp := time.Unix(*c.RegisteredClaims.ExpiresAt, 0); now.After(exp) {
		return ErrTokenExpired
	}

	// Verify audience claim
	if len(c.RegisteredClaims.Audience) == 0 || c.RegisteredClaims.Audience[0] == "" {
		return ErrInvalidToken
	}

	foundAud := false
	for _, aud := range c.RegisteredClaims.Audience {
		if aud == v.audience {
			foundAud = true
		}
	}

	if !foundAud {
		return ErrInvalidToken
	}

	// Verify subject claim
	if c.RegisteredClaims.Subject == "" {
		return ErrInvalidToken
	}

	// Verify client_id claim
	if c.RegisteredClaims.ClientID == "" {
		return ErrInvalidToken
	}

	// Verify iat claim
	if c.RegisteredClaims.IssuedAt == nil {
		return ErrInvalidToken
	}

	// Verify jti claim
	if c.RegisteredClaims.JTI == "" {
		return ErrInvalidToken
	}

	return nil
}

func (v *rfc9068Validator) getKeyFunc(keyID string) (crypto.PublicKey, error) {
	key, exists := v.signingKeyMap[keyID]
	if !exists {
		return nil, ErrInvalidToken
	}

	return key, nil
}

func (v *rfc9068Validator) downloadSigningKey(ctx context.Context) error {
	keyset, err := v.fetchJWKS(ctx)
	if err != nil {
		return err
	}

	singingKeyMap, err := v.processJWKS(ctx, keyset)
	if err != nil {
		return err
	}

	v.signingKeyMap = singingKeyMap

	return nil
}

func (v *rfc9068Validator) fetchJWKS(ctx context.Context) (JWKSet, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.jwksURI, nil)
	if err != nil {
		return JWKSet{}, err
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return JWKSet{}, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return JWKSet{}, fmt.Errorf("got unexpected status code: %d", resp.StatusCode)
	}

	var jwks JWKSet
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return JWKSet{}, fmt.Errorf("could not decode jwks: %w", err)
	}

	return jwks, nil
}

func (v *rfc9068Validator) processJWKS(ctx context.Context, jwks JWKSet) (map[string]crypto.PublicKey, error) {
	signingKeyMap := map[string]crypto.PublicKey{}

	for _, k := range jwks.Keys {
		if k.Use != jwkKeyUseSig {
			continue
		}

		if k.Kty != jwkAlgRSA {
			continue
		}

		if k.KID == "" {
			continue
		}

		if len(k.X5c) < 1 { // If no certificate chain then skip
			continue
		}

		// pre-computing the signing key from the JWK and cert so that it does not need to be done again once cached
		pemBlock, _ := pem.Decode([]byte("-----BEGIN CERTIFICATE-----\n" + k.X5c[0] + "\n-----END CERTIFICATE-----"))
		if pemBlock == nil {
			return nil, fmt.Errorf("could not parse certificate PEM")
		}

		cert, err := x509.ParseCertificate(pemBlock.Bytes)
		if err != nil {
			return nil, err
		}

		signingKeyMap[k.KID] = cert.PublicKey
	}

	if len(signingKeyMap) == 0 {
		return nil, fmt.Errorf("no appropriate JWK found")
	}

	return signingKeyMap, nil
}
