package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	SecretKey = []byte("secret")
)

func NewToken(id, username string) (string, error) {
	token := jwt.New(jwt.SigningMethodES256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = id
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(SecretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := claims["userId"].(string)
		return userId, nil
	} else {
		err := errors.New("their was either an error parsing the token or it was not valid")
		return "", err
	}
}
