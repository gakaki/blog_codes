package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/url"
	"os"
	"testing"
)

func ReadToStr(fileName string) string {
	str, _ := os.ReadFile(fileName)
	return string(str)
}
func TestRsa2Rsa2PriSign(t *testing.T) {

	publicKey := ReadToStr("./publicKey.pem")

	validateCodeEncoded := ReadToStr("./validateCode.txt")
	validateCode, err := url.QueryUnescape(validateCodeEncoded)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("validateCode>>%s\n", validateCode)

	orderDig := ReadToStr("./orderDig.txt")
	fmt.Printf("orderDig>>%s\n", orderDig)

	res := RSAPubCheckSign(orderDig, validateCode, publicKey, crypto.SHA256)
	fmt.Println("验签结果>>", res)
}

// sha256
func SHA256Str(input string) string {
	s := "sha256 this string"
	h := sha256.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	fmt.Println(s)
	return fmt.Sprintf("%x\n", bs)
}

// RSA2私钥签名
func RSASign(signContent string, privateKey string, hash crypto.Hash) string {
	shaNew := hash.New()
	shaNew.Write([]byte(signContent))
	hashed := shaNew.Sum(nil)
	priKey, err := ParsePrivateKey(privateKey)
	if err != nil {
		return ""
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, priKey, hash, hashed)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(signature)
}

// 解析私钥
func ParsePrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("私钥信息错误！")
	}
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priKey, nil
}

// RSA2公钥验证签名
func RSAPubCheckSign(signContent, sign, publicKey string, hash crypto.Hash) bool {
	hashed := sha256.Sum256([]byte(signContent))
	pubKey, err := ParsePublicKey(publicKey)
	if err != nil {
		return false
	}
	sig, _ := base64.StdEncoding.DecodeString(sign)
	err = rsa.VerifyPKCS1v15(pubKey, hash, hashed[:], sig)
	if err != nil {
		return false
	}
	return true
}

// 解析公钥
func ParsePublicKey(publicKey string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("公钥信息错误！")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubKey.(*rsa.PublicKey), nil
}

// 生成RSA私钥和公钥，保存到文件中
// bits 证书大小
func GenerateRSAKey(bits int) {
	//GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	//Reader是一个全局、共享的密码用强随机数生成器
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}
	//保存私钥
	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	//使用pem格式对x509输出的内容进行编码
	//创建文件保存私钥
	privateFile, err := os.Create("private.pem")
	if err != nil {
		panic(err)
	}
	defer privateFile.Close()
	//构建一个pem.Block结构体对象
	privateBlock := pem.Block{Type: "RSA Private Key", Bytes: X509PrivateKey}
	//将数据保存到文件
	pem.Encode(privateFile, &privateBlock)

	//保存公钥
	//获取公钥的数据
	publicKey := privateKey.PublicKey
	//X509对公钥编码
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	//pem格式编码
	//创建用于保存公钥的文件
	publicFile, err := os.Create("public.pem")
	if err != nil {
		panic(err)
	}
	defer publicFile.Close()
	//创建一个pem.Block结构体对象
	publicBlock := pem.Block{Type: "RSA Public Key", Bytes: X509PublicKey}
	//保存到文件
	pem.Encode(publicFile, &publicBlock)
}
