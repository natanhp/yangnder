package auth

import "github.com/golang-jwt/jwt/v5"

var secretKey = []byte("secret")

func CreateToken(id uint) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
	})

	tokenString, err := claims.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
