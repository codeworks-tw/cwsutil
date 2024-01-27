/*
 * File: crypto.go
 * Created Date: Friday, January 26th 2024, 9:49:36 am
 *
 * Last Modified: Sat Jan 27 2024
 * Modified By: Howard Ling-Hao Kung
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package baseutil

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

	return AESCBCPKCS5PaddingEncrypt(data, aes.BlockSize)
}

func DecryptToMap(input string) (map[string]any, error) {
	data, err := AESCBCPKCS5PaddingDecrypt(input)
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

func AESCBCPKCS5PaddingEncrypt(plaintext []byte, blockSize int) (string, error) {
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
	bPlaintext := PKCS5Padding(plaintext, blockSize)
	block, err := aes.NewCipher(bKey)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(ciphertext, bPlaintext)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func AESCBCPKCS5PaddingDecrypt(cipherTextBase64 string) ([]byte, error) {
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
	cipherTextDecoded, err = PKCS5Unpadding(cipherTextDecoded, aes.BlockSize)
	if err != nil {
		return nil, err
	}
	return cipherTextDecoded, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Unpadding(src []byte, blockSize int) ([]byte, error) {
	srcLen := len(src)
	paddingLen := int(src[srcLen-1])
	if paddingLen >= srcLen || paddingLen > blockSize {
		return nil, errors.New("unpadding size error")
	}
	return src[:srcLen-paddingLen], nil
}
