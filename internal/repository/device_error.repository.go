package repository

import (
	"bme/database"
	"bme/internal/service"
	"context"
	"github.com/pkg/errors"
	"strings"
)

type DeviceError struct {
	*database.GormWrapper
}

func NewDeviceError(db *database.GormWrapper) *DeviceError {
	return &DeviceError{db}
}

func (r *DeviceError) BulkCreate(ctx context.Context, req service.DeviceErrorBulkCreateRequest) error {
	entities := deviceErrorEntitiesFromSvcBulkCreateReq(req)

	return errors.WithStack(r.DB(ctx).Create(&entities).Error)
}

func (r *DeviceError) List(ctx context.Context, f service.DeviceErrorListFilter) (service.DeviceErrorListResponse, error) {
	entities := make(DeviceErrorEntities, 0)

	query := r.DB(ctx).Model(&DeviceErrorEntity{}).Where(f.FilterMap())

	if f.TitleStartsWith != nil {
		appendAppendStartsWith(query, "title", strings.ToLower(*f.TitleStartsWith), true)
	}

	if f.IdStartsWith != nil {
		appendAppendStartsWith(query, "id::text", *f.IdStartsWith, false)
	}

	if f.WithDetails {
		query.Preload("DeviceEntity")
	}

	if err := query.Debug().Find(&entities).Error; err != nil {
		return service.DeviceErrorListResponse{}, errors.WithStack(err)
	}

	return service.DeviceErrorListResponse{Entities: entities.toSvc()}, nil
}
