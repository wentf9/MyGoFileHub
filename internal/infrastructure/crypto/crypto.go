package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Encrypt 使用指定的 32 字节密钥将明文加密为 Base64 字符串
func Encrypt(plainText string, key string) (string, error) {
	if plainText == "" {
		return "", nil
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt 使用指定的 32 字节密钥将 Base64 密文解密为明文
func Decrypt(cipherTextB64 string, key string) (string, error) {
	if cipherTextB64 == "" {
		return "", nil
	}

	cipherText, err := base64.StdEncoding.DecodeString(cipherTextB64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(cipherText) < nonceSize {
		return "", errors.New("cipherText too short")
	}

	nonce, encryptedMessage := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := aesGCM.Open(nil, nonce, encryptedMessage, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
