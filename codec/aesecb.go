package codec

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
)

// ErrPaddingSize indicates bad padding size.
var ErrPaddingSize = errors.New("padding size error")

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// NewECBEncrypter returns an ECB encrypter.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}

func (x *ecbEncrypter) BlockSize() int { return x.blockSize }

// why we don't return error is because cipher.BlockMode doesn't allow this
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		return
	}
	if len(dst) < len(src) {
		return
	}

	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// NewECBDecrypter returns an ECB decrypter.
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

func (x *ecbDecrypter) BlockSize() int {
	return x.blockSize
}

// why we don't return error is because cipher.BlockMode doesn't allow this
func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		return
	}
	if len(dst) < len(src) {
		return
	}

	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

// EcbDecrypt decrypts src with the given key.
func EcbDecrypt(key, src []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	decrypter := NewECBDecrypter(block)
	decrypted := make([]byte, len(src))
	decrypter.CryptBlocks(decrypted, src)

	return pkcs5Unpadding(decrypted, decrypter.BlockSize())
}

// EcbDecryptBase64 decrypts base64 encoded src with the given base64 encoded key.
// The returned string is also base64 encoded.
func EcbDecryptBase64(key, src string) (string, error) {
	keyBytes, err := getKeyBytes(key)
	if err != nil {
		return "", err
	}

	encryptedBytes, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}

	decryptedBytes, err := EcbDecrypt(keyBytes, encryptedBytes)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(decryptedBytes), nil
}

// EcbEncrypt encrypts src with the given key.
func EcbEncrypt(key, src []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	padded := pkcs5Padding(src, block.BlockSize())
	crypted := make([]byte, len(padded))
	encrypter := NewECBEncrypter(block)
	encrypter.CryptBlocks(crypted, padded)

	return crypted, nil
}

// EcbEncryptBase64 encrypts base64 encoded src with the given base64 encoded key.
// The returned string is also base64 encoded.
func EcbEncryptBase64(key, src string) (string, error) {
	keyBytes, err := getKeyBytes(key)
	if err != nil {
		return "", err
	}

	srcBytes, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}

	encryptedBytes, err := EcbEncrypt(keyBytes, srcBytes)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

func getKeyBytes(key string) ([]byte, error) {
	if len(key) <= 32 {
		return []byte(key), nil
	}

	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	return keyBytes, nil
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5Unpadding(src []byte, blockSize int) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])
	if unpadding >= length || unpadding > blockSize {
		return nil, ErrPaddingSize
	}

	return src[:length-unpadding], nil
}

func hashBytes(key string) (hash []byte) {
	h := sha256.New()
	io.WriteString(h, key)
	hashStr := hex.EncodeToString(h.Sum(nil))
	hash = []byte(hashStr)[:32]
	return
}

func Encrypt(plainText string, key string) (cipherText string, err error) {
	var block cipher.Block
	keyBytes := hashBytes(key)
	plainTextBytes := []byte(plainText)
	block, err = aes.NewCipher(keyBytes)
	if err != nil {
		return
	}

	cipherTextBytes := make([]byte, aes.BlockSize+len(plainTextBytes))
	iv := cipherTextBytes[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherTextBytes[aes.BlockSize:], plainTextBytes)
	cipherText = "crypt-" + hex.EncodeToString(cipherTextBytes)
	return
}

func Decrypt(cipherText string, key string) (plainText string, err error) {
	if len(cipherText) == 0 || len(cipherText) < 6 || cipherText[:6] != "crypt-" {
		err = errors.New("Illegal ciphertext")
		return
	}
	cipherText = string(cipherText[6:])
	var block cipher.Block
	keyBytes := hashBytes(key)
	cipherTextBytes, _ := hex.DecodeString(cipherText)
	block, err = aes.NewCipher(keyBytes)
	if err != nil {
		return
	}

	if len(cipherTextBytes) < aes.BlockSize {
		err = errors.New("Ciphertext too short")
		return
	}

	iv := cipherTextBytes[:aes.BlockSize]
	cipherTextBytes = cipherTextBytes[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)

	plainTextBytes := make([]byte, len(cipherTextBytes))
	stream.XORKeyStream(plainTextBytes, cipherTextBytes)
	plainText = string(plainTextBytes)
	return
}
