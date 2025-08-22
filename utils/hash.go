package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"strings"
)

// 固定IV，等价于C#中的Keys
var aesIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}

// PKCS7填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// PKCS7去填充
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("invalid padding size")
	}
	unpadding := int(data[length-1])
	if unpadding > length {
		return nil, errors.New("invalid padding size")
	}
	return data[:(length - unpadding)], nil
}

// AES加密
func EncryptAES(plainText, key string) (string, error) {
	if len(key) < 16 {
		return "", errors.New("key length must be at least 16")
	}
	rgbKey := []byte(key[:16])
	block, err := aes.NewCipher(rgbKey)
	if err != nil {
		return "", err
	}
	origData := pkcs7Padding([]byte(plainText), block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, aesIV)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	b64 := base64.StdEncoding.EncodeToString(crypted)
	b64 = strings.ReplaceAll(b64, "+", "_")
	b64 = strings.ReplaceAll(b64, "/", "~")
	return b64, nil
}

// AES解密
func DecryptAES(cipherText, key string) (string, error) {
	if len(key) < 16 {
		return "", errors.New("key length must be at least 16")
	}
	rgbKey := []byte(key[:16])
	b64 := strings.ReplaceAll(cipherText, "_", "+")
	b64 = strings.ReplaceAll(b64, "~", "/")
	crypted, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(rgbKey)
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, aesIV)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData, err = pkcs7UnPadding(origData)
	if err != nil {
		return "", err
	}
	return string(origData), nil
}

// 尝试解密
func TryDecryptAES(cipherText, key string) (string, bool) {
	plain, err := DecryptAES(cipherText, key)
	return plain, err == nil
}
