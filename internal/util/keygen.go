package util

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"

)

// 易读字符集：去除了易混淆字符 0/O, 1/l/I
const ReadableChars = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// GenerateReadableKey 生成人可读的密钥，每4位用连字符分隔
func GenerateReadableKey(length int, groupSize int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be positive")
	}

	bytes := make([]byte, length)
	for i := range bytes {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(ReadableChars))))
		if err != nil {
			return "", err
		}
		bytes[i] = ReadableChars[index.Int64()]
	}

	key := string(bytes)

	// 每 groupSize 个字符加一个连字符（如 4 位一组）
	if groupSize > 0 {
		var result []byte
		for i := 0; i < len(key); i++ {
			if i > 0 && i%groupSize == 0 {
				result = append(result, '-')
			}
			result = append(result, key[i])
		}
		key = string(result)
	}

	return key, nil
}

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
