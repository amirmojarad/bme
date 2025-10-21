package jwt

import "github.com/golang-jwt/jwt/v4"

type Tokens struct {
	RefreshToken string
	AccessToken  string
}

type UserClaims struct {
	UserID   uint   `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
