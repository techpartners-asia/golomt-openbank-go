package openbank

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

//	func hashKey(key string) string {
//		h := sha256.New()
//		io.WriteString(h, key)
//		return fmt.Sprintf("%x", h.Sum(nil))
//	}
func PKCS7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func (g openbank) EncryptAESCBC(text string) (string, error) {
	block, err := aes.NewCipher([]byte(g.sessionKey))
	if err != nil {
		return "", err
	}

	if len(g.ivKey) != aes.BlockSize {
		return "", fmt.Errorf("IV length must be %d bytes", aes.BlockSize)
	}

	// Pad the plaintext
	plaintext := PKCS7Pad([]byte(text), aes.BlockSize)

	encrypted := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, []byte(g.ivKey))
	mode.CryptBlocks(encrypted, plaintext)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func PKCS7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	padding := int(data[length-1])
	if padding > length || padding == 0 {
		return nil, fmt.Errorf("invalid PKCS#7 padding")
	}
	return data[:length-padding], nil
}

func (g openbank) DecryptAESCBC(ciphertext string) (string, error) {
	block, err := aes.NewCipher([]byte(g.sessionKey))
	if err != nil {
		return "", err
	}

	if len(g.ivKey) != aes.BlockSize {
		return "", fmt.Errorf("IV length must be %d bytes", aes.BlockSize)
	}

	// Decode the base64-encoded ciphertext
	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	decrypted := make([]byte, len(decodedCiphertext))
	mode := cipher.NewCBCDecrypter(block, []byte(g.ivKey))
	mode.CryptBlocks(decrypted, decodedCiphertext)

	// Unpad the decrypted data
	unpaddedData, err := PKCS7Unpad(decrypted)
	if err != nil {
		return "", err
	}

	return string(unpaddedData), nil
}
