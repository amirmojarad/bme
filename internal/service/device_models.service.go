package service

import (
	"bme/internal/constants"
	"time"
)

type (
	DeviceEntities            []DeviceEntity
	DeviceErrorCreateRequests []DeviceErrorCreateRequest
	DeviceErrorEntities       []DeviceErrorEntity
)

type DeviceEntity struct {
	ID                  uint
	DeviceErrorEntities DeviceErrorEntities
	Title               string
	Description         string
	Status              constants.DeviceStatus
	CreatedBy           uint
	UpdatedBy           uint
	DeletedBy           uint
}

type CreateDeviceRequest struct {
	Title       string
	Description string
	Status      constants.DeviceStatus
	RequestedBy uint
}
type GetDeviceFilter struct {
	ID uint
}

type ListDevicesFilter struct {
	IdStartsWith    *string
	TitleStartsWith *string
	Status          *constants.DeviceStatus
}

type ListDevicesResponse struct {
	Items DeviceEntities
}

type DeviceErrorEntity struct {
	ID           uint
	DeviceID     uint
	DeviceEntity DeviceEntity
	Title        string
	Description  string
	Status       string
	CreatedBy    uint
	UpdatedBy    uint
	DeletedBy    uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}

type DeviceErrorCreateRequest struct {
	DeviceID    uint
	Title       string
	Description string
	Status      constants.DeviceStatus
}

type DeviceErrorBulkCreateRequest struct {
	Entities    []DeviceErrorCreateRequest
	RequestedBy uint
}

type DeviceErrorListFilter struct {
	DeviceID        *uint
	IdStartsWith    *string
	IDs             []uint
	TitleStartsWith *string
	WithDetails     bool
	Status          *constants.DeviceStatus
}

type DeviceErrorListResponse struct {
	Entities DeviceErrorEntities
}

func (f GetDeviceFilter) FilterMap() map[string]any {
	return map[string]any{
		"id": f.ID,
	}
}

func (f ListDevicesFilter) FilterMap() map[string]any {
	filterMap := make(map[string]any)

	if f.Status != nil {
		filterMap["status"] = f.Status.String()
	}

	return filterMap
}

func (f DeviceErrorListFilter) FilterMap() map[string]any {
	filterMap := make(map[string]any)

	if f.IDs != nil {
		filterMap["id"] = f.IDs
	}

	if f.DeviceID != nil {
		filterMap["device_id"] = f.DeviceID
	}

	return filterMap
}
