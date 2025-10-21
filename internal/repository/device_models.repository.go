package repository

import (
	"bme/internal/constants"
	"bme/internal/service"
	"gorm.io/gorm"
)

type DeviceEntities []DeviceEntity

type DeviceEntity struct {
	gorm.Model
	Title       string
	Description string
	Status      string
	CreatedBy   uint
	UpdatedBy   uint
	DeletedBy   uint
}

func (DeviceEntity) TableName() string {
	return "devices"
}

func deviceEntityFromSvcCreateDeviceRequest(req service.CreateDeviceRequest) DeviceEntity {
	return DeviceEntity{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status.String(),
		CreatedBy:   req.RequestedBy,
		UpdatedBy:   req.RequestedBy,
	}
}

func (entity DeviceEntity) toSvc() service.DeviceEntity {
	return service.DeviceEntity{
		ID:          entity.ID,
		Title:       entity.Title,
		Description: entity.Description,
		Status:      constants.DeviceStatus(entity.Status),
		CreatedBy:   entity.CreatedBy,
		UpdatedBy:   entity.UpdatedBy,
		DeletedBy:   entity.DeletedBy,
	}
}

func (entities DeviceEntities) toSvc() service.DeviceEntities {
	result := make(service.DeviceEntities, 0, len(entities))

	for _, entity := range entities {
		result = append(result, entity.toSvc())
	}

	return result
}
