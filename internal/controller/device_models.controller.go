package controller

import (
	"bme/internal/constants"
	"bme/internal/service"
	"github.com/go-playground/validator/v10"
	"strconv"
)

type DeviceEntities []DeviceEntity

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
	var (
		idStartsWith    *string
		titleStartsWith *string
	)

	numericVal, q := convertQ[uint](f.Q)
	if numericVal != nil {
		tmp := strconv.Itoa(int(*numericVal))
		idStartsWith = &tmp
	}

	if q != nil {
		titleStartsWith = q
	}

	return service.ListDevicesFilter{
		IdStartsWith:    idStartsWith,
		TitleStartsWith: titleStartsWith,
		Status:          nil,
	}
}
