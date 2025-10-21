package service

import "context"

type UserRepository interface {
	First(ctx context.Context, f FirstUserFilter) (UserEntity, error)
	Create(ctx context.Context, req CreateUserRequest) (UserEntity, error)
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
