package jwt

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"go.uber.org/fx"
)

// JWT defines one struct
type JWT struct {
	// SigningKey signature data
	SigningKey []byte
}

// NewJWT initialize for DI
func NewJWT() *JWT {
	return &JWT{
		[]byte("Schedrestd --->>>>"),
	}
}

// AIPClaims ...
type AIPClaims struct {
	Name  string `json:"name"`
	jwt.StandardClaims
}


// GenerateToken generates the token and use jwt.SigningMethodHS256
func (j *JWT) GenerateToken(claims AIPClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken parses the jwt token
func (j *JWT) ParseToken(tokenString string) (*AIPClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AIPClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors & jwt.ValidationErrorMalformed != 0 {
				return nil, fmt.Errorf("token unvailable")
			} else if ve.Errors & jwt.ValidationErrorExpired != 0 {
				return nil, fmt.Errorf("token expired")
			} else if ve.Errors & jwt.ValidationErrorNotValidYet != 0 {
				return nil, fmt.Errorf("token invalid")
			} else {
				return nil, fmt.Errorf("token unvailable")
			}

		}
	}

	if claims, ok := token.Claims.(*AIPClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("token invalid")

}

// Module ...
var Module = fx.Options(fx.Provide(NewJWT))
