package constants

type DeviceStatus string

const (
	DeviceStatusActive   DeviceStatus = "active"
	DeviceStatusInactive DeviceStatus = "inactive"
)

func (enum DeviceStatus) String() string {
	return string(enum)
}

func (enum DeviceStatus) IsEmpty() bool {
	return enum.String() == ""
}
