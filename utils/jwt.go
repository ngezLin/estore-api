package utils

import (
	"estore-api/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokenAdmin(adminId uint) (string, error) {
	claims := jwt.MapClaims{
		"admin_id": adminId,
		"exp":      time.Now().Add(config.GetJWTExpirationDuration()).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString((config.GetJWTSecret()))
}

func GenerateTokenCustomer(customerId uint) (string, error) {
	claims := jwt.MapClaims{
		"customer_id": customerId,
		"exp":         time.Now().Add(config.GetJWTExpirationDuration()).Unix(),
		"iat":         time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString((config.GetJWTSecret()))
}

func ValidateTokenAdmin(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.GetJWTSecret(), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		adminId := uint(claims["admin_id"].(float64))
		return adminId, nil
	}

	return 0, err
}

func ValidateTokenCustomer(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.GetJWTSecret(), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		customerId := uint(claims["customer_id"].(float64))
		return customerId, nil
	}

	return 0, err
}
