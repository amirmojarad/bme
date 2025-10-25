package repository

import (
	"bme/database"
	"bme/internal/service"
	"context"
	"github.com/pkg/errors"
)

type TroubleshootingStepsToSteps struct {
	*database.GormWrapper
}

func NewTroubleshootingStepsToSteps(db *database.GormWrapper) *TroubleshootingStepsToSteps {
	return &TroubleshootingStepsToSteps{db}
}

func (r *TroubleshootingStepsToSteps) BulkCreate(ctx context.Context, req service.TroubleshootingStepsToStepsCreateReq) error {
	entities := troubleshootingStepsToStepsEntitiesFromSvcReq(req)

	return errors.WithStack(r.DB(ctx).Create(&entities).Error)
}
