package controller

import (
	"bme/internal/constants"
	"bme/internal/service"
	"github.com/go-playground/validator/v10"
)

type (
	DeviceEntities            []DeviceEntity
	DeviceErrorCreateRequests []DeviceErrorCreateRequest
	DeviceErrorEntities       []DeviceErrorEntity
)

type DeviceEntity struct {
	ID          uint                   `json:"id,omitempty"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Status      constants.DeviceStatus `json:"status,omitempty"`
}

type CreateDeviceRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	RequestedBy uint   `json:"-" validate:"required"`
}
type GetDeviceFilter struct {
	ID uint `uri:"id" binding:"required"`
}

type ListDevicesFilter struct {
	Q      *string                 `form:"q"`
	Status *constants.DeviceStatus `json:"status,omitempty"`
}

type ListDevicesResponse struct {
	Items DeviceEntities `json:"items"`
}

type DeviceErrorEntity struct {
	ID          uint   `json:"id"`
	DeviceID    uint   `json:"device_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type DeviceErrorCreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type DeviceErrorBulkCreateRequest struct {
	DeviceID    uint                      `uri:"id" binding:"required"`
	Items       DeviceErrorCreateRequests `json:"items"`
	RequestedBy uint                      `json:"-"`
}

type DeviceErrorListFilter struct {
	DeviceID uint                    `uri:"id" binding:"required"`
	Q        *string                 `form:"q"`
	Status   *constants.DeviceStatus `form:"status"`
}

type DeviceErrorListResponse struct {
	Items DeviceErrorEntities `json:"items"`
}

func (req CreateDeviceRequest) validate() error {
	validate := validator.New()

	if err := validate.Struct(&req); err != nil {
		return err
	}

	return nil
}

func (req CreateDeviceRequest) toSvc() service.CreateDeviceRequest {
	return service.CreateDeviceRequest{
		Title:       req.Title,
		Description: req.Description,
		RequestedBy: req.RequestedBy,
		Status:      constants.DeviceStatusDefaultOnCreation,
	}
}

func toViewDeviceEntity(svcEntity service.DeviceEntity) DeviceEntity {
	return DeviceEntity{
		ID:          svcEntity.ID,
		Title:       svcEntity.Title,
		Description: svcEntity.Description,
		Status:      svcEntity.Status,
	}
}

func (req GetDeviceFilter) toSvc() service.GetDeviceFilter {
	return service.GetDeviceFilter{
		ID: req.ID,
	}
}

func toViewDeviceEntities(svcEntities service.DeviceEntities) DeviceEntities {
	result := make(DeviceEntities, 0, len(svcEntities))

	for _, svcEntity := range svcEntities {
		result = append(result, toViewDeviceEntity(svcEntity))
	}

	return result
}

func toViewListDevicesResp(resp service.ListDevicesResponse) ListDevicesResponse {
	return ListDevicesResponse{
		Items: toViewDeviceEntities(resp.Items),
	}
}

func (f ListDevicesFilter) toSvc() service.ListDevicesFilter {
	idStartsWith, titleStartsWith := qAs(f.Q)

	return service.ListDevicesFilter{
		IdStartsWith:    idStartsWith,
		TitleStartsWith: titleStartsWith,
		Status:          nil,
	}
}

func (req DeviceErrorBulkCreateRequest) toSvc() service.DeviceErrorBulkCreateRequest {
	return service.DeviceErrorBulkCreateRequest{
		Entities:    req.Items.toSvc(req.DeviceID),
		RequestedBy: req.RequestedBy,
	}
}

func (requests DeviceErrorCreateRequests) toSvc(deviceID uint) service.DeviceErrorCreateRequests {
	result := make(service.DeviceErrorCreateRequests, 0, len(requests))

	for _, request := range requests {
		result = append(result, request.toSvc(deviceID))
	}

	return result
}

func (req DeviceErrorCreateRequest) toSvc(deviceID uint) service.DeviceErrorCreateRequest {
	return service.DeviceErrorCreateRequest{
		DeviceID:    deviceID,
		Title:       req.Title,
		Description: req.Description,
		Status:      constants.DeviceErrorStatusDefaultOnCreation,
	}
}

func (f DeviceErrorListFilter) toSvc() service.DeviceErrorListFilter {
	idStartsWith, titleStartsWith := qAs(f.Q)

	return service.DeviceErrorListFilter{
		DeviceID:        &f.DeviceID,
		IdStartsWith:    idStartsWith,
		TitleStartsWith: titleStartsWith,
		Status:          f.Status,
	}
}

func toViewListDeviceErrorsResp(response service.DeviceErrorListResponse) DeviceErrorListResponse {
	return DeviceErrorListResponse{
		Items: toViewDeviceErrorEntities(response.Entities),
	}
}

func toViewDeviceErrorEntities(svcEntities service.DeviceErrorEntities) DeviceErrorEntities {
	result := make(DeviceErrorEntities, 0, len(svcEntities))

	for _, svcEntity := range svcEntities {
		result = append(result, toViewDeviceErrorEntity(svcEntity))
	}

	return result
}

func toViewDeviceErrorEntity(svcEntity service.DeviceErrorEntity) DeviceErrorEntity {
	return DeviceErrorEntity{
		ID:          svcEntity.ID,
		DeviceID:    svcEntity.DeviceID,
		Title:       svcEntity.Title,
		Description: svcEntity.Description,
		Status:      svcEntity.Status,
	}
}
