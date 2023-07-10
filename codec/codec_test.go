package codec

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAesEcb(t *testing.T) {
	var (
		key     = []byte("q4t7w!z%C*F-JaNu8x/A?D")
		val     = []byte("hello")
		badKey1 = []byte("aaaaaaaaa")
		// more than 32 chars
		badKey2 = []byte("aaaaaaaaaaaaaaaaaBBBBBBaaaaaaaaaaaaa")
	)
	_, err := EcbEncrypt(badKey1, val)
	assert.NotNil(t, err)
	_, err = EcbEncrypt(badKey2, val)
	assert.NotNil(t, err)
	dst, err := EcbEncrypt(key, val)
	assert.Nil(t, err)
	_, err = EcbDecrypt(badKey1, dst)
	assert.NotNil(t, err)
	_, err = EcbDecrypt(badKey2, dst)
	assert.NotNil(t, err)
	_, err = EcbDecrypt(key, val)
	// not enough block, just nil
	assert.Nil(t, err)
	src, err := EcbDecrypt(key, dst)
	assert.Nil(t, err)
	assert.Equal(t, val, src)
}

const (
	text      = "hello, world!\n"
	md5Digest = "910c8bc73110b0cd1bc5d2bcae782511"
)

func TestMd5(t *testing.T) {
	actual := fmt.Sprintf("%x", Md5([]byte(text)))
	assert.Equal(t, md5Digest, actual)
}

func BenchmarkHashFnv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := fnv.New32()
		new(big.Int).SetBytes(h.Sum([]byte(text))).Int64()
	}
}

func BenchmarkHashMd5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := md5.New()
		bytes := h.Sum([]byte(text))
		new(big.Int).SetBytes(bytes).Int64()
	}
}

func BenchmarkMurmur3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Hash([]byte(text))
	}
}

func TestGzip(t *testing.T) {
	var buf bytes.Buffer
	for i := 0; i < 100000; i++ {
		fmt.Fprint(&buf, i)
	}

	bs := Gzip(buf.Bytes())
	actual, err := Gunzip(bs)

	assert.Nil(t, err)
	assert.True(t, len(bs) < buf.Len())
	assert.Equal(t, buf.Bytes(), actual)
}

func TestHmac(t *testing.T) {
	ret := Hmac([]byte("foo"), "bar")
	assert.Equal(t, "f9320baf0249169e73850cd6156ded0106e2bb6ad8cab01b7bbbebe6d1065317",
		fmt.Sprintf("%x", ret))
}

func TestHmacBase64(t *testing.T) {
	ret := HmacBase64([]byte("foo"), "bar")
	assert.Equal(t, "+TILrwJJFp5zhQzWFW3tAQbiu2rYyrAbe7vr5tEGUxc=", ret)
}

const (
	priKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQC4TJk3onpqb2RYE3wwt23J9SHLFstHGSkUYFLe+nl1dEKHbD+/
Zt95L757J3xGTrwoTc7KCTxbrgn+stn0w52BNjj/kIE2ko4lbh/v8Fl14AyVR9ms
fKtKOnhe5FCT72mdtApr+qvzcC3q9hfXwkyQU32pv7q5UimZ205iKSBmgQIDAQAB
AoGAM5mWqGIAXj5z3MkP01/4CDxuyrrGDVD5FHBno3CDgyQa4Gmpa4B0/ywj671B
aTnwKmSmiiCN2qleuQYASixes2zY5fgTzt+7KNkl9JHsy7i606eH2eCKzsUa/s6u
WD8V3w/hGCQ9zYI18ihwyXlGHIgcRz/eeRh+nWcWVJzGOPUCQQD5nr6It/1yHb1p
C6l4fC4xXF19l4KxJjGu1xv/sOpSx0pOqBDEX3Mh//FU954392rUWDXV1/I65BPt
TLphdsu3AkEAvQJ2Qay/lffFj9FaUrvXuftJZ/Ypn0FpaSiUh3Ak3obBT6UvSZS0
bcYdCJCNHDtBOsWHnIN1x+BcWAPrdU7PhwJBAIQ0dUlH2S3VXnoCOTGc44I1Hzbj
Rc65IdsuBqA3fQN2lX5vOOIog3vgaFrOArg1jBkG1wx5IMvb/EnUN2pjVqUCQCza
KLXtCInOAlPemlCHwumfeAvznmzsWNdbieOZ+SXVVIpR6KbNYwOpv7oIk3Pfm9sW
hNffWlPUKhW42Gc+DIECQQDmk20YgBXwXWRM5DRPbhisIV088N5Z58K9DtFWkZsd
OBDT3dFcgZONtlmR1MqZO0pTh30lA4qovYj3Bx7A8i36
-----END RSA PRIVATE KEY-----`
	pubKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC4TJk3onpqb2RYE3wwt23J9SHL
FstHGSkUYFLe+nl1dEKHbD+/Zt95L757J3xGTrwoTc7KCTxbrgn+stn0w52BNjj/
kIE2ko4lbh/v8Fl14AyVR9msfKtKOnhe5FCT72mdtApr+qvzcC3q9hfXwkyQU32p
v7q5UimZ205iKSBmgQIDAQAB
-----END PUBLIC KEY-----`
	testBody = `this is the content`
)

func TestCryption(t *testing.T) {
	enc, err := NewRsaEncrypter([]byte(pubKey))
	assert.Nil(t, err)
	ret, err := enc.Encrypt([]byte(testBody))
	assert.Nil(t, err)

	file, err := ioutil.TempFile(os.TempDir(), Md5Hex([]byte(text)))
	assert.Nil(t, err)
	ioutil.WriteFile(file.Name(), []byte(text), os.ModeTemporary)
	assert.Nil(t, err)
	filename := file.Name()
	err = file.Close()
	assert.Nil(t, err)

	dec, err := NewRsaDecrypter(filename)
	assert.Nil(t, err)
	actual, err := dec.Decrypt(ret)
	assert.Nil(t, err)
	assert.Equal(t, testBody, string(actual))

	actual, err = dec.DecryptBase64(base64.StdEncoding.EncodeToString(ret))
	assert.Nil(t, err)
	assert.Equal(t, testBody, string(actual))
}

func TestBadPubKey(t *testing.T) {
	_, err := NewRsaEncrypter([]byte("foo"))
	assert.Equal(t, ErrPublicKey, err)
}
