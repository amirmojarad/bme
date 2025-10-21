package repository

import (
	"bme/database"
	"bme/internal/service"
	"bme/pkg/errorext"
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TroubleshootingSteps struct {
	*database.GormWrapper
}

func NewTroubleshootingSteps(db *database.GormWrapper) *TroubleshootingSteps {
	return &TroubleshootingSteps{
		db,
	}
}

func (r *TroubleshootingSteps) BulkCreate(ctx context.Context, req service.TroubleshootingBulkCreateRequest) error {
	entities := troubleshootingEntitiesFromSvcTroubleShootingBulkCreateRequest(req)

	return errors.WithStack(r.DB(ctx).Create(&entities).Error)
}

func (r *TroubleshootingSteps) List(ctx context.Context, f service.TroubleshootingStepListFilter) (service.TroubleshootingStepListResponse, error) {
	entities := make(TroubleShootingStepEntities, 0)

	query := r.DB(ctx).Where(f.FilterMap())

	if f.IdStartsWith != nil {
		appendAppendStartsWith(query, "id::text", *f.IdStartsWith, false)
	}

	if f.TitleStartsWith != nil {
		appendAppendStartsWith(query, "title::text", *f.TitleStartsWith, true)
	}

	if err := query.Find(&entities).Error; err != nil {
		return service.TroubleshootingStepListResponse{}, errors.WithStack(err)
	}

	return service.TroubleshootingStepListResponse{Entities: entities.toSvc()}, nil
}

func (r *TroubleshootingSteps) Get(ctx context.Context, f service.TroubleshootingStepGetFilter) (service.TroubleshootingStepEntity, error) {
	var entity TroubleshootingStepEntity

	if err := r.DB(ctx).Where(f.FilterMap()).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.TroubleshootingStepEntity{}, errorext.NewNotFound(err, errorext.ErrNotFound)
		}
	}

	return entity.toSvc(), nil
}

func (r *TroubleshootingSteps) UpdateHints(ctx context.Context, id uint, hints map[string]any) error {
	return nil
}
