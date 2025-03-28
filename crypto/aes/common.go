package aes

import "crypto/rand"

var (
	randReadFn = rand.Read
)
