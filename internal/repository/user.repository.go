package repository

import (
	"bme/database"
	"bme/internal/service"
	"bme/pkg/errorext"
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type User struct {
	*database.GormWrapper
}

func NewUserRepository(wrapper *database.GormWrapper) User {
	return User{
		wrapper,
	}
}

func (r User) Create(ctx context.Context, req service.CreateUserRequest) (service.UserEntity, error) {
	user := userFromSvcReq(req)

	if err := r.DB(ctx).Create(&user).Error; err != nil {
		return service.UserEntity{}, errors.WithStack(err)
	}

	return user.toSvc(), nil
}

func (r User) First(ctx context.Context, f service.FirstUserFilter) (service.UserEntity, error) {
	var user UserEntity

	if err := r.DB(ctx).Where(f.FilterMap()).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.UserEntity{}, errorext.NewNotFound(err, errorext.ErrNotFound)
		}

		return service.UserEntity{}, errors.WithStack(err)
	}

	return user.toSvc(), nil
}

func (r User) Update(ctx context.Context, req service.UserUpdateRequest) error {
	return errors.WithStack(r.DB(ctx).Model(&UserEntity{}).Where(req.FilterMap()).Updates(req.UpdatesMap()).Debug().Error)
}
