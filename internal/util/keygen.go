package util

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
)

// 生成指定字节长度的随机字节序列
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// 生成Base64编码的随机密钥
func GenerateBase64Key(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// 生成十六进制编码的随机密钥
func GenerateHexKey(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
