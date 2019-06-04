package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

func Md5(str string) (string, error) {
	m := md5.New()
	if _, err := m.Write([]byte(str)); err != nil {
		return "", err
	}
	return hex.EncodeToString(m.Sum(nil)), nil
}

func Hmac(key, str string) (string, error) {
	h := hmac.New(md5.New, []byte(key))
	if _, err := h.Write([]byte(str)); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func Sha1(str string) (string, error) {
	sh := sha1.New()
	if _, err := sh.Write([]byte(str)); err != nil {
		return "", err
	}
	return hex.EncodeToString(sh.Sum(nil)), nil
}
