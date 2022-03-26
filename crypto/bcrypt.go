package crypto

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// BCryptPasswordEncoder

func BCryptPasswordEncoder(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(password, bcryptPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(bcryptPassword), []byte(password))
	return err == nil
}

func GetHashingCost(hashedPassword []byte) int {
	cost, _ := bcrypt.Cost(hashedPassword) // 为了简单忽略错误处理
	return cost
}

func PasswordHashingHandler(w http.ResponseWriter, r *http.Request) {
	password := "secret"
	hash, _ := BCryptPasswordEncoder(password) // 为了简单忽略错误处理

	fmt.Fprintln(w, "Password:", password)
	fmt.Fprintln(w, "Hash:    ", hash)

	match := CheckPassword(password, hash)
	fmt.Fprintln(w, "Match:   ", match)

	cost := GetHashingCost([]byte(hash))
	fmt.Fprintln(w, "Cost:    ", cost)
}
