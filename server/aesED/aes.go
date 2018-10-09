package aesED

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

var key = []byte{0xBA, 0x37, 0x2F, 0x02, 0xC3, 0x72, 0x1F, 0x7D,
	0x7A, 0x3D, 0x5F, 0x06, 0x41, 0x9B, 0x3F, 0x2D,
	0xBA, 0x97, 0x2F, 0x32, 0xC3, 0x92, 0x1F, 0x7D,
	0x7A, 0x3D, 0x5F, 0x04, 0x41, 0x9B, 0x3F, 0x2D,
}

func Encrypt(text string) (string, error) {
	var iv = key[:aes.BlockSize]
	encrypted := make([]byte, len(text))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	encrypter := cipher.NewCFBEncrypter(block, iv)
	encrypter.XORKeyStream(encrypted, []byte(text))
	return hex.EncodeToString(encrypted), nil
}

func Decrypt(encrypted string) (string, error) {
	var err error
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	src, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	var iv = key[:aes.BlockSize]
	decrypted := make([]byte, len(src))
	var block cipher.Block
	block, err = aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	decrypter := cipher.NewCFBDecrypter(block, iv)
	decrypter.XORKeyStream(decrypted, src)
	return string(decrypted), nil
}