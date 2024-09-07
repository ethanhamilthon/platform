package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (a *AuthService) GenerateJWT(hours int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims.(jwt.MapClaims)["exp"] = time.Now().Add(time.Hour * time.Duration(hours)).Unix()
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *AuthService) VerifyJWT(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}
	return fmt.Errorf("invalid token")
}
