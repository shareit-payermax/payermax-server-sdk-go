package payermax

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
	"hash"
	"strings"
)

func FormatPrivateKey(privateKey string) (pKey string) {
	var buffer strings.Builder
	buffer.WriteString("-----BEGIN RSA PRIVATE KEY-----\n")
	rawLen := 64
	keyLen := len(privateKey)
	raws := keyLen / rawLen
	temp := keyLen % rawLen
	if temp > 0 {
		raws++
	}
	start := 0
	end := start + rawLen
	for i := 0; i < raws; i++ {
		if i == raws-1 {
			buffer.WriteString(privateKey[start:])
		} else {
			buffer.WriteString(privateKey[start:end])
		}
		buffer.WriteByte('\n')
		start += rawLen
		end = start + rawLen
	}
	buffer.WriteString("-----END RSA PRIVATE KEY-----\n")
	pKey = buffer.String()
	return
}

func FormatPublicKey(publicKey string) (pKey string) {
	var buffer strings.Builder
	buffer.WriteString("-----BEGIN PUBLIC KEY-----\n")
	rawLen := 64
	keyLen := len(publicKey)
	raws := keyLen / rawLen
	temp := keyLen % rawLen
	if temp > 0 {
		raws++
	}
	start := 0
	end := start + rawLen
	for i := 0; i < raws; i++ {
		if i == raws-1 {
			buffer.WriteString(publicKey[start:])
		} else {
			buffer.WriteString(publicKey[start:end])
		}
		buffer.WriteByte('\n')
		start += rawLen
		end = start + rawLen
	}
	buffer.WriteString("-----END PUBLIC KEY-----\n")
	pKey = buffer.String()
	return
}

func DecodePublicKey(publicKeyStr string) (publicKey *rsa.PublicKey, err error) {
	key := FormatPublicKey(publicKeyStr)
	pemContent := []byte(key)
	block, _ := pem.Decode(pemContent)
	if block == nil {
		return nil, fmt.Errorf("pem.Decode(%s)：pemContent decode error", pemContent)
	}
	switch block.Type {
	case "CERTIFICATE":
		pubKeyCert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("x509.ParseCertificate(%s)：%w", pemContent, err)
		}
		pubKey, ok := pubKeyCert.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("certificate get publicKey error [%s]", pemContent)
		}
		publicKey = pubKey
	case "PUBLIC KEY":
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("x509.ParsePKIXPublicKey(%s),err:%w", pemContent, err)
		}
		pubKey, ok := pub.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("publicKey parse error [%s]", pemContent)
		}
		publicKey = pubKey
	case "RSA PUBLIC KEY":
		pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("x509.ParsePKCS1PublicKey(%s)：%w", pemContent, err)
		}
		publicKey = pubKey
	}
	return publicKey, nil
}

func DecodePrivateKey(privateKeyStr string) (privateKey *rsa.PrivateKey, err error) {
	key := FormatPrivateKey(privateKeyStr)
	pemContent := []byte(key)
	block, _ := pem.Decode(pemContent)
	if block == nil {
		return nil, fmt.Errorf("pem.Decode(%s)：pemContent decode error", pemContent)
	}
	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		pk8, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("privateKey parse error [%s]", pemContent)
		}
		var ok bool
		privateKey, ok = pk8.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("privateKey parse error [%s]", pemContent)
		}
	}
	return privateKey, nil
}

func GetRsaSign(body string, privateKey *rsa.PrivateKey) (sign string, err error) {
	var (
		h              hash.Hash
		hashs          crypto.Hash
		encryptedBytes []byte
	)

	h = sha256.New()
	hashs = crypto.SHA256

	if _, err = h.Write([]byte(body)); err != nil {
		return
	}
	if encryptedBytes, err = rsa.SignPKCS1v15(rand.Reader, privateKey, hashs, h.Sum(nil)); err != nil {
		return "", fmt.Errorf("[%w]: %v", errors.New("sign error"), err)
	}
	sign = base64.StdEncoding.EncodeToString(encryptedBytes)
	return
}

func VerifySign(signData, sign string, publicKey *rsa.PublicKey) (err error) {
	var (
		h     hash.Hash
		hashs crypto.Hash
	)
	if err != nil {
		return err
	}
	signBytes, _ := base64.StdEncoding.DecodeString(sign)

	hashs = crypto.SHA256
	h = hashs.New()
	h.Write([]byte(signData))
	if err = rsa.VerifyPKCS1v15(publicKey, hashs, h.Sum(nil), signBytes); err != nil {
		return fmt.Errorf("[%w]: %v", errors.New("verify sign error"), err)
	}
	return nil
}

func genRsaKey() (prvkey, pubkey string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	prvkey = string(pem.EncodeToMemory(block))
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		panic(err)
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubkey = string(pem.EncodeToMemory(block))
	return
}

func CreateKeyPair() (privateKey, publicKey string) {
	prvkey, pubkey := genRsaKey()
	privateKey = strings.Replace(strings.Replace(strings.Replace(prvkey, "-----BEGIN RSA PRIVATE KEY-----", "", -1), "-----END RSA PRIVATE KEY-----", "", -1), "\n", "", -1)
	publicKey = strings.Replace(strings.Replace(strings.Replace(pubkey, "-----BEGIN PUBLIC KEY-----", "", -1), "-----END PUBLIC KEY-----", "", -1), "\n", "", -1)
	return
}
