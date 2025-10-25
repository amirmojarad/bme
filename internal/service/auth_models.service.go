package service

import "bme/pkg/jwt"

type AuthorizeRequest struct {
	TelegramUserID *uint
	ChatID         *uint
}

type AuthorizeResponse struct {
	TelegramUserID       uint
	UserID               uint
	IsAdmin              bool
	CanSendMessageToChat bool
}

type AuthLoginResponse struct {
	UserEntity UserEntity
}

func (resp AuthLoginResponse) JwtClaims() jwt.UserClaims {
	return jwt.UserClaims{
		UserID:   resp.UserEntity.ID,
		Username: resp.UserEntity.Username,
	}
}
