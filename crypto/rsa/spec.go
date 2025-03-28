package rsa

type Cipher interface {
	Encrypt(plainText string) (string, error)
	Decrypt(encryptedText string) (string, error)
}
