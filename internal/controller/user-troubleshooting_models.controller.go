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
)

type UserTroubleshootingSessionWithDetailsEntity struct {
	ID               uint                                        `json:"id,omitempty"`
	UserID           uint                                        `json:"user_id,omitempty"`
	DeviceID         uint                                        `json:"device_id,omitempty"`
	DeviceTitle      string                                      `json:"device_title,omitempty"`
	DeviceErrorID    uint                                        `json:"device_error_id,omitempty"`
	DeviceErrorTitle string                                      `json:"device_error_title,omitempty"`
	Status           constants.UserTroubleshootingSessionsStatus `json:"status,omitempty"`
	CreatedAt        time.Time                                   `json:"created_at"`
	DeletedAt        *time.Time                                  `json:"deleted_at,omitempty"`
}

type UserTroubleshootingSessionCreateRequest struct {
	UserID        uint `json:"-"`
	DeviceID      uint `json:"device_id" binding:"required"`
	DeviceErrorID uint `json:"device_error_id" binding:"required"`
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

func (req UserTroubleshootingSessionCreateRequest) toSvc() service.UserTroubleshootingSessionCreateRequest {
	return service.UserTroubleshootingSessionCreateRequest{
		UserID:        req.UserID,
		DeviceID:      req.DeviceID,
		DeviceErrorID: req.DeviceErrorID,
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
		ID:               svcEntity.ID,
		UserID:           svcEntity.UserID,
		DeviceID:         svcEntity.DeviceID,
		DeviceTitle:      svcEntity.DeviceTitle,
		DeviceErrorID:    svcEntity.DeviceErrorID,
		DeviceErrorTitle: svcEntity.DeviceErrorTitle,
		Status:           svcEntity.Status,
		CreatedAt:        svcEntity.CreatedAt,
		DeletedAt:        svcEntity.DeletedAt,
	}
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
