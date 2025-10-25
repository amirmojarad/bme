package service

import (
	"bme/pkg/errorext"
	"bme/pkg/utils"
	"context"
	"github.com/pkg/errors"
)

type AuthUserRepository interface {
	Create(ctx context.Context, req CreateUserRequest) (UserEntity, error)
	First(ctx context.Context, req FirstUserFilter) (UserEntity, error)
}

type AuthUserTroubleshootingSessionsRepository interface {
	First(ctx context.Context, f UserTroubleshootingSessionGetFilter) (UserTroubleshootingSessionEntity, error)
}

type Auth struct {
	repo           AuthUserRepository
	userTsSessions AuthUserTroubleshootingSessionsRepository
}

func NewAuth(
	repo AuthUserRepository,
	userTsSessions AuthUserTroubleshootingSessionsRepository,
) Auth {
	return Auth{
		repo:           repo,
		userTsSessions: userTsSessions,
	}
}

func (s Auth) Register(ctx context.Context, req AuthRegisterRequest) (UserEntity, error) {
	if err := req.validate(); err != nil {
		return UserEntity{}, err
	}

	foundedUser, err := s.repo.First(ctx, FirstUserFilter{
		Username: &req.Username,
	})
	if err != nil {
		if !errorext.IsNotFound(err) {
			return UserEntity{}, err
		}
	}

	if !foundedUser.isEmpty() {
		return UserEntity{}, errorext.NewValidation(errors.New("user already exists"), errorext.ErrValidation)
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return UserEntity{}, err
	}

	return s.repo.Create(ctx, CreateUserRequest{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		PhoneNumber:    req.PhoneNumber,
	})
}

func (s Auth) Login(ctx context.Context, req AuthLoginRequest) (AuthLoginResponse, error) {
	if err := req.validate(); err != nil {
		return AuthLoginResponse{}, err
	}

	foundedUser, err := s.repo.First(ctx, FirstUserFilter{
		Username: &req.Username,
	})
	if err != nil {
		return AuthLoginResponse{}, err
	}

	if foundedUser.isEmpty() {
		return AuthLoginResponse{}, errorext.NewNotFound(errors.New("user not found"), errorext.ErrNotFound)
	}

	if err = utils.VerifyPassword(req.Password, foundedUser.HashedPassword); err != nil {
		return AuthLoginResponse{}, err
	}

	return AuthLoginResponse{
		UserEntity: foundedUser,
	}, nil
}
