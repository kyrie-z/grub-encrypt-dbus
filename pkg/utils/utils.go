package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

type pbkdf2_arg struct {
	iter    int
	keyLen  int
	saltLen int
}

var GrubPBKDF2 *pbkdf2_arg

func InitGrubPBKDF2() *pbkdf2_arg {
	GrubPBKDF2 = &pbkdf2_arg{
		iter:    10000,
		keyLen:  64,
		saltLen: 64,
	}
	return GrubPBKDF2
}

func GrubPBKDF2Crypto(usr, passwd string) (string, error) {
	if GrubPBKDF2 == nil {
		return "", fmt.Errorf("Grub PBKDF2 don't init.")
	}

	//生成随机盐
	salt := make([]byte, GrubPBKDF2.saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	//生成密文
	dk := pbkdf2.Key([]byte(passwd), salt, GrubPBKDF2.iter, GrubPBKDF2.keyLen, sha512.New)
	chiperPasswd := "password_pbkdf2 " + usr + " grub.pbkdf2.sha512." +
		ToString(GrubPBKDF2.iter) + "." + hex.EncodeToString(salt) + "." + hex.EncodeToString(dk)

	return chiperPasswd, nil

}

func ToString(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
