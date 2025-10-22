package repository

import (
	"bme/internal/constants"
	"bme/internal/service"
	"gorm.io/gorm"
	"time"
)

type (
	UserTroubleshootingSessionEntities []UserTroubleshootingSessionEntity
)

type UserTroubleshootingSessionEntity struct {
	ID            uint `gorm:"primarykey"`
	UserID        uint
	User          UserEntity `gorm:"foreignkey:UserID; references:ID"`
	DeviceID      uint
	Device        DeviceEntity `gorm:"foreignkey:DeviceID; references:ID"`
	DeviceErrorID uint
	DeviceError   DeviceErrorEntity `gorm:"foreignkey:DeviceErrorID; references:ID"`
	Status        string
	CreatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (UserTroubleshootingSessionEntity) TableName() string {
	return "user_troubleshooting_sessions"
}

func userTroubleshootingSessionEntityFromSvcReq(
	req service.UserTroubleshootingSessionCreateRequest,
) UserTroubleshootingSessionEntity {
	return UserTroubleshootingSessionEntity{
		UserID:        req.UserID,
		DeviceID:      req.DeviceID,
		DeviceErrorID: req.DeviceErrorID,
		Status:        constants.UserTroubleshootingSessionDefaultOnCreation.String(),
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
		ID:               entity.ID,
		UserID:           entity.UserID,
		DeviceID:         entity.DeviceID,
		DeviceTitle:      entity.Device.Title,
		DeviceErrorID:    entity.DeviceErrorID,
		DeviceErrorTitle: entity.DeviceError.Title,
		Status:           constants.UserTroubleshootingSessionsStatus(entity.Status),
		CreatedAt:        entity.CreatedAt,
		DeletedAt:        &entity.DeletedAt.Time,
	}

}

func (entities UserTroubleshootingSessionEntities) toSvc() service.UserTroubleshootingSessionWithDetailsEntities {
	result := make([]service.UserTroubleshootingSessionWithDetailsEntity, 0, len(entities))

	for _, entity := range entities {
		result = append(result, entity.toSvcWithDetails())
	}

	return result
}
