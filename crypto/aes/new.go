package aes

type impl struct {
	key keyBytes
}

// New create AES Key 32 bytes
func New() (Key, error) {
	key := make([]byte, 32) // 256-bit key
	_, err := randReadFn(key)
	if err != nil {
		return nil, err
	}
	return &impl{
		key: key,
	}, nil
}
