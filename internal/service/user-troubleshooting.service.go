package service

import (
	"bme/internal/constants"
	"bme/pkg/errorext"
	"context"
	"github.com/pkg/errors"
)

type UserTroubleshootingSessionsRepository interface {
	Create(ctx context.Context, req UserTroubleshootingSessionCreateRequest) (uint, error)
	ListWithDetails(ctx context.Context, f UserTroubleshootingSessionListFilter) (UserTroubleshootingSessionsListWithDetailsResp, error)
	UpdateStatus(ctx context.Context, req UserTroubleshootingUpdateStatusRequest) error
	UpdateCurrentStepID(ctx context.Context, id uint, currentStepID uint) error
	First(ctx context.Context, f UserTroubleshootingSessionGetFilter) (UserTroubleshootingSessionEntity, error)
	FirstWithDetails(ctx context.Context, f UserTroubleshootingSessionGetFilter) (UserTroubleshootingSessionWithDetailsEntity, error)
	StepsMap(ctx context.Context, f UserTroubleshootingSessionListFilter) (TroubleshootingNextStepsMap, error)
	PrevSteps(ctx context.Context, f UserTroubleshootingSessionListFilter) (TroubleshootingNextStepsMap, error)
	Finish(ctx context.Context, id uint) error
}

type UserTroubleshootingJourneyRepository interface {
	Create(ctx context.Context, req UserTroubleshootingJourneyCreateRequest) error
	Latest(ctx context.Context, sessionID uint) (UserTroubleshootingJourneyEntity, error)
	List(ctx context.Context, sessionID uint) (UserTroubleshootingJourneyEntities, error)
	Finish(ctx context.Context, id uint) error
}

type UserTroubleshootingTroubleshootingStepsRepository interface {
	Get(ctx context.Context, f TroubleshootingStepGetFilter) (TroubleshootingStepEntity, error)
	List(ctx context.Context, f TroubleshootingStepListFilter) (TroubleshootingStepListResponse, error)
}

type UserTroubleshooting struct {
	userTroubleshootingSessionsRepo      UserTroubleshootingSessionsRepository
	userTroubleshootingJourneyRepository UserTroubleshootingJourneyRepository
	tsStepRepo                           UserTroubleshootingTroubleshootingStepsRepository
	txRepo                               TransactionalRepository
}

func NewUserTroubleshooting(
	userTroubleshootingSessionsRepo UserTroubleshootingSessionsRepository,
	userTroubleshootingJourneyRepository UserTroubleshootingJourneyRepository,
	tsStepRepo UserTroubleshootingTroubleshootingStepsRepository,
	txRepo TransactionalRepository,
) *UserTroubleshooting {
	return &UserTroubleshooting{
		userTroubleshootingJourneyRepository: userTroubleshootingJourneyRepository,
		userTroubleshootingSessionsRepo:      userTroubleshootingSessionsRepo,
		tsStepRepo:                           tsStepRepo,
		txRepo:                               txRepo,
	}
}

func (s *UserTroubleshooting) CreateSession(ctx context.Context, req UserTroubleshootingSessionCreateRequest) (UserTroubleshootingSessionWithDetailsEntity, error) {
	if req.StartStepID != nil {
		if _, err := s.tsStepRepo.Get(ctx, TroubleshootingStepGetFilter{
			DeviceID:      &req.DeviceID,
			DeviceErrorID: &req.DeviceErrorID,
		}); err != nil {
			return UserTroubleshootingSessionWithDetailsEntity{}, err
		}
	} else {
		troubleshootingStep, err := s.tsStepRepo.Get(ctx, TroubleshootingStepGetFilter{
			DeviceID:      &req.DeviceID,
			DeviceErrorID: &req.DeviceErrorID,
			Sort:          true,
		})
		if err != nil {
			if errorext.IsNotFound(err) {
				return UserTroubleshootingSessionWithDetailsEntity{}, errorext.New(errors.New("device error has not steps"), errorext.ErrNotFound)
			}

			return UserTroubleshootingSessionWithDetailsEntity{}, err
		}

		req.StartStepID = &troubleshootingStep.ID
	}

	if _, err := s.userTroubleshootingSessionsRepo.Create(ctx, req); err != nil {
		return UserTroubleshootingSessionWithDetailsEntity{}, err
	}

	return s.CurrentActiveSession(ctx, UserTroubleshootingCurrentActiveSessionReq{UserID: req.UserID})
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

	if session.FinishedAt != nil {
		return errorext.NewValidation(errors.New("session is already finished"), errorext.ErrSessionIsAlreadyFinished)
	}

	if session.Status == req.NewStatus {
		return nil
	}

	return s.txRepo.BeginTx(ctx, func(tx context.Context) error {
		if err = s.userTroubleshootingSessionsRepo.UpdateStatus(ctx, req); err != nil {
			return err
		}

		return s.userTroubleshootingSessionsRepo.Finish(ctx, session.ID)
	})
}

func (s *UserTroubleshooting) CurrentActiveSession(ctx context.Context, req UserTroubleshootingCurrentActiveSessionReq) (UserTroubleshootingSessionWithDetailsEntity, error) {
	var (
		activeStatus = constants.UserTroubleshootingSessionActive
		nextSteps    = make(TroubleshootingStepTitleAndIDEntities, 0)
		prevSteps    = make(TroubleshootingStepTitleAndIDEntities, 0)
	)

	currentActiveSession, err := s.userTroubleshootingSessionsRepo.FirstWithDetails(ctx, UserTroubleshootingSessionGetFilter{
		UserID: &req.UserID,
		Status: &activeStatus,
	})
	if err != nil {
		return UserTroubleshootingSessionWithDetailsEntity{}, err
	}
	if currentActiveSession.CurrentStepID == nil {
		return currentActiveSession, errorext.NewValidation(errors.New("session does not have current step"), errorext.ErrSessionIsAlreadyDone)
	}

	nextStepsMap, err := s.userTroubleshootingSessionsRepo.StepsMap(ctx, UserTroubleshootingSessionListFilter{
		UserID:        &req.UserID,
		DeviceID:      &currentActiveSession.DeviceID,
		DeviceErrorID: &currentActiveSession.DeviceErrorID,
	})
	if err != nil {
		return UserTroubleshootingSessionWithDetailsEntity{}, err
	}

	prevStepsMap, err := s.userTroubleshootingSessionsRepo.PrevSteps(ctx, UserTroubleshootingSessionListFilter{
		UserID:        &req.UserID,
		DeviceID:      &currentActiveSession.DeviceID,
		DeviceErrorID: &currentActiveSession.DeviceErrorID,
	})
	if err != nil {
		return UserTroubleshootingSessionWithDetailsEntity{}, err
	}

	nextSteps = nextStepsMap.toTroubleshootingStepTitleAndIDEntities(*currentActiveSession.CurrentStepID)
	prevSteps = prevStepsMap.toTroubleshootingStepTitleAndIDEntities(*currentActiveSession.CurrentStepID)

	go func() {
		if currentNextSteps, ok := nextStepsMap.Map[*currentActiveSession.CurrentStepID]; ok {
			if len(currentNextSteps) == 0 {
				if err = s.userTroubleshootingSessionsRepo.UpdateStatus(ctx, UserTroubleshootingUpdateStatusRequest{
					ID:          currentActiveSession.ID,
					RequestedBy: req.UserID,
					NewStatus:   constants.UserTroubleshootingSessionDone,
				}); err != nil {
				}
			}
		}
	}()

	return UserTroubleshootingSessionWithDetailsEntity{
		ID:                         currentActiveSession.ID,
		UserID:                     currentActiveSession.UserID,
		DeviceID:                   currentActiveSession.DeviceID,
		DeviceTitle:                currentActiveSession.DeviceTitle,
		DeviceErrorID:              currentActiveSession.DeviceErrorID,
		DeviceErrorTitle:           currentActiveSession.DeviceErrorTitle,
		Status:                     currentActiveSession.Status,
		CurrentStepID:              currentActiveSession.CurrentStepID,
		CurrentTroubleshootingStep: currentActiveSession.CurrentTroubleshootingStep,
		NextSteps:                  nextSteps,
		PrevStep:                   prevSteps,
	}, nil
}

func (s *UserTroubleshooting) NextStep(ctx context.Context, req UserTroubleshootingNextStepRequest) error {
	var (
		activeSession = constants.UserTroubleshootingSessionActive
	)

	session, err := s.userTroubleshootingSessionsRepo.FirstWithDetails(ctx, UserTroubleshootingSessionGetFilter{
		UserID: &req.UserID,
		Status: &activeSession,
	})
	if err != nil {
		return err
	}

	if session.CurrentStepID == nil {
		return errors.New("no current step found")
	}

	stepsMap, err := s.userTroubleshootingSessionsRepo.StepsMap(ctx, UserTroubleshootingSessionListFilter{
		DeviceID:      &session.DeviceID,
		DeviceErrorID: &session.DeviceErrorID,
	})
	if err != nil {
		return err
	}

	if _, ok := stepsMap.Map[*session.CurrentStepID][req.NextStepID]; !ok {
		return errorext.NewValidation(errors.New("current step is not related to next step"), errorext.ErrValidation)
	}

	return s.txRepo.BeginTx(ctx, func(tx context.Context) error {
		latestStep, err := s.userTroubleshootingJourneyRepository.Latest(tx, session.ID)
		if err != nil {
			return err
		}

		if err = s.userTroubleshootingJourneyRepository.Finish(tx, latestStep.ID); err != nil {
			return err
		}

		if err = s.userTroubleshootingJourneyRepository.Create(tx, UserTroubleshootingJourneyCreateRequest{
			SessionID:                 session.ID,
			FromTroubleshootingStepID: *session.CurrentStepID,
			ToTroubleshootingStepID:   req.NextStepID,
		}); err != nil {
			return err
		}

		return s.userTroubleshootingSessionsRepo.UpdateCurrentStepID(tx, session.ID, req.NextStepID)
	})
}

func (s *UserTroubleshooting) DeclineSession(ctx context.Context, userID uint) error {
	currentActiveSession, err := s.CurrentActiveSession(ctx, UserTroubleshootingCurrentActiveSessionReq{userID})
	if err != nil {
		return err
	}

	return s.UpdateStatus(ctx, UserTroubleshootingUpdateStatusRequest{
		ID:          currentActiveSession.ID,
		RequestedBy: userID,
		NewStatus:   constants.UserTroubleshootingSessionDeclined,
	})
}

func (s *UserTroubleshooting) DoneSession(ctx context.Context, userID uint) error {
	currentActiveSession, err := s.CurrentActiveSession(ctx, UserTroubleshootingCurrentActiveSessionReq{userID})
	if err != nil {
		return err
	}

	return s.UpdateStatus(ctx, UserTroubleshootingUpdateStatusRequest{
		ID:          currentActiveSession.ID,
		RequestedBy: userID,
		NewStatus:   constants.UserTroubleshootingSessionDone,
	})
}

func (s *UserTroubleshooting) PrevStep(ctx context.Context, req UserTroubleshootingPrevStepRequest) error {
	var (
		activeSession = constants.UserTroubleshootingSessionActive
	)

	session, err := s.userTroubleshootingSessionsRepo.FirstWithDetails(ctx, UserTroubleshootingSessionGetFilter{
		UserID: &req.UserID,
		Status: &activeSession,
	})
	if err != nil {
		return err
	}
	if session.CurrentStepID == nil {
		return errors.New("no current step found")
	}

	prevStepsMap, err := s.userTroubleshootingSessionsRepo.PrevSteps(ctx, UserTroubleshootingSessionListFilter{
		DeviceID:      &session.DeviceID,
		DeviceErrorID: &session.DeviceErrorID,
	})
	if err != nil {
		return err
	}

	if _, ok := prevStepsMap.Map[*session.CurrentStepID][req.PrevStepID]; !ok {
		return errorext.NewValidation(errors.New("current step is not related to prev step"), errorext.ErrValidation)
	}

	return s.txRepo.BeginTx(ctx, func(tx context.Context) error {
		if err = s.userTroubleshootingJourneyRepository.Create(tx, UserTroubleshootingJourneyCreateRequest{
			SessionID:                 session.ID,
			FromTroubleshootingStepID: *session.CurrentStepID,
			ToTroubleshootingStepID:   req.PrevStepID,
		}); err != nil {
			return err
		}

		return s.userTroubleshootingSessionsRepo.UpdateCurrentStepID(ctx, session.ID, req.PrevStepID)
	})
}
