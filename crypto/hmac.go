package crypto

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
)

func HMACSHA1(secret, value string) string {
	//hmac ,use sha1
	key := []byte(secret)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(value))
	return fmt.Sprintf("%x", mac.Sum(nil))
}
