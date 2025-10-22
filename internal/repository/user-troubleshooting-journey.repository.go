package repository

import "bme/database"

type UserTroubleshootingJourney struct {
	*database.GormWrapper
}

func NewUserTroubleshootingJourney(db *database.GormWrapper) *UserTroubleshootingJourney {
	return &UserTroubleshootingJourney{
		db,
	}
}
