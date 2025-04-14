package middleware

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type JWTObj struct {
	Secret []byte
}

func (j *JWTObj) ValidateJWT(tokenStr string) (string, error) {
	errToken := errors.New("invalid token")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errToken
		}
		return j.Secret, nil
	})

	if err != nil || !token.Valid {
		return "", errToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errToken
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return "", errToken
	}

	return sub, nil
}
