package service

import (
	"bme/internal/constants"
	"bme/pkg/errorext"
	"context"
	"github.com/pkg/errors"
)

type UserTroubleshootingSessionsRepository interface {
	Create(ctx context.Context, req UserTroubleshootingSessionCreateRequest) error
	ListWithDetails(ctx context.Context, f UserTroubleshootingSessionListFilter) (UserTroubleshootingSessionsListWithDetailsResp, error)
	UpdateStatus(ctx context.Context, req UserTroubleshootingUpdateStatusRequest) error
	First(ctx context.Context, f UserTroubleshootingSessionGetFilter) (UserTroubleshootingSessionEntity, error)
	FirstWithDetails(ctx context.Context, f UserTroubleshootingSessionGetFilter) (UserTroubleshootingSessionWithDetailsEntity, error)
}

type UserTroubleshootingJourneyRepository interface {
}

type UserTroubleshooting struct {
	userTroubleshootingSessionsRepo      UserTroubleshootingSessionsRepository
	userTroubleshootingJourneyRepository UserTroubleshootingJourneyRepository
}

func NewUserTroubleshooting(
	userTroubleshootingSessionsRepo UserTroubleshootingSessionsRepository,
	userTroubleshootingJourneyRepository UserTroubleshootingJourneyRepository,
) *UserTroubleshooting {
	return &UserTroubleshooting{
		userTroubleshootingJourneyRepository: userTroubleshootingJourneyRepository,
		userTroubleshootingSessionsRepo:      userTroubleshootingSessionsRepo,
	}
}

func (s *UserTroubleshooting) CreateSession(ctx context.Context, req UserTroubleshootingSessionCreateRequest) error {
	return s.userTroubleshootingSessionsRepo.Create(ctx, req)
}

func (s *UserTroubleshooting) ListSessions(
	ctx context.Context,
	f UserTroubleshootingSessionListFilter,
) (UserTroubleshootingSessionsListWithDetailsResp, error) {
	return s.userTroubleshootingSessionsRepo.ListWithDetails(ctx, f)
}

func (s *UserTroubleshooting) UpdateStatus(ctx context.Context, req UserTroubleshootingUpdateStatusRequest) error {
	session, err := s.userTroubleshootingSessionsRepo.First(ctx, UserTroubleshootingSessionGetFilter{
		ID:     &req.ID,
		UserID: &req.RequestedBy,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	if session.Status == constants.UserTroubleshootingSessionDone {
		return errorext.NewValidation(errors.New("session is already done"), errorext.ErrSessionIsAlreadyDone)
	}

	if session.Status == req.NewStatus {
		return nil
	}

	return s.userTroubleshootingSessionsRepo.UpdateStatus(ctx, req)
}

func (s *UserTroubleshooting) CurrentActiveSession(ctx context.Context, req UserTroubleshootingCurrentActiveSessionReq) (UserTroubleshootingSessionWithDetailsEntity, error) {
	activeStatus := constants.UserTroubleshootingSessionActive

	return s.userTroubleshootingSessionsRepo.FirstWithDetails(ctx, UserTroubleshootingSessionGetFilter{
		UserID: &req.UserID,
		Status: &activeStatus,
	})
}
