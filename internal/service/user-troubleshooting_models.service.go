package service

import (
	"bme/internal/constants"
	"time"
)

type (
	UserTroubleshootingSessionEntities            []UserTroubleshootingSessionEntity
	UserTroubleshootingSessionWithDetailsEntities []UserTroubleshootingSessionWithDetailsEntity
	UserTroubleshootingJourneyEntities            []UserTroubleshootingJourneyEntity
	SessionStepEntities                           []SessionStepEntity
	TroubleshootingStepsToStepEntities            []TroubleshootingStepsToStepEntity
)
type UserTroubleshootingSessionEntity struct {
	ID            uint
	UserID        uint
	DeviceID      uint
	DeviceErrorID uint
	Status        constants.UserTroubleshootingSessionsStatus
	CreatedAt     time.Time
	DeletedAt     *time.Time
	FinishedAt    *time.Time
}

type UserTroubleshootingSessionCreateRequest struct {
	UserID        uint
	DeviceID      uint
	DeviceErrorID uint
	StartStepID   *uint
}

type UserTroubleshootingSessionWithDetailsEntity struct {
	ID                         uint
	UserID                     uint
	DeviceID                   uint
	DeviceTitle                string
	DeviceErrorID              uint
	DeviceErrorTitle           string
	Status                     constants.UserTroubleshootingSessionsStatus
	CurrentStepID              *uint
	CurrentTroubleshootingStep TroubleshootingStepEntity
	NextSteps                  TroubleshootingStepTitleAndIDEntities
	PrevStep                   TroubleshootingStepTitleAndIDEntities
	CreatedAt                  time.Time
	DeletedAt                  *time.Time
	FinishedAt                 *time.Time
}
type UserTroubleshootingJourneyEntity struct {
	ID                           uint
	UserTroubleshootingSessionID uint
	FromTroubleshootingStepID    uint
	FromTroubleshootingStepTitle string
	ToTroubleshootingStepID      uint
	ToTroubleshootingStepTitle   string
	Description                  string
	CreatedAt                    time.Time
	FinishedAt                   *time.Time
}

type SessionByIdFilter struct {
	UserID    uint
	SessionID uint
}

type SessionStepEntity struct {
	FromStepID    uint
	FromStepTitle string
	ToStepID      uint
	ToStepTitle   string
	CreatedAt     time.Time
	FinishedAt    *time.Time
}

type SessionByIdResponse struct {
	UserTroubleshootingSessionWithDetailsEntity
	Steps SessionStepEntities
}

type UserTroubleshootingSessionListFilter struct {
	UserID            *uint
	DeviceID          *uint
	DeviceErrorID     *uint
	TitleStartsWith   *string
	Status            *constants.UserTroubleshootingSessionsStatus
	PaginationRequest *PaginationRequest
}

type UserTroubleshootingUpdateStatusRequest struct {
	ID          uint
	RequestedBy uint
	NewStatus   constants.UserTroubleshootingSessionsStatus
}

type UserTroubleshootingSessionGetFilter struct {
	ID     *uint
	UserID *uint
	Status *constants.UserTroubleshootingSessionsStatus
}

type UserTroubleshootingSessionsListWithDetailsResp struct {
	Items UserTroubleshootingSessionWithDetailsEntities
	PaginationMeta
}

type UserTroubleshootingCurrentActiveSessionReq struct {
	UserID uint
}

type UserTroubleshootingNextStepRequest struct {
	UserID     uint
	NextStepID uint
}

type UserTroubleshootingPrevStepRequest struct {
	UserID     uint
	PrevStepID uint
}

type UserTroubleshootingJourneyCreateRequest struct {
	SessionID                 uint
	FromTroubleshootingStepID uint
	ToTroubleshootingStepID   uint
}

type TroubleshootingStepsListStepsFilter struct {
	DeviceID      *uint
	DeviceErrorID *uint
}

type TroubleshootingStepsToStepEntity struct {
	DeviceID         uint
	DeviceTitle      string
	DeviceErrorID    uint
	DeviceErrorTitle string
	FromStepID       uint
	FromStepTitle    string
	ToStepID         uint
	ToStepTitle      string
}

func (f UserTroubleshootingSessionListFilter) FilterMap() map[string]any {
	filter := make(map[string]any)

	if f.UserID != nil {
		filter["user_id"] = *f.UserID
	}

	if f.DeviceID != nil {
		filter["device_id"] = *f.DeviceID
	}

	if f.DeviceErrorID != nil {
		filter["device_error_id"] = *f.DeviceErrorID
	}

	if f.Status != nil {
		filter["status"] = *f.Status
	}

	return filter
}

func (req UserTroubleshootingUpdateStatusRequest) FilterMap() map[string]any {
	return map[string]any{
		"id":      req.ID,
		"user_id": req.RequestedBy,
	}
}

func (req UserTroubleshootingUpdateStatusRequest) UpdatesMap() map[string]any {
	return map[string]any{
		"status": req.NewStatus.String(),
	}
}

func (f UserTroubleshootingSessionGetFilter) FilterMap() map[string]any {
	filter := make(map[string]any)

	if f.UserID != nil {
		filter["user_id"] = *f.UserID
	}

	if f.ID != nil {
		filter["id"] = *f.ID
	}

	if f.Status != nil {
		filter["status"] = *f.Status
	}

	return filter
}

func (entities UserTroubleshootingJourneyEntities) toSessionStepEntities() SessionStepEntities {
	result := make(SessionStepEntities, 0, len(entities))

	for _, entity := range entities {
		result = append(result, SessionStepEntity{
			FromStepID:    entity.FromTroubleshootingStepID,
			FromStepTitle: entity.FromTroubleshootingStepTitle,
			ToStepID:      entity.ToTroubleshootingStepID,
			ToStepTitle:   entity.ToTroubleshootingStepTitle,
			CreatedAt:     entity.CreatedAt,
			FinishedAt:    entity.FinishedAt,
		})
	}

	return result
}
