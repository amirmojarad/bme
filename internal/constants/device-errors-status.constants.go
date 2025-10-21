package constants

type DeviceErrorStatus string

const (
	DeviceErrorStatusActive   DeviceStatus = "active"
	DeviceErrorStatusInactive DeviceStatus = "inactive"

	DeviceErrorStatusDefaultOnCreation = DeviceErrorStatusActive
)

func (enum DeviceErrorStatus) String() string {
	return string(enum)
}

func (enum DeviceErrorStatus) IsEmpty() bool {
	return enum.String() == ""
}
