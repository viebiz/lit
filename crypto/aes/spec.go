package aes

type Key interface {
	GCMEncrypt(plainText string, nonce []byte) (string, error)
	GCMDecrypt(encryptedText string, nonce []byte) (string, error)
	CBCEncrypt(plainText string, iv []byte) (string, error)
	CBCDecrypt(encryptedText string, iv []byte) (string, error)
}
