package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"hash"
	"regexp"

	"golang.org/x/crypto/pbkdf2"
)

type pbkdf2_arg struct {
	iter    int              //迭代次数
	keyLen  int              //密文长度
	saltLen int              //盐长度
	hmac    func() hash.Hash //哈希算法
}

var GrubPBKDF2 *pbkdf2_arg

func InitGrubPBKDF2() *pbkdf2_arg {
	GrubPBKDF2 = &pbkdf2_arg{
		iter:    10000,
		keyLen:  64,
		saltLen: 64,
		hmac:    sha512.New,
	}
	return GrubPBKDF2
}

func GrubPBKDF2Crypto(usr, passwd string) (string, error) {
	ok := CheckUsername(usr)
	if !ok {
		return "", errors.New("Account cannot contain special characters.")
	}
	if GrubPBKDF2 == nil {
		return "", errors.New("Grub PBKDF2 don't init.")
	}

	//生成随机盐
	salt := make([]byte, GrubPBKDF2.saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	//生成密文
	dk := pbkdf2.Key([]byte(passwd), salt, GrubPBKDF2.iter, GrubPBKDF2.keyLen, GrubPBKDF2.hmac)
	chiperPasswd := "password_pbkdf2 " + usr + " grub.pbkdf2.sha512." +
		ToString(GrubPBKDF2.iter) + "." + hex.EncodeToString(salt) + "." + hex.EncodeToString(dk)

	return chiperPasswd, nil

}

//帐号校验
func CheckUsername(username string) (b bool) {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{3,}$", username); !ok {
		return false
	}
	return true
}

func ToString(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

func IsInSlice(element string, elements []string) (isIn bool) {
	for _, item := range elements {
		if element == item {
			isIn = true
			return
		}
	}
	return
}
