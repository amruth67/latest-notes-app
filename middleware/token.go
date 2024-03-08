package middleware

import (
	"fmt"
	"mongodb/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("Shilpa67")

var user models.User

func GenerateToken(username string) (tokenString string, msg string) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", "nil"
	}

	return tokenString, "nil"
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil
}
