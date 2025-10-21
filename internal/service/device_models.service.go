package service

import (
	"bme/internal/constants"
	"time"
)

type (
	DeviceEntities                                 []DeviceEntity
	DeviceErrorCreateRequests                      []DeviceErrorCreateRequest
	DeviceErrorEntities                            []DeviceErrorEntity
	TroubleshootingStepEntities                    []TroubleshootingStepEntity
	TroubleShootingStepCreateRequests              []TroubleShootingStepCreateEntity
	TroubleshootingStepsToStepsEntities            []TroubleshootingStepsToStepsEntity
	TroubleshootingStepsToStepsCreateEntities      []TroubleshootingStepsToStepsCreateEntity
	TroubleshootingStepNextStepEntities            []TroubleshootingStepNextStepEntity
	TroubleshootingStepsToStepsWithDetailsEntities []TroubleshootingStepsToStepsWithDetailsEntity
)

type DeviceEntity struct {
	ID                  uint
	DeviceErrorEntities DeviceErrorEntities
	Title               string
	Description         string
	Status              constants.DeviceStatus
	CreatedBy           uint
	UpdatedBy           uint
	DeletedBy           uint
}

type CreateDeviceRequest struct {
	Title       string
	Description string
	Status      constants.DeviceStatus
	RequestedBy uint
}
type GetDeviceFilter struct {
	ID uint
}

type ListDevicesFilter struct {
	IdStartsWith    *string
	TitleStartsWith *string
	Status          *constants.DeviceStatus
}

type ListDevicesResponse struct {
	Items DeviceEntities
}

type DeviceErrorEntity struct {
	ID           uint
	DeviceID     uint
	DeviceEntity DeviceEntity
	Title        string
	Description  string
	Status       string
	CreatedBy    uint
	UpdatedBy    uint
	DeletedBy    uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}

type DeviceErrorCreateRequest struct {
	DeviceID    uint
	Title       string
	Description string
	Status      constants.DeviceErrorStatus
}

type DeviceErrorBulkCreateRequest struct {
	Entities    []DeviceErrorCreateRequest
	RequestedBy uint
}

type DeviceErrorListFilter struct {
	DeviceID        *uint
	IdStartsWith    *string
	IDs             []uint
	TitleStartsWith *string
	WithDetails     bool
	Status          *constants.DeviceErrorStatus
}

type DeviceErrorListResponse struct {
	Entities DeviceErrorEntities
}

type TroubleshootingStepEntity struct {
	ID            uint
	DeviceID      uint
	DeviceErrorID uint
	Title         string
	Description   string
	Hints         map[string]any
	Status        constants.TroubleshootingStepStatus
	NextSteps     TroubleshootingStepEntities
	CreatedBy     uint
	UpdatedBy     uint
	DeletedBy     uint
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

type TroubleShootingStepCreateEntity struct {
	DeviceID      uint
	DeviceErrorID uint
	Title         string
	Description   string
	Hints         map[string]any
	Status        constants.TroubleshootingStepStatus
	CreatedBy     uint
	UpdatedBy     uint
}
type TroubleshootingStepListResponse struct {
	Entities TroubleshootingStepEntities
}
type TroubleshootingBulkCreateRequest struct {
	Entities      TroubleShootingStepCreateRequests
	RequestedBy   uint
	DeviceID      uint
	DeviceErrorID uint
}
type TroubleshootingStepListFilter struct {
	IdStartsWith    *string
	IDs             []uint
	TitleStartsWith *string
	DeviceErrorID   *uint
	DeviceID        *uint
	Status          *constants.TroubleshootingStepStatus
}

type TroubleshootingStepGetFilter struct {
	DeviceID      *uint
	DeviceErrorID *uint
	ID            *uint
	WithNextSteps bool
}

type TroubleshootingStepsToStepsWithDetailsEntity struct {
	ID                            uint
	FromStepID                    uint
	FromTroubleshootingStepEntity TroubleshootingStepEntity
	ToStepID                      uint
	ToTroubleshootingStepEntity   TroubleshootingStepEntity
	Priority                      constants.TroubleshootingStepsToStepsPriority
	PriorityTitle                 string
	CreatedBy                     uint
	UpdatedBy                     uint
	DeletedBy                     uint
	CreatedAt                     time.Time
	UpdatedAt                     time.Time
	DeletedAt                     *time.Time
}

type TroubleshootingStepsToStepsEntity struct {
	ID         uint
	FromStepID uint
	ToStepID   uint
	Priority   int
	CreatedBy  uint
	UpdatedBy  uint
	DeletedBy  uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

type TroubleshootingStepsToStepsCreateEntity struct {
	FromStepID uint
	ToStepID   uint
	Priority   constants.TroubleshootingStepsToStepsPriority
	CreatedBy  uint
	UpdatedBy  uint
}

type TroubleshootingStepsToStepsCreateReq struct {
	Entities    TroubleshootingStepsToStepsCreateEntities
	RequestedBy uint
}

type TroubleshootingStepNextStepEntity struct {
	ToStepID uint
	Priority constants.TroubleshootingStepsToStepsPriority
}

type CreateTroubleshootingNextStepsReq struct {
	DeviceID      uint
	DeviceErrorID uint
	ID            uint
	NextSteps     TroubleshootingStepNextStepEntities
	RequestedBy   uint
}

func (f GetDeviceFilter) FilterMap() map[string]any {
	return map[string]any{
		"id": f.ID,
	}
}

func (f ListDevicesFilter) FilterMap() map[string]any {
	filterMap := make(map[string]any)

	if f.Status != nil {
		filterMap["status"] = f.Status.String()
	}

	return filterMap
}

func (f DeviceErrorListFilter) FilterMap() map[string]any {
	filterMap := make(map[string]any)

	if f.IDs != nil {
		filterMap["id"] = f.IDs
	}

	if f.DeviceID != nil {
		filterMap["device_id"] = f.DeviceID
	}

	if f.Status != nil {
		filterMap["status"] = f.Status.String()
	}

	return filterMap
}

func (f TroubleshootingStepGetFilter) FilterMap() map[string]any {
	filterMap := make(map[string]any)

	if f.DeviceID != nil {
		filterMap["device_id"] = f.DeviceID
	}

	if f.DeviceErrorID != nil {
		filterMap["device_error_id"] = f.DeviceErrorID
	}

	if f.ID != nil {
		filterMap["id"] = f.ID
	}

	return filterMap
}

func (f TroubleshootingStepListFilter) FilterMap() map[string]any {
	filterMap := make(map[string]any)

	if f.DeviceID != nil {
		filterMap["device_id"] = f.DeviceID
	}

	if f.Status != nil {
		filterMap["status"] = f.Status.String()
	}

	if f.DeviceErrorID != nil {
		filterMap["device_error_id"] = f.DeviceErrorID
	}

	if f.IDs != nil && len(f.IDs) > 0 {
		filterMap["id"] = f.IDs
	}

	return filterMap
}

func (req CreateTroubleshootingNextStepsReq) toTroubleshootingStepToStepsBulkCreateRequest() TroubleshootingStepsToStepsCreateReq {
	entities := make(TroubleshootingStepsToStepsCreateEntities, 0)

	for _, entity := range req.NextSteps {
		entities = append(entities, TroubleshootingStepsToStepsCreateEntity{
			FromStepID: req.ID,
			ToStepID:   entity.ToStepID,
			Priority:   entity.Priority,
			CreatedBy:  req.RequestedBy,
			UpdatedBy:  req.RequestedBy,
		})
	}

	return TroubleshootingStepsToStepsCreateReq{
		Entities:    entities,
		RequestedBy: 0,
	}
}
