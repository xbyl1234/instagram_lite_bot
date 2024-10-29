package shadowsocks

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"github.com/mzz2017/softwind/pool"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"
	"io"
)

type CipherConf struct {
	KeyLen    int
	SaltLen   int
	NonceLen  int
	TagLen    int
	NewCipher func(key []byte) (cipher.AEAD, error)
}

const (
	MaxNonceSize = 12
	ATypeIPv4    = 1
	ATypeDomain  = 3
	ATypeIpv6    = 4
)

var (
	CiphersConf = map[string]CipherConf{
		"chacha20-ietf-poly1305": {KeyLen: 32, SaltLen: 32, NonceLen: 12, TagLen: 16, NewCipher: chacha20poly1305.New},
		"chacha20-poly1305":      {KeyLen: 32, SaltLen: 32, NonceLen: 12, TagLen: 16, NewCipher: chacha20poly1305.New},
		"aes-256-gcm":            {KeyLen: 32, SaltLen: 32, NonceLen: 12, TagLen: 16, NewCipher: NewGcm},
		"aes-128-gcm":            {KeyLen: 16, SaltLen: 16, NonceLen: 12, TagLen: 16, NewCipher: NewGcm},
	}
	ZeroNonce  [MaxNonceSize]byte
	ReusedInfo = []byte("ss-subkey")
)

func NewGcm(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func MD5Sum(d []byte) []byte {
	h := md5.New()
	h.Write(d)
	return h.Sum(nil)
}

func EVPBytesToKey(password string, keyLen int) (key []byte) {
	const md5Len = 16

	cnt := (keyLen-1)/md5Len + 1
	m := make([]byte, cnt*md5Len)
	copy(m, MD5Sum([]byte(password)))

	// Repeatedly call md5 until bytes generated is enough.
	// Each call to md5 uses data: prev md5 sum + password.
	d := make([]byte, md5Len+len(password))
	start := 0
	for i := 1; i < cnt; i++ {
		start += md5Len
		copy(d, m[start-md5Len:start])
		copy(d[md5Len:], password)
		copy(m[start:], MD5Sum(d))
	}
	return m[:keyLen]
}

func (conf *CipherConf) Verify(buf []byte, masterKey []byte, salt []byte, cipherText []byte, subKey *[]byte) ([]byte, bool) {
	var sk []byte
	if subKey != nil && len(*subKey) == conf.KeyLen {
		sk = *subKey
	} else {
		sk = pool.Get(conf.KeyLen)
		defer pool.Put(sk)
		kdf := hkdf.New(
			sha1.New,
			masterKey,
			salt,
			ReusedInfo,
		)
		io.ReadFull(kdf, sk)
		if subKey != nil && cap(*subKey) >= conf.KeyLen {
			*subKey = (*subKey)[:conf.KeyLen]
			copy(*subKey, sk)
		}
	}

	ciph, _ := conf.NewCipher(sk)

	if _, err := ciph.Open(buf[:0], ZeroNonce[:conf.NonceLen], cipherText, nil); err != nil {
		return nil, false
	}
	return buf[:len(cipherText)-ciph.Overhead()], true
}

// EncryptUDPFromPool returns shadowBytes from pool.
// the shadowBytes MUST be put back.
func EncryptUDPFromPool(key Key, b []byte, salt []byte) (shadowBytes []byte, err error) {
	var buf = pool.Get(key.CipherConf.SaltLen + len(b) + key.CipherConf.TagLen)
	defer func() {
		if err != nil {
			pool.Put(buf)
		}
	}()
	copy(buf, salt)
	subKey := pool.Get(key.CipherConf.KeyLen)
	defer pool.Put(subKey)
	kdf := hkdf.New(
		sha1.New,
		key.MasterKey,
		buf[:key.CipherConf.SaltLen],
		ReusedInfo,
	)
	_, err = io.ReadFull(kdf, subKey)
	if err != nil {
		return nil, err
	}
	ciph, err := key.CipherConf.NewCipher(subKey)
	if err != nil {
		return nil, err
	}
	_ = ciph.Seal(buf[key.CipherConf.SaltLen:key.CipherConf.SaltLen], ZeroNonce[:key.CipherConf.NonceLen], b, nil)
	return buf, nil
}

// DecryptUDP will decrypt the data in place
func DecryptUDP(key Key, shadowBytes []byte) (n int, err error) {
	if len(shadowBytes) < key.CipherConf.SaltLen {
		return 0, fmt.Errorf("short length to decrypt")
	}
	subKey := pool.Get(key.CipherConf.KeyLen)
	defer pool.Put(subKey)
	kdf := hkdf.New(
		sha1.New,
		key.MasterKey,
		shadowBytes[:key.CipherConf.SaltLen],
		ReusedInfo,
	)
	_, err = io.ReadFull(kdf, subKey)
	if err != nil {
		return
	}
	ciph, err := key.CipherConf.NewCipher(subKey)
	if err != nil {
		return
	}
	plainText, err := ciph.Open(shadowBytes[key.CipherConf.SaltLen:key.CipherConf.SaltLen], ZeroNonce[:key.CipherConf.NonceLen], shadowBytes[key.CipherConf.SaltLen:], nil)
	if err != nil {
		return 0, err
	}
	copy(shadowBytes, plainText)
	return len(plainText), nil
}
