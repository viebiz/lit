// Package jwt provided method for signing and verify JWT
//
// Supported signing methods: RSA, HMAC, ECDSA
//
// Usage example:
//
//	func main() {
//		// This jwt using RSA signing method
//		signingMethod := jwt.NewRS256()
//		exp := time.Now().UTC().Add(72 * time.Hour).Unix()
//
//		// Create a RSA private key for sign token
//		privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		var payload jwt.Claims
//		payload = jwt.RegisteredClaims{
//			Issuer:    "https://limitless.mukagen.com",
//			Audience:  []string{"https://resource-api.com"},
//			ExpiresAt: &exp,
//		}
//
//		// Create a new token with provided signing method and payload, and sign it with the private key
//		tokenString, err := jwt.NewToken(signingMethod, payload).SignedString(privateKey)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		fmt.Printf("JWT: %s\n", tokenString)
//
//		tokenBody := strings.Split(tokenString, ".")[1]
//		claims, err := base64.RawURLEncoding.DecodeString(tokenBody)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		fmt.Printf("Claims: %s\n", claims)
//
//		// Parse and verify the token
//		tk, err := jwt.NewParser[jwt.RegisteredClaims]().Parse(tokenString, func(_ string) (crypto.PublicKey, error) {
//			return privateKey.Public(), nil
//		})
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		fmt.Printf("Token: %+v\n", tk)
//	}
package jwt
