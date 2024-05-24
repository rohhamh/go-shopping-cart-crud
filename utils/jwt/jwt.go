package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rohhamh/go-shopping-cart-crud/config"
)

func SignToken(claim jwt.Claims) string {
    t := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)
    s, err := t.SignedString(config.JWTKey)
    if err != nil {
        fmt.Printf("jwt got error %v\n", err)
    }
    return s
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte(config.JWTKey), nil
    })
    return token, err
}
