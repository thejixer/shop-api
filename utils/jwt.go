package utils

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

func SignToken(id int) (string, error) {

	secret := os.Getenv("JWT_SECRET")
	mySigningKey := []byte(secret)

	claims := jwt.MapClaims{
		"id": fmt.Sprintf("%v", id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)

}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}
