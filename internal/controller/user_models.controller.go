package controller

import "bme/internal/service"

type UserEntity struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}

type UserResetPasswordRequest struct {
	RequestedBy     uint   `json:"-"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type UserUpdateRequest struct {
	RequestedBy uint    `json:"-"`
	PhoneNumber *string `json:"phone_number"`
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	Username    *string `json:"username"`
}

func (req UserUpdateRequest) toSvc() service.UserUpdateRequest {
	return service.UserUpdateRequest{
		RequestedBy: req.RequestedBy,
		PhoneNumber: req.PhoneNumber,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Username:    req.Username,
	}
}

func (req UserResetPasswordRequest) toSvc() service.UserResetPasswordRequest {
	return service.UserResetPasswordRequest{
		RequestedBy:     req.RequestedBy,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}
}

func toViewUserResponse(user service.UserEntity) UserEntity {
	return UserEntity{
		ID:          user.ID,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
	}
}
