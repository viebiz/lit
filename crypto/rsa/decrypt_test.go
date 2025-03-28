package rsa

import (
	"crypto/rsa"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_impl_Decrypt(t *testing.T) {
	type arg struct {
		givenPrivateKeyFn  func() *rsa.PrivateKey
		givenPrivateKey    *rsa.PrivateKey
		givenEncryptedText string
		expDecryptedText   string
		expErr             error
	}
	tcs := map[string]arg{
		"err: illegal base64 data at input": {
			givenPrivateKeyFn: func() *rsa.PrivateKey {
				return nil
			},
			givenEncryptedText: "invalid base64 string",
			expErr:             errors.New("illegal base64 data at input byte 7"),
		},
		"err: decryption error": {
			givenPrivateKeyFn: func() *rsa.PrivateKey {
				privateKey, err := newPrivateKey("./testdata/test_rsa_private_key.pem")
				require.NoError(t, err)
				return privateKey
			},
			givenEncryptedText: "DYUQ+JQdxMaVIZt9k6e3vRsD48tOKcT95kQ5R8NcqywoU6YHB8/UotQmG1Iqq8ZHGslqCogu0cGgSvp0Ifam8yWcNEPy+mcp1gXhkyx3b0DTClsFvCxl5lWTDwwmwq4snaSmR+u4amcWSGg0brs2oqjWRFpI4YEttSHsrCxO1k4dg9WKd55N64Dz8hdEO3bCJHFaY307uhVhpzfHGoCvDyHHI4y15k5DnUxMun9W6ckPeN0pWCAXXi+9HmJjLx1ZZkRZ19ii6PoJwH9/rmb+/FQk0pEr0DtUF9ghfS5OyJLbTN0dmMAE1ithwIeoIqUc0+OzlHUWAdNzgxv/hXmACQ==",
			expErr:             errors.New("crypto/rsa: decryption error"),
		},
		"ok": {
			givenPrivateKeyFn: func() *rsa.PrivateKey {
				privateKey, err := newPrivateKey("./testdata/rsa_private_key.pem")
				require.NoError(t, err)
				return privateKey
			},
			givenEncryptedText: "EoNoHBJNqLqelvRh8z20L3AyP+JkSmu7r4NbC+hY1DNPnbIhzUrYbzcY6pyHdU4YEhk/SmhNLRoc2rBtwNjTk0qPhxJKUuHMVU9QfVsMNZSW/xLD+hauuE4zebULufZVEbMNXSgn9kBKdWZSVUatvt5pv+8BeTZcTMP+b8fAGzVxzXFIRCSo9kI2eiVXPEIRj4Gd8a/JxdFjtHZFhE5rRMlTWT/5ZcE2+EW5UWoIv3yAS1S+OE4LKxW/z+n1PfUhw+P/iuUEj9h58t4RgrmiZkhjM0qcmnYHRCHftoD9fNUdpfZ4lhBomvYRxs13NI5JRveAS17CbylVFTHJ0P7Cng==",
			expDecryptedText:   "This is a plain text",
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Given

			// When
			instance := &impl{
				privateKey: tc.givenPrivateKeyFn(),
			}
			result, err := instance.Decrypt(tc.givenEncryptedText)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expDecryptedText, result)
			}
		})
	}
}
