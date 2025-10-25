package controller

import (
	"bme/internal/service"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserTroubleshootingService interface {
	UpdateStatus(ctx context.Context, req service.UserTroubleshootingUpdateStatusRequest) error
	CreateSession(ctx context.Context, req service.UserTroubleshootingSessionCreateRequest) (service.UserTroubleshootingSessionWithDetailsEntity, error)
	ListSessions(ctx context.Context, f service.UserTroubleshootingSessionListFilter) (service.UserTroubleshootingSessionsListWithDetailsResp, error)
	CurrentActiveSession(ctx context.Context, req service.UserTroubleshootingCurrentActiveSessionReq) (service.UserTroubleshootingSessionWithDetailsEntity, error)
	NextStep(ctx context.Context, req service.UserTroubleshootingNextStepRequest) error
	DeclineSession(ctx context.Context, userID uint) error
	DoneSession(ctx context.Context, userID uint) error
	PrevStep(ctx context.Context, req service.UserTroubleshootingPrevStepRequest) error
}

type UserTroubleshooting struct {
	svc    UserTroubleshootingService
	logger *logrus.Entry
}

func NewUserTroubleshooting(svc UserTroubleshootingService, logger *logrus.Entry) *UserTroubleshooting {
	return &UserTroubleshooting{
		svc:    svc,
		logger: logger,
	}
}

func (c UserTroubleshooting) CreateSession(ctx *gin.Context) {
	var (
		req    UserTroubleshootingSessionCreateRequest
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

	req.UserID = header.UserID

	sessionID, err := c.svc.CreateSession(ctx, req.toSvc())
	if err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"session_id": sessionID})
}

func (c UserTroubleshooting) ListSessions(ctx *gin.Context) {
	var (
		f      UserTroubleshootingSessionListFilter
		header HeaderEntityBindingRequired
	)

	if err := ctx.ShouldBindQuery(&f); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	if err := ctx.ShouldBindHeader(&header); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	resp, err := c.svc.ListSessions(ctx, f.toSvc())
	if err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	ctx.JSON(http.StatusOK, toViewUserTroubleshootingSessionsListWithDetailsResp(resp))
}

func (c UserTroubleshooting) UpdateStatus(ctx *gin.Context) {
	var (
		req    UserTroubleshootingUpdateStatusRequest
		header HeaderEntityBindingRequired
	)

	if err := ctx.ShouldBindUri(&req); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	if err := ctx.ShouldBindHeader(&header); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	req.UserID = header.UserID

	if err := req.validate(); err != nil {
		writeBadRequestErrorResponse(ctx, err)

		return
	}

	if err := c.svc.UpdateStatus(ctx, req.toSvc()); err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (c UserTroubleshooting) DeclineSession(ctx *gin.Context) {
	var (
		header HeaderEntityBindingRequired
	)

	if err := ctx.ShouldBindHeader(&header); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	if err := c.svc.DeclineSession(ctx, header.UserID); err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (c UserTroubleshooting) DoneSession(ctx *gin.Context) {
	var (
		header HeaderEntityBindingRequired
	)

	if err := ctx.ShouldBindHeader(&header); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	if err := c.svc.DoneSession(ctx, header.UserID); err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (c UserTroubleshooting) CurrentActiveSession(ctx *gin.Context) {
	var (
		header HeaderEntityBindingRequired
	)

	if err := ctx.ShouldBindHeader(&header); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	resp, err := c.svc.CurrentActiveSession(ctx, service.UserTroubleshootingCurrentActiveSessionReq{UserID: header.UserID})
	if err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	ctx.JSON(http.StatusOK, toViewUserTroubleshootingSessionWithDetailsEntityFromSvc(resp))
}

func (c UserTroubleshooting) NextStep(ctx *gin.Context) {
	var (
		header HeaderEntityBindingRequired
		req    UserTroubleshootingNextStepRequest
	)

	if err := ctx.ShouldBindHeader(&header); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	req.UserID = header.UserID

	if err := c.svc.NextStep(ctx, req.toSvc()); err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (c UserTroubleshooting) PrevStep(ctx *gin.Context) {
	var (
		header HeaderEntityBindingRequired
		req    UserTroubleshootingPrevStepRequest
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	if err := ctx.ShouldBindHeader(&header); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	req.UserID = header.UserID

	if err := c.svc.PrevStep(ctx, req.toSvc()); err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
