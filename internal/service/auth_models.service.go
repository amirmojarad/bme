package service

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
