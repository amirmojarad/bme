package constants

type DeviceErrorStatus string

const (
	DeviceErrorStatusActive   DeviceErrorStatus = "active"
	DeviceErrorStatusInactive DeviceErrorStatus = "inactive"

	DeviceErrorStatusDefaultOnCreation = DeviceErrorStatusActive
)

func (enum DeviceErrorStatus) String() string {
	return string(enum)
}

func (enum DeviceErrorStatus) IsEmpty() bool {
	return enum.String() == ""
}

func (enum DeviceErrorStatus) OrDefault() *DeviceErrorStatus {
	if enum == "" {
		defaultValue := DeviceErrorStatusDefaultOnCreation

		return &defaultValue
	}

	return &enum
}
