package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// hash password waktu register
func (u *Customer) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}

// DTO
type LoginRequestCustomer struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// check password
func (u *Customer) CheckPasswordCustomer(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
