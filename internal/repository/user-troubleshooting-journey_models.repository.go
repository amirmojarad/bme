package repository

import "time"

type UserTroubleshootingJourneyEntity struct {
	ID                           uint `gorm:"primarykey"`
	UserTroubleshootingSessionID uint
	Session                      UserTroubleshootingSessionEntity `gorm:"foreignkey:UserTroubleshootingSessionID; references:ID"`
	FromTroubleshootingSessionID uint
	ToTroubleshootingSessionID   uint
	Description                  string
	CreatedAt                    time.Time
}

func (UserTroubleshootingJourneyEntity) TableName() string {
	return "user_troubleshooting_journey"
}
