package repository

import (
	"bme/internal/constants"
	"bme/internal/service"
	"gorm.io/gorm"
)

type (
	TroubleshootingStepsToStepsEntities []TroubleshootingStepsToStepsEntity
)

type TroubleshootingStepsToStepsEntity struct {
	gorm.Model
	FromStepID                    uint
	FromTroubleshootingStepEntity TroubleshootingStepEntity `gorm:"foreignkey:FromStepID" references:"ID"`
	ToStepID                      uint
	ToTroubleshootingStepEntity   TroubleshootingStepEntity `gorm:"foreignkey:ToStepID" references:"ID"`
	Priority                      constants.TroubleshootingStepsToStepsPriority
	CreatedBy                     uint
	UpdatedBy                     uint
	DeletedBy                     uint
}

func (TroubleshootingStepsToStepsEntity) TableName() string {
	return "troubleshooting_steps_to_steps"
}

func troubleshootingStepsToStepsEntityFromSvcReq(req service.TroubleshootingStepsToStepsCreateEntity) TroubleshootingStepsToStepsEntity {
	return TroubleshootingStepsToStepsEntity{
		FromStepID: req.FromStepID,
		ToStepID:   req.ToStepID,
		Priority:   req.Priority,
		CreatedBy:  req.CreatedBy,
		UpdatedBy:  req.UpdatedBy,
	}
}

func troubleshootingStepsToStepsEntitiesFromSvcReq(req service.TroubleshootingStepsToStepsCreateReq) TroubleshootingStepsToStepsEntities {
	result := make(TroubleshootingStepsToStepsEntities, 0)

	for _, entity := range req.Entities {
		result = append(result, troubleshootingStepsToStepsEntityFromSvcReq(entity))
	}

	return result
}

func (entities TroubleshootingStepsToStepsEntities) toSvc() service.TroubleshootingStepsToStepsWithDetailsEntities {
	result := make(service.TroubleshootingStepsToStepsWithDetailsEntities, 0, len(entities))

	for _, entity := range entities {
		result = append(result, entity.toSvc())
	}

	return result
}

func (entity TroubleshootingStepsToStepsEntity) toSvc() service.TroubleshootingStepsToStepsWithDetailsEntity {
	return service.TroubleshootingStepsToStepsWithDetailsEntity{
		ID:                            entity.ID,
		FromStepID:                    entity.FromStepID,
		FromTroubleshootingStepEntity: entity.FromTroubleshootingStepEntity.toSvc(),
		ToStepID:                      entity.ToStepID,
		ToTroubleshootingStepEntity:   entity.ToTroubleshootingStepEntity.toSvc(),
		Priority:                      entity.Priority.OrDefault(),
		PriorityTitle:                 entity.Priority.String(),
		CreatedBy:                     entity.CreatedBy,
		UpdatedBy:                     entity.UpdatedBy,
		DeletedBy:                     entity.DeletedBy,
		CreatedAt:                     entity.CreatedAt,
		UpdatedAt:                     entity.UpdatedAt,
		DeletedAt:                     &entity.DeletedAt.Time,
	}
}

func (entities TroubleshootingStepsToStepsEntities) toSvcTroubleshootingEntities() service.TroubleshootingStepEntities {
	result := make(service.TroubleshootingStepEntities, 0, len(entities))

	for _, entity := range entities {
		result = append(result, entity.ToTroubleshootingStepEntity.toSvc())
	}

	return result
}
