package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Phone string `json:"phone"`
	jwt.StandardClaims
}

// HashPassword хеширует пароль с использованием bcrypt.
// @param password string Исходный пароль для хеширования.
// @return string Захешированный пароль.
// @return error Ошибка, если произошла ошибка хеширования.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash проверяет соответствие пароля и его хеша с использованием bcrypt.
// @param password string Исходный пароль для проверки.
// @param hash string Захешированный пароль для сравнения.
// @return bool Результат проверки: true - если пароль совпадает с хешем, false - в противном случае.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT генерирует JWT на основе телефонного номера.
// @param phone string Телефонный номер пользователя.
// @return string Сгенерированный JWT токен.
// @return error Ошибка, если не удалось сгенерировать токен.
func GenerateJWT(phone string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Phone: phone,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
