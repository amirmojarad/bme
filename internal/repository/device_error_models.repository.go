package repository

import (
	"bme/internal/service"
	"gorm.io/gorm"
)

type DeviceErrorEntities []DeviceErrorEntity

type DeviceErrorEntity struct {
	gorm.Model
	DeviceID     uint
	DeviceEntity DeviceEntity `gorm:"foreignkey:DeviceID; references:ID"`
	Title        string
	Description  string
	Status       string
	CreatedBy    uint
	UpdatedBy    uint
	DeletedBy    uint
}

func (DeviceErrorEntity) TableName() string {
	return "device_errors"
}

func deviceErrorEntitiesFromSvcBulkCreateReq(req service.DeviceErrorBulkCreateRequest) DeviceErrorEntities {
	result := make(DeviceErrorEntities, 0)

	for _, entity := range req.Entities {
		result = append(result, deviceErrorEntityFromSvcCreateReq(entity, req.RequestedBy))
	}

	return result
}

func deviceErrorEntityFromSvcCreateReq(req service.DeviceErrorCreateRequest, requestedBy uint) DeviceErrorEntity {
	return DeviceErrorEntity{
		DeviceID:    req.DeviceID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status.String(),
		CreatedBy:   requestedBy,
		UpdatedBy:   requestedBy,
	}
}

func (entities DeviceErrorEntities) toSvc() service.DeviceErrorEntities {
	result := make(service.DeviceErrorEntities, 0, len(entities))

	for _, entity := range entities {
		result = append(result, entity.toSvc())
	}

	return result
}

func (entity DeviceErrorEntity) toSvc() service.DeviceErrorEntity {
	return service.DeviceErrorEntity{
		ID:           entity.ID,
		DeviceID:     entity.DeviceID,
		DeviceEntity: entity.DeviceEntity.toSvc(),
		Title:        entity.Title,
		Description:  entity.Description,
		Status:       entity.Status,
		CreatedBy:    entity.CreatedBy,
		UpdatedBy:    entity.CreatedBy,
		DeletedBy:    entity.DeletedBy,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
		DeletedAt:    entity.DeletedAt.Time,
	}
}
