package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(password string) (pwd string) {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	pwd = string(bytes)
	return
}
func CheckPassword(hashedpwd string, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedpwd), []byte(pwd))
	return err == nil
}
