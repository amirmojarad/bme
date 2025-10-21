package constants

type DeviceStatus string

const (
	DeviceStatusActive   DeviceStatus = "active"
	DeviceStatusInactive DeviceStatus = "inactive"

	DeviceStatusDefaultOnCreation = DeviceStatusActive
)

func (enum DeviceStatus) String() string {
	return string(enum)
}

func (enum DeviceStatus) IsEmpty() bool {
	return enum.String() == ""
}

func (enum DeviceStatus) OrDefault() *DeviceStatus {
	if enum == "" {
		defaultValue := DeviceStatusDefaultOnCreation

		return &defaultValue
	}

	return &enum
}
