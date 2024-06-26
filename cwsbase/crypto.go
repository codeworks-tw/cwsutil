/*
 * File: crypto.go
 * Created Date: Thursday, April 11th 2024, 10:31:37 am
 *
 * Last Modified: Tue Jun 04 2024
 * Modified By: hsky77
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwsbase

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
)

func EncryptMap(input map[string]any) (string, error) {
	data, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	return aESCBCPKCS5PaddingEncrypt(data, aes.BlockSize)
}

func DecryptToMap(input string) (map[string]any, error) {
	data, err := aESCBCPKCS5PaddingDecrypt(input)
	if err != nil {
		return nil, err
	}

	var output map[string]any
	err = json.Unmarshal(data, &output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func aESCBCPKCS5PaddingEncrypt(plaintext []byte, blockSize int) (string, error) {
	key, err := hex.DecodeString(GetEnv[string]("CRYPTO_KEY_HEX"))
	if err != nil {
		return "", err
	}
	iv, err := hex.DecodeString(GetEnv[string]("CRYPTO_IV_HEX"))
	if err != nil {
		return "", err
	}

	bKey := []byte(key)
	bIV := []byte(iv)
	bPlaintext := pKCS5Padding(plaintext, blockSize)
	block, err := aes.NewCipher(bKey)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(ciphertext, bPlaintext)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func aESCBCPKCS5PaddingDecrypt(cipherTextBase64 string) ([]byte, error) {
	key, err := hex.DecodeString(GetEnv[string]("CRYPTO_KEY_HEX"))
	if err != nil {
		return nil, err
	}
	iv, err := hex.DecodeString(GetEnv[string]("CRYPTO_IV_HEX"))
	if err != nil {
		return nil, err
	}

	bKey := []byte(key)
	bIV := []byte(iv)
	cipherTextDecoded, err := base64.URLEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(bKey)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, bIV)
	mode.CryptBlocks([]byte(cipherTextDecoded), []byte(cipherTextDecoded))
	cipherTextDecoded, err = pKCS5Unpadding(cipherTextDecoded, aes.BlockSize)
	if err != nil {
		return nil, err
	}
	return cipherTextDecoded, nil
}

func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pKCS5Unpadding(src []byte, blockSize int) ([]byte, error) {
	srcLen := len(src)
	paddingLen := int(src[srcLen-1])
	if paddingLen >= srcLen || paddingLen > blockSize {
		return nil, errors.New("unpadding size error")
	}
	return src[:srcLen-paddingLen], nil
}
