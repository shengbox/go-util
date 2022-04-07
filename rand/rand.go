package rand

import (
	"math/rand"
	"strconv"
	"time"
)

func GenCaptcha() string {
	rand.Seed(time.Now().Unix())
	smsCaptcha := 100000 + rand.Intn(899999)
	return strconv.Itoa(smsCaptcha)
}
