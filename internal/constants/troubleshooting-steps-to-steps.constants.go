package constants

type TroubleshootingStepsToStepsPriority int

const (
	TroubleshootingStepsToStepsPriorityHigh TroubleshootingStepsToStepsPriority = 1 + iota
	TroubleshootingStepsToStepsPriorityMedium
	TroubleshootingStepsToStepsPriorityLow
)

var TroubleshootingStepsToStepsPriorityNames = map[TroubleshootingStepsToStepsPriority]string{
	TroubleshootingStepsToStepsPriorityHigh:   "High",
	TroubleshootingStepsToStepsPriorityMedium: "Medium",
	TroubleshootingStepsToStepsPriorityLow:    "Low",
}

func (enum TroubleshootingStepsToStepsPriority) OrDefault() TroubleshootingStepsToStepsPriority {
	if enum <= 0 {
		return TroubleshootingStepsToStepsPriorityLow
	}

	return enum
}

func (enum TroubleshootingStepsToStepsPriority) String() string {
	return TroubleshootingStepsToStepsPriorityNames[enum]
}
