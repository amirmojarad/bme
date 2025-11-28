package repository

import (
	"bme/database"
	"bme/internal/service"
	"bme/pkg/errorext"
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"
	"time"
)

type UserTroubleshootingSessions struct {
	*database.GormWrapper
}

func NewUserTroubleshootingSessions(db *database.GormWrapper) *UserTroubleshootingSessions {
	return &UserTroubleshootingSessions{
		db,
	}
}

func (r UserTroubleshootingSessions) Create(
	ctx context.Context,
	req service.UserTroubleshootingSessionCreateRequest,
) (uint, error) {
	entity := userTroubleshootingSessionEntityFromSvcReq(req)

	if err := r.DB(ctx).Create(&entity).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return 0, errorext.NewValidation(err, errorext.ErrUserHasActiveSessionAlready)
		}

		return 0, err
	}

	return entity.ID, nil
}

func (r UserTroubleshootingSessions) ListWithDetails(
	ctx context.Context,
	f service.UserTroubleshootingSessionListFilter,
) (service.UserTroubleshootingSessionsListWithDetailsResp, error) {
	var (
		entities       = make(UserTroubleshootingSessionEntities, 0)
		paginationMeta = &service.PaginationMeta{}
	)

	query := r.DB(ctx).
		Model(&UserTroubleshootingSessionEntity{}).
		Where(f.FilterMap()).Preload("Device").Preload("DeviceError")

	if f.TitleStartsWith != nil {
		prefix := *f.TitleStartsWith + "%"

		query = query.
			Joins("LEFT JOIN device_errors DeviceError ON DeviceError.id = user_troubleshooting_sessions.device_error_id").
			Joins("LEFT JOIN devices Device ON Device.id = user_troubleshooting_sessions.device_id").
			Where(
				"(LOWER(DeviceError.title) LIKE LOWER(?) OR LOWER(Device.title) LIKE LOWER(?))",
				prefix, prefix,
			)
	}

	if f.PaginationRequest != nil {
		paginationMeta = f.PaginationRequest.PaginationMeta()
		query.Scopes(Paginate(paginationMeta, WithFullQ(query)))
	}

	if err := query.Order("CASE WHEN user_troubleshooting_sessions.status = 'active' THEN 1    WHEN user_troubleshooting_sessions.status = 'done'   THEN 2    ELSE 3  END").Order("created_at desc").Debug().Find(&entities).Error; err != nil {
		return service.UserTroubleshootingSessionsListWithDetailsResp{}, errors.WithStack(err)
	}

	return service.UserTroubleshootingSessionsListWithDetailsResp{
		Items:          entities.toSvc(),
		PaginationMeta: *paginationMeta,
	}, nil
}

func (r UserTroubleshootingSessions) UpdateStatus(
	ctx context.Context,
	req service.UserTroubleshootingUpdateStatusRequest,
) error {
	return errors.WithStack(r.DB(ctx).Model(&UserTroubleshootingSessionEntity{}).Where(req.FilterMap()).Updates(req.UpdatesMap()).Error)
}

func (r UserTroubleshootingSessions) First(ctx context.Context, f service.UserTroubleshootingSessionGetFilter) (service.UserTroubleshootingSessionEntity, error) {
	var entity UserTroubleshootingSessionEntity

	err := r.DB(ctx).Model(&UserTroubleshootingSessionEntity{}).Where(f.FilterMap()).First(&entity).Error
	if err != nil {
		return service.UserTroubleshootingSessionEntity{}, errors.WithStack(err)
	}

	return entity.toSvc(), nil
}

func (r UserTroubleshootingSessions) FirstWithDetails(ctx context.Context, f service.UserTroubleshootingSessionGetFilter) (service.UserTroubleshootingSessionWithDetailsEntity, error) {
	var entity UserTroubleshootingSessionEntity

	if err := r.DB(ctx).
		Model(&UserTroubleshootingSessionEntity{}).
		Where(f.FilterMap()).
		Preload("Device").
		Preload("DeviceError").
		Preload("CurrentTroubleshootingStep").
		First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.UserTroubleshootingSessionWithDetailsEntity{}, errorext.NewNotFound(err, errorext.ErrNotFound)
		}

		return service.UserTroubleshootingSessionWithDetailsEntity{}, errors.WithStack(err)
	}

	return entity.toSvcWithDetails(), nil
}

func (r UserTroubleshootingSessions) StepsMap(ctx context.Context, f service.UserTroubleshootingSessionListFilter) (service.TroubleshootingNextStepsMap, error) {
	transitions := make(TroubleshootingStepTransitions, 0)

	err := r.DB(ctx).
		Select(`
        ss.from_step_id,
        ts.title AS title,
        ss.to_step_id,
        to_ts.title AS to_title`).
		Table("troubleshooting_steps ts").
		Joins("JOIN troubleshooting_steps_to_steps ss ON ts.id = ss.from_step_id").
		Joins("JOIN troubleshooting_steps to_ts ON to_ts.id = ss.to_step_id").
		Where("ts.device_id = ?", f.DeviceID).
		Where("ts.device_error_id = ?", f.DeviceErrorID).Debug().
		Scan(&transitions).Error

	if err != nil {
		return service.TroubleshootingNextStepsMap{}, errors.WithStack(err)
	}

	return transitions.toSvcTroubleshootingStepsMap(), nil
}

func (r UserTroubleshootingSessions) PrevSteps(ctx context.Context, f service.UserTroubleshootingSessionListFilter) (service.TroubleshootingNextStepsMap, error) {
	transitions := make(TroubleshootingStepPrevStepsEntities, 0)

	err := r.DB(ctx).
		Select(`
        tsts.to_step_id, to_step.title as to_step_title, tsts.from_step_id, from_step.title as from_step_title`).
		Table("troubleshooting_steps ts").
		Joins("join public.troubleshooting_steps_to_steps tsts on ts.id = tsts.to_step_id").
		Joins("join troubleshooting_steps from_step on from_step.id = tsts.from_step_id").
		Joins("join troubleshooting_steps to_step on to_step.id = tsts.to_step_id").
		Where("ts.device_id = ?", f.DeviceID).
		Where("ts.device_error_id = ?", f.DeviceErrorID).Debug().
		Scan(&transitions).Error

	if err != nil {
		return service.TroubleshootingNextStepsMap{}, errors.WithStack(err)
	}

	return transitions.toSvcTroubleshootingStepsMap(), nil
}

func (r UserTroubleshootingSessions) UpdateCurrentStepID(ctx context.Context, id uint, currentStepID uint) error {
	return errors.WithStack(r.DB(ctx).
		Model(UserTroubleshootingSessionEntity{}).
		Where("id = ?", id).
		Update("current_troubleshooting_step_id", currentStepID).Error)
}

func (r UserTroubleshootingSessions) Finish(ctx context.Context, id uint) error {
	return errors.WithStack(r.DB(ctx).Model(&UserTroubleshootingSessionEntity{}).Where("id = ?", id).Update("finished_at", time.Now()).Error)
}
