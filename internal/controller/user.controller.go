package controller

import (
	"bme/internal/service"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserService interface {
	ResetPassword(ctx context.Context, req service.UserResetPasswordRequest) error
	Update(ctx context.Context, req service.UserUpdateRequest) error
	First(ctx context.Context, f service.FirstUserFilter) (service.UserEntity, error)
}

type User struct {
	svc    UserService
	logger *logrus.Entry
}

func NewUser(svc UserService, logger *logrus.Entry) *User {
	return &User{
		svc:    svc,
		logger: logger,
	}
}

func (c *User) ResetPassword(ctx *gin.Context) {
	var (
		req    UserResetPasswordRequest
		header HeaderEntityBindingRequired
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeBadRequestErrorResponse(ctx, err)

		return
	}

	if err := ctx.ShouldBindHeader(&header); err != nil {
		writeBadRequestErrorResponse(ctx, err)

		return
	}

	req.RequestedBy = header.UserID

	if err := c.svc.ResetPassword(ctx, req.toSvc()); err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (c *User) Update(ctx *gin.Context) {
	var (
		req    UserUpdateRequest
		header HeaderEntityBindingRequired
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	if err := ctx.ShouldBindHeader(&header); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	req.RequestedBy = header.UserID

	if err := c.svc.Update(ctx, req.toSvc()); err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (c *User) Get(ctx *gin.Context) {
	var (
		header HeaderEntityBindingRequired
	)

	if err := ctx.ShouldBindHeader(&header); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	resp, err := c.svc.First(ctx, service.FirstUserFilter{ID: &header.UserID})
	if err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	ctx.JSON(http.StatusOK, toViewUserResponse(resp))

}
