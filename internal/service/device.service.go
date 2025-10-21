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

type Device struct {
	repo            DeviceRepository
	deviceErrorRepo DeviceDeviceErrorRepository
}

func NewDevice(repo DeviceRepository, deviceErrorRepo DeviceDeviceErrorRepository) *Device {
	return &Device{
		repo:            repo,
		deviceErrorRepo: deviceErrorRepo,
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
