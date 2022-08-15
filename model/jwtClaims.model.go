package model

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	User_Id int    `json:"user_id"`
	Email   string `json:"user_email"`
	Role    string `json:""`
	jwt.StandardClaims
}

func (claims JWTClaims) Valid() error {
	var now = time.Now().UTC().Unix()
	if claims.VerifyExpiresAt(now, true) {
		return nil
	}
	return fmt.Errorf("Invalid token")
}
