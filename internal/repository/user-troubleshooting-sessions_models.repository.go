package repository

import (
	"bme/internal/constants"
	"bme/internal/service"
	"gorm.io/gorm"
	"time"
)

type (
	UserTroubleshootingSessionEntities   []UserTroubleshootingSessionEntity
	TroubleshootingStepTransitions       []TroubleshootingStepTransition
	TroubleshootingStepPrevStepsEntities []TroubleshootingStepPrevStepsEntity
)

type UserTroubleshootingSessionEntity struct {
	ID                           uint `gorm:"primarykey"`
	UserID                       uint
	User                         UserEntity `gorm:"foreignkey:UserID; references:ID"`
	DeviceID                     uint
	Device                       DeviceEntity `gorm:"foreignkey:DeviceID; references:ID"`
	DeviceErrorID                uint
	DeviceError                  DeviceErrorEntity `gorm:"foreignkey:DeviceErrorID; references:ID"`
	Status                       string
	CurrentTroubleshootingStepID *uint
	CurrentTroubleshootingStep   TroubleshootingStepEntity `gorm:"foreignkey:CurrentTroubleshootingStepID; references:ID"`
	CreatedAt                    time.Time
	DeletedAt                    gorm.DeletedAt `gorm:"index"`
}

type TroubleshootingStepTransition struct {
	FromStepID    uint   `gorm:"column:from_step_id"`
	FromStepTitle string `gorm:"column:title"`
	ToStepID      uint   `gorm:"column:to_step_id"`
	ToStepTitle   string `gorm:"column:to_title"`
}

type TroubleshootingStepPrevStepsEntity struct {
	FromStepID    uint   `gorm:"column:to_step_id"`
	FromStepTitle string `gorm:"column:to_step_title"`
	ToStepID      uint   `gorm:"column:from_step_id"`
	ToStepTitle   string `gorm:"column:from_step_title"`
}

func (UserTroubleshootingSessionEntity) TableName() string {
	return "user_troubleshooting_sessions"
}

func userTroubleshootingSessionEntityFromSvcReq(
	req service.UserTroubleshootingSessionCreateRequest,
) UserTroubleshootingSessionEntity {
	return UserTroubleshootingSessionEntity{
		UserID:                       req.UserID,
		DeviceID:                     req.DeviceID,
		DeviceErrorID:                req.DeviceErrorID,
		CurrentTroubleshootingStepID: req.StartStepID,
		Status:                       constants.UserTroubleshootingSessionDefaultOnCreation.String(),
	}
}

func (entity UserTroubleshootingSessionEntity) toSvc() service.UserTroubleshootingSessionEntity {
	return service.UserTroubleshootingSessionEntity{
		ID:            entity.ID,
		UserID:        entity.UserID,
		DeviceID:      entity.DeviceID,
		DeviceErrorID: entity.DeviceErrorID,
		Status:        constants.UserTroubleshootingSessionsStatus(entity.Status),
		CreatedAt:     entity.CreatedAt,
		DeletedAt:     &entity.DeletedAt.Time,
	}
}

func (entity UserTroubleshootingSessionEntity) toSvcWithDetails() service.UserTroubleshootingSessionWithDetailsEntity {
	return service.UserTroubleshootingSessionWithDetailsEntity{
		ID:                         entity.ID,
		UserID:                     entity.UserID,
		DeviceID:                   entity.DeviceID,
		DeviceTitle:                entity.Device.Title,
		DeviceErrorID:              entity.DeviceErrorID,
		DeviceErrorTitle:           entity.DeviceError.Title,
		CurrentTroubleshootingStep: entity.CurrentTroubleshootingStep.toSvc(),
		Status:                     constants.UserTroubleshootingSessionsStatus(entity.Status),
		CreatedAt:                  entity.CreatedAt,
		DeletedAt:                  &entity.DeletedAt.Time,
		CurrentStepID:              entity.CurrentTroubleshootingStepID,
	}

}

func (entities UserTroubleshootingSessionEntities) toSvc() service.UserTroubleshootingSessionWithDetailsEntities {
	result := make([]service.UserTroubleshootingSessionWithDetailsEntity, 0, len(entities))

	for _, entity := range entities {
		result = append(result, entity.toSvcWithDetails())
	}

	return result
}

func (transitions TroubleshootingStepTransitions) toSvcTroubleshootingStepsMap() service.TroubleshootingNextStepsMap {
	transitionsMap := make(map[uint]map[uint]string)

	for _, transition := range transitions {
		if _, ok := transitionsMap[transition.FromStepID]; !ok {
			transitionsMap[transition.FromStepID] = make(map[uint]string)
		}

		transitionsMap[transition.FromStepID][transition.ToStepID] = transition.ToStepTitle
	}

	return service.TroubleshootingNextStepsMap{
		Map: transitionsMap,
	}
}

func (transitions TroubleshootingStepPrevStepsEntities) toSvcTroubleshootingStepsMap() service.TroubleshootingNextStepsMap {
	transitionsMap := make(map[uint]map[uint]string)

	for _, transition := range transitions {
		if _, ok := transitionsMap[transition.FromStepID]; !ok {
			transitionsMap[transition.FromStepID] = make(map[uint]string)
		}

		transitionsMap[transition.FromStepID][transition.ToStepID] = transition.ToStepTitle
	}

	return service.TroubleshootingNextStepsMap{
		Map: transitionsMap,
	}
}
