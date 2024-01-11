package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
)

// https://www.codeover.cn/go_rsa/
// RsaEncryptBase64 使用 RSA 公钥加密数据, 返回加密后并编码为 base64 的数据
func RsaEncryptBase64(originalData, publicKey string) (string, error) {
	block, _ := pem.Decode([]byte(publicKey))
	pubKey, parseErr := x509.ParsePKIXPublicKey(block.Bytes)
	if parseErr != nil {
		fmt.Println(parseErr)
		return "", errors.New("解析公钥失败")
	}
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey.(*rsa.PublicKey), []byte(originalData))
	return base64.StdEncoding.EncodeToString(encryptedData), err
}

// RsaDecryptBase64 使用 RSA 私钥解密数据
func RsaDecryptBase64(encryptedData, privateKey string) (string, error) {
	encryptedDecodeBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode([]byte(privateKey))
	priKey, parseErr := x509.ParsePKCS8PrivateKey(block.Bytes)
	if parseErr != nil {
		fmt.Println(parseErr)
		return "", errors.New("解析私钥失败")
	}

	originalData, encryptErr := rsa.DecryptPKCS1v15(rand.Reader, priKey.(*rsa.PrivateKey), encryptedDecodeBytes)
	return string(originalData), encryptErr
}

func main() {}

func EncryptResponseBody(rw http.ResponseWriter, r *http.Request) {

	outputBody := "output_data"
	encryptBase64, err := RsaEncryptBase64(outputBody, "public_key")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(encryptBase64))
}
