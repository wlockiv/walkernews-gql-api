package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/wlockiv/walkernews/graph/model"
	"time"
)

var (
	SecretKey = []byte("secret")
)

// Generates a JWT token.
//   - email:  user's email address
//   - userKey: FaunaDB Key from the Login Query
func GenerateToken(user *model.User, userKey string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = user.ID
	claims["email"] = user.Email
	claims["userKey"] = userKey
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(SecretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenStr string) (map[string]string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// ! Is there a better way of doing this?
		results := map[string]string{
			"userId":  claims["userId"].(string),
			"email":   claims["email"].(string),
			"userKey": claims["userKey"].(string),
		}

		return results, nil
	} else {
		err := errors.New("their was either an error parsing the token or it was not valid")
		return nil, err
	}
}
