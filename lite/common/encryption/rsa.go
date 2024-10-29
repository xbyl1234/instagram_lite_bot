package encryption

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

//	pri, _ := pem.Decode(common.PriKey)
//	pub, _ := pem.Decode(common.PubKey)
//	common.GenerateRSAKey(256)
//	cipher := common.RSAEncrypt(pub, []byte("hello zhaoyingkui "))
//	fmt.Println(cipher)
//	plain := common.RSADecrypt(pri, cipher)
//	fmt.Println(string(plain))

func GenerateRSAKey(bits int) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	privateFile, err := os.Create("private.pem")
	if err != nil {
		panic(err)
	}
	defer privateFile.Close()
	privateBlock := pem.Block{Type: "RSA Private Key", Bytes: X509PrivateKey}

	pem.Encode(privateFile, &privateBlock)

	publicKey := privateKey.PublicKey
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	publicFile, err := os.Create("public.pem")
	if err != nil {
		panic(err)
	}
	defer publicFile.Close()
	publicBlock := pem.Block{Type: "RSA Public Key", Bytes: X509PublicKey}
	pem.Encode(publicFile, &publicBlock)
}

func DecodePriKey(key []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(key)
	pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	return pri
}

func DecodePubKey(key []byte) *rsa.PublicKey {
	block, _ := pem.Decode(key)
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	return pubInterface.(*rsa.PublicKey)
}

const (
	BASE_64_FORMAT         = "UrlSafeNoPadding"
	RSA_ALGORITHM_KEY_TYPE = "PKCS8"
	RSA_ALGORITHM_SIGN     = crypto.SHA256
)

//// 生成密钥对
//func CreateKeys(publicKeyWriter, privateKeyWriter io.Writer, keyLength int) error {
//	// 生成私钥文件
//	privateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
//	if err != nil {
//		return err
//	}
//	derStream := MarshalPKCS8PrivateKey(privateKey)
//	block := &pem.Block{
//		Type:  "PRIVATE KEY",
//		Bytes: derStream,
//	}
//	err = pem.Encode(privateKeyWriter, block)
//	if err != nil {
//		return err
//	}
//	// 生成公钥文件
//	publicKey := &privateKey.PublicKey
//	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
//	if err != nil {
//		return err
//	}
//	block = &pem.Block{
//		Type:  "PUBLIC KEY",
//		Bytes: derPkix,
//	}
//	err = pem.Encode(publicKeyWriter, block)
//	if err != nil {
//		return err
//	}
//	return nil
//}

//func NewXRsa(publicKey []byte, privateKey []byte) (*XRsa, error) {
//	block, _ := pem.Decode(publicKey)
//	if block == nil {
//		return nil, errors.New("public key error")
//	}
//	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
//	if err != nil {
//		return nil, err
//	}
//	pub := pubInterface.(*rsa.PublicKey)
//	block, _ = pem.Decode(privateKey)
//	if block == nil {
//		return nil, errors.New("private key error!")
//	}
//	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
//	if err != nil {
//		return nil, err
//	}
//	pri, ok := priv.(*rsa.PrivateKey)
//	if ok {
//		return &XRsa{
//			publicKey:  pub,
//			privateKey: pri,
//		}, nil
//	} else {
//		return nil, errors.New("private key not supported")
//	}
//}

// 公钥加密
func RSAEncrypt(key *rsa.PublicKey, data []byte) ([]byte, error) {
	partLen := key.N.BitLen()/8 - 11
	chunks := split(data, partLen)
	buffer := bytes.NewBuffer([]byte{})
	for _, chunk := range chunks {
		bytes, err := rsa.EncryptPKCS1v15(rand.Reader, key, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(bytes)
	}
	return buffer.Bytes(), nil
}

// 私钥解密
func RSADecrypt(key *rsa.PrivateKey, encrypted []byte) ([]byte, error) {
	partLen := key.N.BitLen() / 8
	chunks := split(encrypted, partLen)
	buffer := bytes.NewBuffer([]byte{})
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, key, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(decrypted)
	}
	return buffer.Bytes(), nil
}

//func Sign(data string) (string, error) {
//	h := RSA_ALGORITHM_SIGN.New()
//	h.Write([]byte(data))
//	hashed := h.Sum(nil)
//	sign, err := rsa.SignPKCS1v15(rand.Reader, r.privateKey, RSA_ALGORITHM_SIGN, hashed)
//	if err != nil {
//		return "", err
//	}
//	return base64.RawURLEncoding.EncodeToString(sign), err
//}
//
//func Verify(data string, sign string) error {
//	h := RSA_ALGORITHM_SIGN.New()
//	h.Write([]byte(data))
//	hashed := h.Sum(nil)
//	decodedSign, err := base64.RawURLEncoding.DecodeString(sign)
//	if err != nil {
//		return err
//	}
//	return rsa.VerifyPKCS1v15(r.publicKey, RSA_ALGORITHM_SIGN, hashed, decodedSign)
//}

//func MarshalPKCS8PrivateKey(key *rsa.PrivateKey) []byte {
//	info := struct {
//		Version             int
//		PrivateKeyAlgorithm []asn1.ObjectIdentifier
//		PrivateKey          []byte
//	}{}
//	info.Version = 0
//	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
//	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
//	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)
//	k, _ := asn1.Marshal(info)
//	return k
//}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:])
	}
	return chunks
}
