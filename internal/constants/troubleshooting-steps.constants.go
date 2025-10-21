package constants

type TroubleshootingStepStatus string

const (
	TroubleshootingStepStatusActive   TroubleshootingStepStatus = "active"
	TroubleshootingStepStatusInactive TroubleshootingStepStatus = "inactive"

	TroubleshootingStepStatusDefaultOnCreation = TroubleshootingStepStatusActive
)

func (enum TroubleshootingStepStatus) String() string {
	return string(enum)
}

func (enum TroubleshootingStepStatus) IsEmpty() bool {
	return enum.String() == ""
}

func (enum TroubleshootingStepStatus) OrDefault() *TroubleshootingStepStatus {
	if enum == "" {
		defaultValue := TroubleshootingStepStatusDefaultOnCreation

		return &defaultValue
	}

	return &enum
}
