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

func (r *TroubleshootingStepsToSteps) List(ctx context.Context,
	filter service.TroubleshootingStepsListStepsFilter) (service.TroubleshootingStepsToStepEntities, error) {
	entities := make(TroubleshootingStepsToStepsWithDetailsEntities, 0)

	query := r.DB(ctx).
		Table("devices d").
		Joins("join public.device_errors de on d.id = de.device_id").
		Joins("join public.troubleshooting_steps from_step on de.id = from_step.device_error_id").
		Joins("join public.troubleshooting_steps_to_steps step_to_steps on step_to_steps.from_step_id = from_step.id").
		Joins("join troubleshooting_steps to_step on step_to_steps.to_step_id = to_step.id").Order("from_step.device_error_id, from_step.sort")
	if filter.DeviceID != nil {
		query.Where("d.id = ?", *filter.DeviceID)
	}

	if filter.DeviceErrorID != nil {
		query.Where("de.id = ?", *filter.DeviceErrorID)
	}

	if err := query.Find(&entities).Error; err != nil {
		return service.TroubleshootingStepsToStepEntities{}, errors.WithStack(err)
	}

	return entities.toSvc(), nil
}
