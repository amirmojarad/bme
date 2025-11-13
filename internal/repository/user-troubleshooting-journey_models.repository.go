package repository

import (
	"bme/internal/service"
	"time"
)

type UserTroubleshootingJourneyEntities []UserTroubleshootingJourneyEntity

type UserTroubleshootingJourneyEntity struct {
	ID                           uint `gorm:"primarykey"`
	UserTroubleshootingSessionID uint
	Session                      UserTroubleshootingSessionEntity `gorm:"foreignkey:UserTroubleshootingSessionID; references:ID"`
	FromTroubleshootingStepID    uint
	FromTroubleshootingStep      TroubleshootingStepEntity `gorm:"foreignkey:FromTroubleshootingStepID; references:ID"`
	ToTroubleshootingStepID      uint
	Description                  string
	CreatedAt                    time.Time
	FinishedAt                   *time.Time
}

func (UserTroubleshootingJourneyEntity) TableName() string {
	return "user_troubleshooting_journey"
}

func userTroubleshootingJourneyEntityFromSvcReq(req service.UserTroubleshootingJourneyCreateRequest) UserTroubleshootingJourneyEntity {
	return UserTroubleshootingJourneyEntity{
		UserTroubleshootingSessionID: req.SessionID,
		FromTroubleshootingStepID:    req.FromTroubleshootingStepID,
		ToTroubleshootingStepID:      req.ToTroubleshootingStepID,
	}
}

func (entity UserTroubleshootingJourneyEntity) toSvc() service.UserTroubleshootingJourneyEntity {
	return service.UserTroubleshootingJourneyEntity{
		ID:                           entity.ID,
		UserTroubleshootingSessionID: entity.UserTroubleshootingSessionID,
		FromTroubleshootingStepID:    entity.FromTroubleshootingStepID,
		ToTroubleshootingStepID:      entity.ToTroubleshootingStepID,
		Description:                  entity.Description,
		CreatedAt:                    entity.CreatedAt,
		FromTroubleshootingStepTitle: entity.FromTroubleshootingStep.Title,
		FinishedAt:                   entity.FinishedAt,
	}
}

func (entities UserTroubleshootingJourneyEntities) toSvc() service.UserTroubleshootingJourneyEntities {
	result := make(service.UserTroubleshootingJourneyEntities, 0, len(entities))

	for _, entity := range entities {
		result = append(result, entity.toSvc())
	}

	return result
}
