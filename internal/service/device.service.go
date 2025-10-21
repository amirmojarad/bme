package service

import (
	"bme/internal/constants"
	"context"
)

type DeviceRepository interface {
	Create(ctx context.Context, req CreateDeviceRequest) error
	Get(ctx context.Context, f GetDeviceFilter) (DeviceEntity, error)
	List(ctx context.Context, f ListDevicesFilter) (ListDevicesResponse, error)
}

type Device struct {
	repo DeviceRepository
}

func NewDevice(repo DeviceRepository) *Device {
	return &Device{
		repo: repo,
	}
}

func (s *Device) Create(ctx context.Context, req CreateDeviceRequest) error {
	if req.Status.IsEmpty() {
		req.Status = constants.DeviceStatusActive
	}
	
	return s.repo.Create(ctx, req)
}

func (s *Device) Get(ctx context.Context, f GetDeviceFilter) (DeviceEntity, error) {
	return s.repo.Get(ctx, f)
}

func (s *Device) List(ctx context.Context, f ListDevicesFilter) (ListDevicesResponse, error) {
	return s.repo.List(ctx, f)
}
