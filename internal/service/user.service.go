package service

import (
	"bme/pkg/errorext"
	"bme/pkg/utils"
	"context"
)

type UserRepository interface {
	First(ctx context.Context, f FirstUserFilter) (UserEntity, error)
	Create(ctx context.Context, req CreateUserRequest) (UserEntity, error)
	Update(ctx context.Context, req UserUpdateRequest) error
}

type User struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) User {
	return User{repo: repo}
}

func (s User) First(ctx context.Context, f FirstUserFilter) (UserEntity, error) {
	return s.repo.First(ctx, f)
}

func (s User) ResetPassword(ctx context.Context, req UserResetPasswordRequest) error {
	user, err := s.First(ctx, FirstUserFilter{
		ID: &req.RequestedBy,
	})
	if err != nil {
		return err
	}

	if err = utils.VerifyPassword(req.CurrentPassword, user.HashedPassword); err != nil {
		return errorext.NewValidation(err, errorext.ErrInvalidPassword)
	}

	hashedNewPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errorext.NewValidation(err, errorext.ErrGeneralOccurrence)
	}

	return s.repo.Update(ctx, UserUpdateRequest{
		RequestedBy:    req.RequestedBy,
		HashedPassword: &hashedNewPassword,
	})
}

func (s User) Update(ctx context.Context, req UserUpdateRequest) error {
	return s.repo.Update(ctx, req)
}
