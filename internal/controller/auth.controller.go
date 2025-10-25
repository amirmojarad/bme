package controller

import (
	"bme/internal/service"
	"bme/pkg/jwt"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AuthService interface {
	Register(ctx context.Context, req service.AuthRegisterRequest) (service.UserEntity, error)
	Login(ctx context.Context, req service.AuthLoginRequest) (service.AuthLoginResponse, error)
}

type Jwt interface {
	GenerateTokens(claims jwt.UserClaims) (jwt.Tokens, error)
	ValidateToken(tokenString string, secret []byte) (*jwt.UserClaims, error)
}

type Auth struct {
	svc    AuthService
	jwt    Jwt
	logger *logrus.Entry
}

func NewAuth(svc AuthService, logger *logrus.Entry, zeusJwt jwt.ZeusJwt) Auth {
	return Auth{
		svc:    svc,
		logger: logger,
		jwt:    zeusJwt,
	}
}

func (c Auth) Register(ctx *gin.Context) {
	var req UserRegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	validationErr, err := req.validate()
	if err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	if validationErr != nil {
		writeValidationErrors(ctx, validationErr)

		return
	}

	newUser, err := c.svc.Register(ctx, req.toSvc())
	if err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	tokens, err := c.jwt.GenerateTokens(newUser.UserClaims())
	if err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	ctx.JSON(http.StatusOK, tokensFromJwtTokens(tokens))
}

func (c Auth) Login(ctx *gin.Context) {
	var req UserLoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	validationErr, err := req.validate()
	if err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	if validationErr != nil {
		writeValidationErrors(ctx, validationErr)

		return
	}

	svcResp, err := c.svc.Login(ctx, req.toSvc())
	if err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	tokens, err := c.jwt.GenerateTokens(svcResp.JwtClaims())
	if err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	ctx.JSON(http.StatusOK, tokensFromJwtTokens(tokens))
}
