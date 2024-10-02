package helper

import (
	"fmt"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
)

func JwtParse(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		fmt.Println("Hai")
		// Return the secret key
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	return token, err

}
