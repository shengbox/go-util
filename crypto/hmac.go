package crypto

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

func HMACSHA1(secret, value string) string {
	key := []byte(secret)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(value))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

func HmacSHA1Base64(keyStr, value string) string {
	key := []byte(keyStr)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(value))
	res := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return res
}
