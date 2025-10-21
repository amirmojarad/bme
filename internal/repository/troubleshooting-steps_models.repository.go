package repository

import (
	"bme/internal/constants"
	"bme/internal/service"
	"gorm.io/gorm"
)

type (
	TroubleShootingStepEntities []TroubleshootingStepEntity
)

type TroubleshootingStepEntity struct {
	gorm.Model
	DeviceID          uint
	DeviceEntity      DeviceEntity `gorm:"foreignkey:DeviceID" reference:"ID"`
	DeviceErrorID     uint
	DeviceErrorEntity DeviceErrorEntity                   `gorm:"foreignKey:DeviceErrorID" reference:"ID"`
	NextSteps         TroubleshootingStepsToStepsEntities `gorm:"foreignKey:FromStepID;references:ID"`
	Title             string
	Description       string
	Hints             map[string]any `gorm:"serializer:json"`
	Status            string
	CreatedBy         uint
	UpdatedBy         uint
	DeletedBy         uint
}

func (TroubleshootingStepEntity) TableName() string {
	return "troubleshooting_steps"
}

func troubleshootingEntitiesFromSvcTroubleShootingBulkCreateRequest(
	svcReq service.TroubleshootingBulkCreateRequest,
) TroubleShootingStepEntities {
	entities := make(TroubleShootingStepEntities, 0, len(svcReq.Entities))

	for _, entity := range svcReq.Entities {
		entities = append(entities, troubleShootingEntityFromSvcTroubleShootingBulkCreateRequest(entity, svcReq.RequestedBy))
	}

	return entities
}

func troubleShootingEntityFromSvcTroubleShootingBulkCreateRequest(svcCreateEntity service.TroubleShootingStepCreateEntity, requestedBy uint) TroubleshootingStepEntity {
	return TroubleshootingStepEntity{
		DeviceID:      svcCreateEntity.DeviceID,
		DeviceErrorID: svcCreateEntity.DeviceErrorID,
		Title:         svcCreateEntity.Title,
		Description:   svcCreateEntity.Description,
		Hints:         svcCreateEntity.Hints,
		Status:        svcCreateEntity.Status.String(),
		CreatedBy:     requestedBy,
		UpdatedBy:     requestedBy,
	}
}

func (entity TroubleshootingStepEntity) toSvc() service.TroubleshootingStepEntity {
	return service.TroubleshootingStepEntity{
		ID:            entity.ID,
		DeviceID:      entity.DeviceID,
		DeviceErrorID: entity.DeviceErrorID,
		Title:         entity.Title,
		Description:   entity.Description,
		Hints:         entity.Hints,
		Status:        constants.TroubleshootingStepStatus(entity.Status),
		CreatedBy:     entity.CreatedBy,
		UpdatedBy:     entity.CreatedBy,
		DeletedBy:     entity.DeletedBy,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		DeletedAt:     &entity.DeletedAt.Time,
		NextSteps:     entity.NextSteps.toSvcTroubleshootingEntities(),
	}
}

func (entities TroubleShootingStepEntities) toSvc() service.TroubleshootingStepEntities {
	result := make(service.TroubleshootingStepEntities, 0, len(entities))

	for _, entity := range entities {
		result = append(result, entity.toSvc())
	}

	return result
}
