package chinese

import (
	"math/rand"
	"strconv"
	"time"
)

//随机生成手机号码
func NewPhone() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	array := [15]string{"130", "131", "133", "136", "139", "150", "151", "156", "177", "178", "185", "186", "188", "189", "158"}
	phone := array[r.Intn(15)]
	for i := 0; i < 8; i++ {
		phone += strconv.Itoa(r.Intn(10))
	}
	return phone
}
