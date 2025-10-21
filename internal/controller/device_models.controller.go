package controller

import (
	"bme/internal/constants"
	"bme/internal/service"
	"github.com/go-playground/validator/v10"
)

type (
	DeviceEntities            []DeviceEntity
	DeviceErrorCreateRequests []DeviceErrorCreateRequest
	DeviceErrorEntities       []DeviceErrorEntity

	TroubleShootingStepEntities            []TroubleShootingStepEntity
	TroubleShootingStepCreateRequests      []TroubleShootingStepCreateEntity
	TroubleshootingStepNextStepEntities    []TroubleshootingStepNextStepEntity
	TroubleShootingStepWithDetailsEntities []TroubleShootingStepWithDetailsEntity
)

type DeviceEntity struct {
	ID          uint                   `json:"id,omitempty"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Status      constants.DeviceStatus `json:"status,omitempty"`
}

type CreateDeviceRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	RequestedBy uint   `json:"-" validate:"required"`
}
type GetDeviceFilter struct {
	ID uint `uri:"id" binding:"required"`
}

type ListDevicesFilter struct {
	Q      *string                `form:"q"`
	Status constants.DeviceStatus `form:"status"`
}

type ListDevicesResponse struct {
	Items DeviceEntities `json:"items"`
}

type DeviceErrorEntity struct {
	ID          uint   `json:"id"`
	DeviceID    uint   `json:"device_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type DeviceErrorCreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type DeviceErrorBulkCreateRequest struct {
	DeviceID    uint                      `uri:"id" binding:"required"`
	Items       DeviceErrorCreateRequests `json:"items"`
	RequestedBy uint                      `json:"-"`
}

type DeviceErrorListFilter struct {
	DeviceID uint                        `uri:"id" binding:"required"`
	Q        *string                     `form:"q"`
	Status   constants.DeviceErrorStatus `form:"status"`
}

type DeviceErrorListResponse struct {
	Items DeviceErrorEntities `json:"items"`
}

type TroubleShootingStepEntity struct {
	ID            uint                                `json:"id"`
	DeviceID      uint                                `json:"device_id"`
	DeviceErrorID uint                                `json:"device_error_id"`
	Title         string                              `json:"title"`
	Description   string                              `json:"description"`
	Hints         map[string]any                      `json:"hints"`
	Status        constants.TroubleshootingStepStatus `json:"status"`
	NextSteps     TroubleShootingStepEntities         `json:"next_steps"`
}

type TroubleShootingStepWithDetailsEntity struct {
	ID            uint                                `json:"id"`
	DeviceID      uint                                `json:"device_id"`
	DeviceErrorID uint                                `json:"device_error_id"`
	Title         string                              `json:"title"`
	Description   string                              `json:"description"`
	Hints         map[string]any                      `json:"hints"`
	Status        constants.TroubleshootingStepStatus `json:"status"`
	NextSteps     TroubleShootingStepEntities         `json:"next_steps"`
}

type TroubleShootingStepCreateEntity struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Hints       map[string]any `json:"hints"`
}
type TroubleshootingStepListResponse struct {
	Items TroubleShootingStepEntities `json:"items"`
}

type TroubleshootingBulkCreateRequest struct {
	Items         TroubleShootingStepCreateRequests `json:"items"`
	DeviceID      uint                              `uri:"id" binding:"required"`
	DeviceErrorID uint                              `uri:"error_id" binding:"required"`
	RequestedBy   uint                              `json:"-"`
}

type TroubleshootingStepGetFilter struct {
	DeviceID      uint `uri:"id" binding:"required"`
	DeviceErrorID uint `uri:"error_id" binding:"required"`
	ID            uint `uri:"troubleshooting_step_id" binding:"required"`
	WithDetails   bool `form:"with_details"`
}

type TroubleshootingStepsListFilter struct {
	DeviceID      uint                                `uri:"id" binding:"required"`
	DeviceErrorID uint                                `uri:"error_id" binding:"required"`
	Q             *string                             `form:"q"`
	Status        constants.TroubleshootingStepStatus `form:"status"`
}

type TroubleshootingStepNextStepEntity struct {
	ToStepID uint                                          `json:"to_step_id"`
	Priority constants.TroubleshootingStepsToStepsPriority `json:"priority"`
}

type CreateTroubleshootingNextStepsReq struct {
	DeviceID      uint                                `uri:"id" binding:"required"`
	DeviceErrorID uint                                `uri:"error_id" binding:"required"`
	ID            uint                                `uri:"troubleshooting_step_id" binding:"required"`
	NextSteps     TroubleshootingStepNextStepEntities `json:"next_steps"`
	RequestedBy   uint                                `json:"-"`
}

func (req CreateDeviceRequest) validate() error {
	validate := validator.New()

	if err := validate.Struct(&req); err != nil {
		return err
	}

	return nil
}

func (req CreateDeviceRequest) toSvc() service.CreateDeviceRequest {
	return service.CreateDeviceRequest{
		Title:       req.Title,
		Description: req.Description,
		RequestedBy: req.RequestedBy,
		Status:      constants.DeviceStatusDefaultOnCreation,
	}
}

func toViewDeviceEntity(svcEntity service.DeviceEntity) DeviceEntity {
	return DeviceEntity{
		ID:          svcEntity.ID,
		Title:       svcEntity.Title,
		Description: svcEntity.Description,
		Status:      svcEntity.Status,
	}
}

func (req GetDeviceFilter) toSvc() service.GetDeviceFilter {
	return service.GetDeviceFilter{
		ID: req.ID,
	}
}

func toViewDeviceEntities(svcEntities service.DeviceEntities) DeviceEntities {
	result := make(DeviceEntities, 0, len(svcEntities))

	for _, svcEntity := range svcEntities {
		result = append(result, toViewDeviceEntity(svcEntity))
	}

	return result
}

func toViewListDevicesResp(resp service.ListDevicesResponse) ListDevicesResponse {
	return ListDevicesResponse{
		Items: toViewDeviceEntities(resp.Items),
	}
}

func (f ListDevicesFilter) toSvc() service.ListDevicesFilter {
	idStartsWith, titleStartsWith := qAs(f.Q)

	return service.ListDevicesFilter{
		IdStartsWith:    idStartsWith,
		TitleStartsWith: titleStartsWith,
		Status:          f.Status.OrDefault(),
	}
}

func (req DeviceErrorBulkCreateRequest) toSvc() service.DeviceErrorBulkCreateRequest {
	return service.DeviceErrorBulkCreateRequest{
		Entities:    req.Items.toSvc(req.DeviceID),
		RequestedBy: req.RequestedBy,
	}
}

func (requests DeviceErrorCreateRequests) toSvc(deviceID uint) service.DeviceErrorCreateRequests {
	result := make(service.DeviceErrorCreateRequests, 0, len(requests))

	for _, request := range requests {
		result = append(result, request.toSvc(deviceID))
	}

	return result
}

func (req DeviceErrorCreateRequest) toSvc(deviceID uint) service.DeviceErrorCreateRequest {
	return service.DeviceErrorCreateRequest{
		DeviceID:    deviceID,
		Title:       req.Title,
		Description: req.Description,
		Status:      constants.DeviceErrorStatusDefaultOnCreation,
	}
}

func (f DeviceErrorListFilter) toSvc() service.DeviceErrorListFilter {
	idStartsWith, titleStartsWith := qAs(f.Q)

	return service.DeviceErrorListFilter{
		DeviceID:        &f.DeviceID,
		IdStartsWith:    idStartsWith,
		TitleStartsWith: titleStartsWith,
		Status:          f.Status.OrDefault(),
	}
}

func toViewListDeviceErrorsResp(response service.DeviceErrorListResponse) DeviceErrorListResponse {
	return DeviceErrorListResponse{
		Items: toViewDeviceErrorEntities(response.Entities),
	}
}

func toViewDeviceErrorEntities(svcEntities service.DeviceErrorEntities) DeviceErrorEntities {
	result := make(DeviceErrorEntities, 0, len(svcEntities))

	for _, svcEntity := range svcEntities {
		result = append(result, toViewDeviceErrorEntity(svcEntity))
	}

	return result
}

func toViewDeviceErrorEntity(svcEntity service.DeviceErrorEntity) DeviceErrorEntity {
	return DeviceErrorEntity{
		ID:          svcEntity.ID,
		DeviceID:    svcEntity.DeviceID,
		Title:       svcEntity.Title,
		Description: svcEntity.Description,
		Status:      svcEntity.Status,
	}
}

func (req TroubleshootingBulkCreateRequest) toSvc() service.TroubleshootingBulkCreateRequest {
	return service.TroubleshootingBulkCreateRequest{
		Entities:      req.Items.toSvc(req.DeviceID, req.DeviceErrorID, req.RequestedBy),
		RequestedBy:   req.RequestedBy,
		DeviceID:      req.DeviceID,
		DeviceErrorID: req.DeviceErrorID,
	}
}

func (requests TroubleShootingStepCreateRequests) toSvc(deviceID, deviceErrorID, requestedBy uint) service.TroubleShootingStepCreateRequests {
	result := make(service.TroubleShootingStepCreateRequests, 0, len(requests))

	for _, request := range requests {
		result = append(result, request.toSvc(deviceID, deviceErrorID, requestedBy))
	}

	return result
}

func (req TroubleShootingStepCreateEntity) toSvc(deviceID, deviceErrorID, requestedBy uint) service.TroubleShootingStepCreateEntity {
	return service.TroubleShootingStepCreateEntity{
		DeviceID:      deviceID,
		DeviceErrorID: deviceErrorID,
		Title:         req.Title,
		Description:   req.Description,
		Hints:         req.Hints,
		Status:        constants.TroubleshootingStepStatusDefaultOnCreation,
		CreatedBy:     requestedBy,
		UpdatedBy:     requestedBy,
	}
}

func toViewTroubleshootingStepEntity(entity service.TroubleshootingStepEntity) TroubleShootingStepEntity {
	return TroubleShootingStepEntity{
		ID:            entity.ID,
		DeviceID:      entity.DeviceID,
		DeviceErrorID: entity.DeviceID,
		Title:         entity.Title,
		Description:   entity.Description,
		Hints:         entity.Hints,
		Status:        entity.Status,
		NextSteps:     toViewTroubleshootingStepEntities(entity.NextSteps),
	}
}

func (f TroubleshootingStepGetFilter) toSvc() service.TroubleshootingStepGetFilter {
	return service.TroubleshootingStepGetFilter{
		DeviceID:      &f.DeviceID,
		DeviceErrorID: &f.DeviceErrorID,
		ID:            &f.ID,
		WithNextSteps: f.WithDetails,
	}
}

func (f TroubleshootingStepsListFilter) toSvc() service.TroubleshootingStepListFilter {
	idStartsWith, titleStartsWith := qAs(f.Q)

	return service.TroubleshootingStepListFilter{
		IdStartsWith:    idStartsWith,
		TitleStartsWith: titleStartsWith,
		DeviceErrorID:   &f.DeviceErrorID,
		DeviceID:        &f.DeviceID,
		Status:          f.Status.OrDefault(),
	}
}

func toViewTroubleshootingStepEntities(entities service.TroubleshootingStepEntities) TroubleShootingStepEntities {
	result := make(TroubleShootingStepEntities, 0, len(entities))

	for _, entity := range entities {
		result = append(result, toViewTroubleshootingStepEntity(entity))
	}

	return result
}

func toViewTroubleshootingStepListResponse(response service.TroubleshootingStepListResponse) TroubleshootingStepListResponse {
	return TroubleshootingStepListResponse{
		Items: toViewTroubleshootingStepEntities(response.Entities),
	}
}

func (req CreateTroubleshootingNextStepsReq) toSvc() service.CreateTroubleshootingNextStepsReq {
	return service.CreateTroubleshootingNextStepsReq{
		DeviceID:      req.DeviceID,
		DeviceErrorID: req.DeviceErrorID,
		ID:            req.ID,
		NextSteps:     req.NextSteps.toSvc(),
		RequestedBy:   req.RequestedBy,
	}
}

func (entities TroubleshootingStepNextStepEntities) toSvc() service.TroubleshootingStepNextStepEntities {
	result := make(service.TroubleshootingStepNextStepEntities, 0, len(entities))

	for _, entity := range entities {
		result = append(result, entity.toSvc())
	}

	return result
}

func (entity TroubleshootingStepNextStepEntity) toSvc() service.TroubleshootingStepNextStepEntity {
	return service.TroubleshootingStepNextStepEntity{
		ToStepID: entity.ToStepID,
		Priority: entity.Priority.OrDefault(),
	}
}
