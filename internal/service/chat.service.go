package service

import "github.com/sirupsen/logrus"

type ChatUserTroubleshootingService interface {
}

type Chat struct {
	logger *logrus.Entry
}

func NewChat(logger *logrus.Entry) *Chat {
	return &Chat{
		logger: logger,
	}
}
