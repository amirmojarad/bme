package controller

import (
	"bme/internal/constants"
	"bme/internal/service"
	"bme/pkg/errorext"
	"github.com/pkg/errors"
	"time"
)

type (
	UserTroubleshootingSessionWithDetailsEntities []UserTroubleshootingSessionWithDetailsEntity
	TroubleshootingStepTitleAndIDEntities         []TroubleshootingStepTitleAndIDEntity
)

type UserTroubleshootingSessionWithDetailsEntity struct {
	ID                         uint                                        `json:"id,omitempty"`
	UserID                     uint                                        `json:"user_id,omitempty"`
	DeviceID                   uint                                        `json:"device_id,omitempty"`
	DeviceTitle                string                                      `json:"device_title,omitempty"`
	DeviceErrorID              uint                                        `json:"device_error_id,omitempty"`
	DeviceErrorTitle           string                                      `json:"device_error_title,omitempty"`
	Status                     constants.UserTroubleshootingSessionsStatus `json:"status,omitempty"`
	CreatedAt                  time.Time                                   `json:"created_at"`
	CurrentStepID              *uint                                       `json:"current_step_id,omitempty"`
	CurrentTroubleshootingStep TroubleShootingStepEntity                   `json:"current_troubleshooting_step,omitempty"`
	NextSteps                  TroubleshootingStepTitleAndIDEntities       `json:"next_steps"`
	PrevStep                   TroubleshootingStepTitleAndIDEntities       `json:"prev_steps"`
}

type TroubleshootingStepTitleAndIDEntity struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type UserTroubleshootingSessionCreateRequest struct {
	UserID        uint  `json:"-"`
	DeviceID      uint  `json:"device_id" binding:"required"`
	DeviceErrorID uint  `json:"device_error_id" binding:"required"`
	StartStepID   *uint `json:"start_step_id"`
}

type UserTroubleshootingSessionListFilter struct {
	UserID        *uint                                        `json:"-"`
	DeviceID      *uint                                        `form:"device_id"`
	DeviceErrorID *uint                                        `form:"device_error_id"`
	Status        *constants.UserTroubleshootingSessionsStatus `form:"status,omitempty"`
	Q             *string                                      `form:"q,omitempty"`
	PaginationRequest
}

type UserTroubleshootingUpdateStatusRequest struct {
	UserID    uint   `json:"-"`
	SessionID uint   `uri:"id" binding:"required"`
	Status    string `json:"status"`
}

type UserTroubleshootingSessionsListWithDetailsResp struct {
	Items UserTroubleshootingSessionWithDetailsEntities `json:"items"`
	PaginationMeta
}

type UserTroubleshootingNextStepRequest struct {
	UserID     uint `json:"-"`
	NextStepID uint `json:"next_step_id" binding:"required"`
}

type UserTroubleshootingPrevStepRequest struct {
	UserID     uint `json:"-"`
	PrevStepID uint `json:"prev_step_id" binding:"required"`
}

func (req UserTroubleshootingSessionCreateRequest) toSvc() service.UserTroubleshootingSessionCreateRequest {
	return service.UserTroubleshootingSessionCreateRequest{
		UserID:        req.UserID,
		DeviceID:      req.DeviceID,
		DeviceErrorID: req.DeviceErrorID,
		StartStepID:   req.StartStepID,
	}
}

func (f UserTroubleshootingSessionListFilter) toSvc() service.UserTroubleshootingSessionListFilter {
	_, titleStartsWith := qAs(f.Q)

	return service.UserTroubleshootingSessionListFilter{
		UserID:            f.UserID,
		DeviceID:          f.DeviceID,
		DeviceErrorID:     f.DeviceErrorID,
		Status:            f.Status,
		PaginationRequest: f.PaginationRequest.toSvc(),
		TitleStartsWith:   titleStartsWith,
	}
}

func toViewUserTroubleshootingSessionsListWithDetailsResp(resp service.UserTroubleshootingSessionsListWithDetailsResp) UserTroubleshootingSessionsListWithDetailsResp {
	return UserTroubleshootingSessionsListWithDetailsResp{
		Items:          toViewUserTroubleshootingSessionsWithDetailsEntities(resp.Items),
		PaginationMeta: toViewPaginationMetaFromSvc(resp.PaginationMeta),
	}
}

func toViewUserTroubleshootingSessionsWithDetailsEntities(svcEntities service.UserTroubleshootingSessionWithDetailsEntities) UserTroubleshootingSessionWithDetailsEntities {
	entities := make([]UserTroubleshootingSessionWithDetailsEntity, 0, len(svcEntities))

	for _, svcEntity := range svcEntities {
		entities = append(entities, toViewUserTroubleshootingSessionWithDetailsEntityFromSvc(svcEntity))
	}

	return entities
}

func toViewUserTroubleshootingSessionWithDetailsEntityFromSvc(svcEntity service.UserTroubleshootingSessionWithDetailsEntity) UserTroubleshootingSessionWithDetailsEntity {
	return UserTroubleshootingSessionWithDetailsEntity{
		ID:                         svcEntity.ID,
		UserID:                     svcEntity.UserID,
		DeviceID:                   svcEntity.DeviceID,
		DeviceTitle:                svcEntity.DeviceTitle,
		DeviceErrorID:              svcEntity.DeviceErrorID,
		DeviceErrorTitle:           svcEntity.DeviceErrorTitle,
		Status:                     svcEntity.Status,
		CurrentTroubleshootingStep: toViewTroubleshootingStepEntity(svcEntity.CurrentTroubleshootingStep),
		CreatedAt:                  svcEntity.CreatedAt,
		CurrentStepID:              svcEntity.CurrentStepID,
		NextSteps:                  toViewTroubleshootingStepTitleAndIDEntities(svcEntity.NextSteps),
		PrevStep:                   toViewTroubleshootingStepTitleAndIDEntities(svcEntity.PrevStep),
	}
}

func toViewTroubleshootingStepTitleAndIDEntity(entity service.TroubleshootingStepTitleAndIDEntity) TroubleshootingStepTitleAndIDEntity {
	return TroubleshootingStepTitleAndIDEntity{
		ID:    entity.ID,
		Title: entity.Title,
	}
}

func toViewTroubleshootingStepTitleAndIDEntities(entities service.TroubleshootingStepTitleAndIDEntities) TroubleshootingStepTitleAndIDEntities {
	result := make([]TroubleshootingStepTitleAndIDEntity, 0, len(entities))

	for _, entity := range entities {
		result = append(result, toViewTroubleshootingStepTitleAndIDEntity(entity))
	}

	return result
}

func (req UserTroubleshootingUpdateStatusRequest) toSvc() service.UserTroubleshootingUpdateStatusRequest {
	return service.UserTroubleshootingUpdateStatusRequest{
		ID:          req.SessionID,
		RequestedBy: req.UserID,
		NewStatus:   constants.UserTroubleshootingSessionsStatus(req.Status),
	}
}

func (req UserTroubleshootingUpdateStatusRequest) validate() error {
	if req.Status == "" {
		return errorext.NewValidation(errors.New("status is required"), errorext.ErrValidation)
	}

	return nil
}

func (req UserTroubleshootingNextStepRequest) toSvc() service.UserTroubleshootingNextStepRequest {
	return service.UserTroubleshootingNextStepRequest{
		UserID:     req.UserID,
		NextStepID: req.NextStepID,
	}
}

func (req UserTroubleshootingPrevStepRequest) toSvc() service.UserTroubleshootingPrevStepRequest {
	return service.UserTroubleshootingPrevStepRequest{
		UserID:     req.UserID,
		PrevStepID: req.PrevStepID,
	}
}
