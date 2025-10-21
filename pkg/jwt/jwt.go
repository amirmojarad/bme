package jwt

import (
	"bme/conf"
	"bme/pkg/errorext"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type ZeusJwt struct {
	cfg *conf.AppConfig
}

func New(cfg *conf.AppConfig) ZeusJwt {
	return ZeusJwt{
		cfg: cfg,
	}
}

func (z ZeusJwt) accessTokenClaims(claims UserClaims) UserClaims {
	return UserClaims{
		UserID:   claims.UserID,
		Username: claims.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(z.cfg.Jwt.AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    z.cfg.Jwt.Issuer,
			Subject:   fmt.Sprintf("%d", claims.UserID),
		},
	}
}

func (z ZeusJwt) refreshTokenClaims(claims UserClaims) UserClaims {
	return UserClaims{
		UserID:   claims.UserID,
		Username: claims.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(z.cfg.Jwt.RefreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    z.cfg.Jwt.Issuer,
			Subject:   fmt.Sprintf("%d", claims.UserID),
		},
	}
}

// GenerateTokens creates access and refresh tokens with custom claims
func (z ZeusJwt) GenerateTokens(claims UserClaims) (Tokens, error) {
	accessClaims := z.accessTokenClaims(claims)

	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := accessJWT.SignedString(z.cfg.Jwt.AccessSecret)

	if err != nil {
		return Tokens{}, errorext.New(err, errorext.ErrGeneralOccurrence)
	}

	// Refresh token (3 days expiry)
	refreshClaims := z.refreshTokenClaims(claims)

	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refreshJWT.SignedString(z.cfg.Jwt.RefreshSecret)
	if err != nil {
		return Tokens{}, errorext.New(err, errorext.ErrGeneralOccurrence)
	}

	return Tokens{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}

// ValidateToken verifies a token and returns its claims
func (z ZeusJwt) ValidateToken(tokenString string, secret []byte) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
