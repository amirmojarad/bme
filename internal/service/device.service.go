package service

import (
	"context"
)

type DeviceRepository interface {
	Create(ctx context.Context, req CreateDeviceRequest) error
	Get(ctx context.Context, f GetDeviceFilter) (DeviceEntity, error)
	List(ctx context.Context, f ListDevicesFilter) (ListDevicesResponse, error)
}

type DeviceDeviceErrorRepository interface {
	BulkCreate(ctx context.Context, req DeviceErrorBulkCreateRequest) error
	List(ctx context.Context, f DeviceErrorListFilter) (DeviceErrorListResponse, error)
}

type DeviceTroubleshootingStepsRepository interface {
	BulkCreate(ctx context.Context, req TroubleshootingBulkCreateRequest) error
	List(ctx context.Context, req TroubleshootingStepListFilter) (TroubleshootingStepListResponse, error)
	Get(ctx context.Context, f TroubleshootingStepGetFilter) (TroubleshootingStepEntity, error)
	UpdateHints(ctx context.Context, id uint, hints map[string]any) error
}

type DeviceTroubleshootingStepsToStepsRepository interface {
	BulkCreate(ctx context.Context, req TroubleshootingStepsToStepsCreateReq) error
	List(ctx context.Context,
		filter TroubleshootingStepsListStepsFilter) (TroubleshootingStepsToStepEntities, error)
}
type Device struct {
	repo                            DeviceRepository
	deviceErrorRepo                 DeviceDeviceErrorRepository
	troubleshootingStepRepo         DeviceTroubleshootingStepsRepository
	troubleshootingStepsToStepsRepo DeviceTroubleshootingStepsToStepsRepository
}

func NewDevice(
	repo DeviceRepository,
	deviceErrorRepo DeviceDeviceErrorRepository,
	troubleshootingStepRepo DeviceTroubleshootingStepsRepository,
	troubleshootingStepsToStepsRepo DeviceTroubleshootingStepsToStepsRepository,
) *Device {
	return &Device{
		repo:                            repo,
		deviceErrorRepo:                 deviceErrorRepo,
		troubleshootingStepRepo:         troubleshootingStepRepo,
		troubleshootingStepsToStepsRepo: troubleshootingStepsToStepsRepo,
	}
}

func (s *Device) Create(ctx context.Context, req CreateDeviceRequest) error {
	return s.repo.Create(ctx, req)
}

func (s *Device) Get(ctx context.Context, f GetDeviceFilter) (DeviceEntity, error) {
	return s.repo.Get(ctx, f)
}

func (s *Device) List(ctx context.Context, f ListDevicesFilter) (ListDevicesResponse, error) {
	return s.repo.List(ctx, f)
}

func (s *Device) BulkCreateDeviceErrors(ctx context.Context, req DeviceErrorBulkCreateRequest) error {
	return s.deviceErrorRepo.BulkCreate(ctx, req)
}

func (s *Device) ListDeviceErrors(ctx context.Context, f DeviceErrorListFilter) (DeviceErrorListResponse, error) {
	return s.deviceErrorRepo.List(ctx, f)
}

func (s *Device) BulkCreateTroubleshootingSteps(ctx context.Context, req TroubleshootingBulkCreateRequest) error {
	return s.troubleshootingStepRepo.BulkCreate(ctx, req)
}

func (s *Device) GetTroubleshootingStep(ctx context.Context, f TroubleshootingStepGetFilter) (TroubleshootingStepEntity, error) {
	return s.troubleshootingStepRepo.Get(ctx, f)
}

func (s *Device) ListTroubleshootingSteps(ctx context.Context, f TroubleshootingStepListFilter) (TroubleshootingStepListResponse, error) {
	return s.troubleshootingStepRepo.List(ctx, f)
}

func (s *Device) CreateTroubleshootingNextSteps(ctx context.Context, req CreateTroubleshootingNextStepsReq) error {
	return s.troubleshootingStepsToStepsRepo.BulkCreate(ctx, req.toTroubleshootingStepToStepsBulkCreateRequest())
}
