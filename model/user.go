package models

import (
	"fmt"
	"grpcapp/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"-"`
	ID         uint       `gorm:"primaryKey;<-:create" json:"id"`
	Name       string     `json:"name"`
	Email      string     `json:"email" gorm:"unique;not null"`
	Gender     string     `json:"gender"`
	Address    string     `json:"address"`
	Username   string     `json:"username" gorm:"unique;not null"`
	Password   string     `json:"password" gorm:"not null"`
	UpdatedAt  *time.Time `json:"updated_at" gorm:"default:null"`
}

var config, _ = utils.LoadConfig(".")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

func (u *User) CreateJWT() (string, error) {
	expirationTime := time.Now().Add(24000 * time.Hour)
	claims := &Claims{
		Email:    u.Email,
		ID:       u.ID,
		Username: u.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	return tokenString, err
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	return hookTask(user, tx)
}

func (user *User) BeforeUpdate(tx *gorm.DB) (err error) {
	return hookTask(user, tx)
}

func hookTask(user *User, tx *gorm.DB) error {
	today := time.Now()
	user.UpdatedAt = &today
	// Check if password is modified and has a value
	if user.Password != "" {
		fmt.Println("password changed")
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
		tx.Statement.SetColumn("password", user.Password)
	}

	if tx.Statement.Changed() {
		tx.Statement.SetColumn("updated_at", user.UpdatedAt)
	}
	return nil
}
