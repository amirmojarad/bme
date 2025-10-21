package service

import "bme/internal/constants"

type DeviceEntities []DeviceEntity

type DeviceEntity struct {
	ID          uint
	Title       string
	Description string
	Status      constants.DeviceStatus
	CreatedBy   uint
	UpdatedBy   uint
	DeletedBy   uint
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
