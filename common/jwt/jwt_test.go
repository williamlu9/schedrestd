package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
)

func TestJWT_GenerateToken(t *testing.T) {
	t.Run("Test jwt generate token", func(t *testing.T) {
		gJwt := NewJWT()
		claims := AIPClaims{
			"ppju",
			jwt.StandardClaims{
				//ExpiresAt: 0,
				Issuer: "Schedrestd",
			},
		}
		token, _ := gJwt.GenerateToken(claims)

		t.Log(token)
	})
}

func TestJWT_ParseToken(t *testing.T) {
	t.Run("Test parsing jwt token", func(t *testing.T) {
		gJwt := NewJWT()

		claims, err := gJwt.ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoicHBqdSIsImlzcyI6IlNreUZvcm0gQUlQIn0.6b6oxkf3hiBl_MtCoBBCpGi3dQV0kDzmRUa2gFXMHHs")
		if err != nil {
			t.Errorf(err.Error())
		} else {
			t.Log(claims)
		}
	})
}
