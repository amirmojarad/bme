package middleware

import (
	"bme/pkg/errorext"
	"bme/pkg/jwt"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

const (
	AuthorizationHeader = "Authorization"

	AuthorizationUserID = "X-Auth-User-ID"
)

type Jwt interface {
	ValidateToken(tokenString string, secret []byte) (*jwt.UserClaims, error)
}
type Auth struct {
	jwt    Jwt
	secret []byte
}

func NewAuth(jwt Jwt, secret []byte) *Auth {
	return &Auth{
		jwt:    jwt,
		secret: secret,
	}
}

func (m *Auth) Authorize(ctx *gin.Context) {
	authorizationHeader := ctx.GetHeader(AuthorizationHeader)

	if authorizationHeader == "" {
		writeUnAuthorizedError(ctx, errorext.NewAuth(errors.New(fmt.Sprintf("%s header is empty", AuthorizationHeader)), errorext.ErrUnAuthorized))

		return
	}

	token := removeBearerFromAuthorizationHeader(authorizationHeader)

	userClaims, err := m.jwt.ValidateToken(token, m.secret)
	if err != nil {
		writeUnAuthorizedError(ctx, errorext.NewAuth(errors.New(fmt.Sprintf("%s is empty", AuthorizationHeader)), errorext.ErrUnAuthorized))

		return
	}

	ctx.Request.Header.Set(AuthorizationUserID, fmt.Sprintf("%d", userClaims.UserID))
	ctx.Request.Header.Del(AuthorizationHeader)
	ctx.Next()
}

func removeBearerFromAuthorizationHeader(authorizationHeader string) string {
	if token, founded := strings.CutPrefix(authorizationHeader, "Bearer "); founded {
		return token
	}

	return ""
}

func writeUnAuthorizedError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
}
