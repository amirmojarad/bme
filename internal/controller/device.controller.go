package controller

import (
	"bme/internal/service"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type DeviceService interface {
	Create(ctx context.Context, req service.CreateDeviceRequest) error
	Get(ctx context.Context, f service.GetDeviceFilter) (service.DeviceEntity, error)
	List(ctx context.Context, f service.ListDevicesFilter) (service.ListDevicesResponse, error)
	BulkCreateDeviceErrors(ctx context.Context, req service.DeviceErrorBulkCreateRequest) error
	ListDeviceErrors(ctx context.Context, f service.DeviceErrorListFilter) (service.DeviceErrorListResponse, error)
	BulkCreateTroubleshootingSteps(ctx context.Context, req service.TroubleshootingBulkCreateRequest) error
	GetTroubleshootingStep(ctx context.Context, f service.TroubleshootingStepGetFilter) (service.TroubleshootingStepEntity, error)
	ListTroubleshootingSteps(ctx context.Context, f service.TroubleshootingStepListFilter) (service.TroubleshootingStepListResponse, error)
}

type Device struct {
	svc    DeviceService
	logger *logrus.Entry
}

func NewDevice(svc DeviceService, logger *logrus.Entry) *Device {
	return &Device{
		svc:    svc,
		logger: logger,
	}
}

func (c *Device) Create(ctx *gin.Context) {
	var (
		req    CreateDeviceRequest
		header HeaderEntityBindingRequired
	)

	if err := ctx.ShouldBindHeader(&header); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeBindingErrorResponse(ctx, err)
	}

	req.RequestedBy = header.UserID

	if err := req.validate(); err != nil {
		writeBadRequestErrorResponse(ctx, err)

		return
	}

	if err := c.svc.Create(ctx, req.toSvc()); err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{})
}

func (c *Device) Get(ctx *gin.Context) {
	var req GetDeviceFilter

	if err := ctx.ShouldBindUri(&req); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		writeBadRequestErrorResponse(ctx, err)

		return
	}

	resp, err := c.svc.Get(ctx, req.toSvc())
	if err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	ctx.JSON(http.StatusOK, toViewDeviceEntity(resp))
}

func (c *Device) List(ctx *gin.Context) {
	var f ListDevicesFilter

	if err := ctx.ShouldBindQuery(&f); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	resp, err := c.svc.List(ctx, f.toSvc())
	if err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	ctx.JSON(http.StatusOK, toViewListDevicesResp(resp))
}

func (c *Device) BulkCreateDeviceErrors(ctx *gin.Context) {
	var (
		req    DeviceErrorBulkCreateRequest
		header HeaderEntityBindingRequired
	)

	if err := ctx.ShouldBindUri(&req); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	if err := ctx.ShouldBindHeader(&header); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	req.RequestedBy = header.UserID

	if err := c.svc.BulkCreateDeviceErrors(ctx, req.toSvc()); err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{})
}

func (c *Device) ListDeviceErrors(ctx *gin.Context) {
	var f DeviceErrorListFilter

	if err := ctx.ShouldBindUri(&f); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	if err := ctx.ShouldBindQuery(&f); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	resp, err := c.svc.ListDeviceErrors(ctx, f.toSvc())
	if err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	ctx.JSON(http.StatusOK, toViewListDeviceErrorsResp(resp))
}

func (c *Device) BulkCreateTroubleshootingSteps(ctx *gin.Context) {
	var (
		req    TroubleshootingBulkCreateRequest
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

	req.RequestedBy = header.UserID

	if err := c.svc.BulkCreateTroubleshootingSteps(ctx, req.toSvc()); err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{})
}

func (c *Device) GetTroubleshootingStep(ctx *gin.Context) {
	var f TroubleshootingStepGetFilter

	if err := ctx.ShouldBindUri(&f); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	resp, err := c.svc.GetTroubleshootingStep(ctx, f.toSvc())
	if err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	ctx.JSON(http.StatusOK, toViewTroubleshootingStepEntity(resp))
}

func (c *Device) ListTroubleshootingSteps(ctx *gin.Context) {
	var f TroubleshootingStepsListFilter

	if err := ctx.ShouldBindUri(&f); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	if err := ctx.ShouldBindQuery(&f); err != nil {
		writeBindingErrorResponse(ctx, err)

		return
	}

	resp, err := c.svc.ListTroubleshootingSteps(ctx, f.toSvc())
	if err != nil {
		writeErrorResponse(ctx, err, c.logger)

		return
	}

	ctx.JSON(http.StatusOK, toViewTroubleshootingStepListResponse(resp))
}
