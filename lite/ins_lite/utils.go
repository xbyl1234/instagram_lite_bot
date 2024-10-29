package ins_lite

import (
	"CentralizedControl/common/utils"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"
)

func RSAEncrypt(pubKey []byte, plainText []byte) ([]byte, error) {
	//block, _ := pem.Decode(pubKey)
	publicKeyInterface, err := x509.ParsePKIXPublicKey(pubKey)
	if err != nil {
		return nil, err
	}
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		return nil, err
	}
	return cipherText, nil
}

func AesGcmEncrypt(key []byte, iv []byte, plainText []byte, add []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ciphertext := aesgcm.Seal(nil, iv, plainText, add)
	return ciphertext, nil
}

func (this *InsLiteClient) MakePassword(passwd string) string {
	encId := this.Cookies.PropStore54.GetInt(PropKeyPasswdPubKeyId, -1)
	if encId == -1 {
		panic("encId is null")
	}
	encPubKey := this.Cookies.PropStore54.GetBytes(PropKeyPasswdPubKey, nil)
	if encPubKey == nil {
		panic("encPubKey is null")
	}
	flag := this.Cookies.PropStore54.GetInt(PropKeyPasswdFlag, -1)
	if flag == -1 {
		panic("flag is null")
	}
	_time := strconv.FormatInt(time.Now().Unix(), 10)
	randKey := utils.GenString(utils.CharSet_All, 32)
	iv := utils.GenString(utils.CharSet_All, 12)
	randKeyEncrypted, err := RSAEncrypt(encPubKey, []byte(randKey))
	if err != nil {
		panic(err)
	}
	passwordEncrypted, err := AesGcmEncrypt([]byte(randKey),
		[]byte(iv), []byte(passwd), []byte(_time))
	if err != nil {
		panic(err)
	}
	buff := bytes.Buffer{}
	buff.WriteByte(1)
	buff.WriteByte(byte(encId))
	buff.WriteString(iv)
	lenByte := make([]byte, 2)
	binary.LittleEndian.PutUint16(lenByte, uint16(len(randKeyEncrypted)))
	buff.Write(lenByte)
	buff.Write(randKeyEncrypted)
	buff.Write(passwordEncrypted[len(passwordEncrypted)-16:])
	buff.Write(passwordEncrypted[:len(passwordEncrypted)-16])
	return fmt.Sprintf("#PWD_LITE4A:%d:%d:%s",
		flag, time.Now().Unix(), base64.StdEncoding.EncodeToString(buff.Bytes()))
}

func (this *InsLiteClient) getTrackingInstanceKey() int32 {
	return int32(atomic.AddUint32(&this.Cookies.Temp.InstanceKey, 1))
}
