package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// []byte(register.Password) đổi kiểu "string" sang kiểu "mảng byte"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("có lỗi khi hash password")
	}
	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateToken(userId int, jwtSecret string) (string, error) {
	// dùng thuật toán  SHA256 để khởi tạo tokens
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"expr":   time.Now().Add(time.Hour * 72).Unix(), // 3 ngày
	})

	// ký token với khóa bí mật
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", errors.New("lỗi tạo token")

	}
	return tokenString, nil
}

func ValidateToken(tokenString, jwtSecret string) (*jwt.Token, error) {
	// kiểm tra tính hợp lệ của token với jwtSecret
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})
	return token, err
}

func GetClaims(token *jwt.Token) (jwt.MapClaims, bool) {
	// lấy thông tin (claims) đã ký vào token lúc tạo
	claims, ok := token.Claims.(jwt.MapClaims)
	return claims, ok
}
