package repository

import (
	"bme/database"
	"bme/internal/service"
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type UserTroubleshootingJourney struct {
	*database.GormWrapper
}

func NewUserTroubleshootingJourney(db *database.GormWrapper) *UserTroubleshootingJourney {
	return &UserTroubleshootingJourney{
		db,
	}
}

func (r *UserTroubleshootingJourney) Create(ctx context.Context, req service.UserTroubleshootingJourneyCreateRequest) error {
	entity := userTroubleshootingJourneyEntityFromSvcReq(req)

	return errors.WithStack(r.DB(ctx).Model(&UserTroubleshootingJourneyEntity{}).Create(&entity).Error)
}

func (r *UserTroubleshootingJourney) Latest(ctx context.Context, sessionID uint) (service.UserTroubleshootingJourneyEntity, error) {
	var entity UserTroubleshootingJourneyEntity

	if err := r.DB(ctx).
		Model(&UserTroubleshootingJourneyEntity{}).
		Preload("FromTroubleshootingStep").
		Where("user_troubleshooting_session_id = ?", sessionID).
		Order("created_at desc").
		First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.UserTroubleshootingJourneyEntity{}, nil
		}

		return service.UserTroubleshootingJourneyEntity{}, errors.WithStack(err)
	}

	return entity.toSvc(), nil
}

func (r *UserTroubleshootingJourney) List(ctx context.Context, sessionID uint) (service.UserTroubleshootingJourneyEntities, error) {
	entities := make(UserTroubleshootingJourneyEntities, 0)

	if err := r.DB(ctx).Model(&UserTroubleshootingJourneyEntity{}).Where("user_troubleshooting_session_id = ?", sessionID).Order("created_at").Find(&entities).Error; err != nil {
		return service.UserTroubleshootingJourneyEntities{}, errors.WithStack(err)
	}

	return entities.toSvc(), nil
}

func (r *UserTroubleshootingJourney) Finish(ctx context.Context, id uint) error {
	return errors.WithStack(r.DB(ctx).Model(&UserTroubleshootingJourneyEntity{}).Where("id = ?", id).Update("finished_at", time.Now()).Error)
}
