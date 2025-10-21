package repository

import (
	"bme/database"
	"bme/internal/service"
	"bme/pkg/errorext"
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Device struct {
	*database.GormWrapper
}

func NewDevice(db *database.GormWrapper) *Device {
	return &Device{
		db,
	}
}

func (r *Device) Create(ctx context.Context, req service.CreateDeviceRequest) error {
	entity := deviceEntityFromSvcCreateDeviceRequest(req)

	return errors.WithStack(r.DB(ctx).Create(&entity).Error)
}

func (r *Device) Get(ctx context.Context, f service.GetDeviceFilter) (service.DeviceEntity, error) {
	var entity DeviceEntity

	query := r.DB(ctx).Model(&entity).Where(f.FilterMap())

	if err := query.Debug().First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.DeviceEntity{}, errorext.NewNotFound(err, errorext.ErrNotFound)
		}

		return service.DeviceEntity{}, errors.WithStack(err)
	}

	return entity.toSvc(), nil
}

func (r *Device) List(ctx context.Context, f service.ListDevicesFilter) (service.ListDevicesResponse, error) {
	entities := make(DeviceEntities, 0)

	query := r.DB(ctx).Model(&entities).Where(f.FilterMap())

	if f.IdStartsWith != nil {
		appendAppendStartsWith(query, "id::text", *f.IdStartsWith, false)
	}

	if f.TitleStartsWith != nil {
		appendAppendStartsWith(query, "title", *f.TitleStartsWith, true)
	}

	if err := query.Debug().Find(&entities).Error; err != nil {
		return service.ListDevicesResponse{}, errors.WithStack(err)
	}

	return service.ListDevicesResponse{Items: entities.toSvc()}, nil
}
