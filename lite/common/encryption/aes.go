package encryption

import (
	"crypto/aes"
	"crypto/cipher"
)

func AesGcmEncrypt(key []byte, iv []byte, plainText []byte, add []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, iv, plainText, add)
	return ciphertext, nil
}
