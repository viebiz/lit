package jwt

type ParserOptions func(*Parser[Claims])

// WithSigningMethod add more signing method to parser
func WithSigningMethod(method SigningMethod) ParserOptions {
	return func(p *Parser[Claims]) {
		p.signingMethods[method.Alg()] = method
	}
}

// WithSigningMethods overrides all current supported signing methods
func WithSigningMethods(signingMethods map[string]SigningMethod) ParserOptions {
	return func(p *Parser[Claims]) {
		p.signingMethods = signingMethods
	}
}
