package repository

import (
	"bme/internal/service"
	"gorm.io/gorm"
)

type UserEntity struct {
	gorm.Model
	Username       string
	FirstName      string
	LastName       string
	HashedPassword string
	PhoneNumber    string
}

func (UserEntity) TableName() string {
	return "users"
}

func userFromSvcReq(req service.CreateUserRequest) UserEntity {
	return UserEntity{
		Username:       req.Username,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		PhoneNumber:    req.PhoneNumber,
		HashedPassword: req.HashedPassword,
	}
}

func (entity UserEntity) toSvc() service.UserEntity {
	return service.UserEntity{
		ID:             entity.ID,
		Username:       entity.Username,
		FirstName:      entity.FirstName,
		LastName:       entity.LastName,
		HashedPassword: entity.HashedPassword,
		PhoneNumber:    entity.PhoneNumber,
	}
}
