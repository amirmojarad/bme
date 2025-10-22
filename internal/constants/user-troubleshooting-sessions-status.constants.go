package constants

type UserTroubleshootingSessionsStatus string

const (
	UserTroubleshootingSessionActive   UserTroubleshootingSessionsStatus = "active"
	UserTroubleshootingSessionDone     UserTroubleshootingSessionsStatus = "done"
	UserTroubleshootingSessionDeclined UserTroubleshootingSessionsStatus = "declined"

	UserTroubleshootingSessionDefaultOnCreation = UserTroubleshootingSessionActive
)

func (enum UserTroubleshootingSessionsStatus) String() string {
	return string(enum)
}

func (enum UserTroubleshootingSessionsStatus) IsEmpty() bool {
	return enum.String() == ""
}

func (enum UserTroubleshootingSessionsStatus) OrDefault() *UserTroubleshootingSessionsStatus {
	if enum == "" {
		defaultValue := UserTroubleshootingSessionDefaultOnCreation

		return &defaultValue
	}

	return &enum
}
