package errorext

type ErrorMessage struct {
	En string
	Fa string
}

// GetFaOrEn returns ErrorMessage.Fa if it's not empty, otherwise returns ErrorMessage.En.
func (e ErrorMessage) GetFaOrEn() string {
	if e.Fa == "" {
		return e.En
	}

	return e.Fa
}

var ErrorMessages = map[ErrorKey]ErrorMessage{
	ErrGeneralOccurrence: {
		En: "an error occurred, please try again later",
		Fa: "خطایی ناشناخته رخ داده است،‌ لطفا مجددا تلاش کنید",
	},
	ErrConnectingToPostgres: {
		En: "error while connecting to postgres",
	},
	ErrGooseSetDialect: {
		En: "error while setting goose dialect",
	},
	ErrGooseUp: {
		En: "error while goose up",
	},
	ErrGormOpen: {
		En: "error while opening gorm connection",
	},
	ErrValidation: {
		En: "error on validation",
	},
	ErrNotFound: {
		En: "not found",
	},
	EnvValueIsEmpty: {
		En: "environment variable is empty",
	},
	ErrFailedGenerateToken: {
		En: "failed to generate token",
	},
	ErrInvalidPassword: {
		En: "invalid password",
	},
	ErrUnAuthorized: {
		En: "unauthorized",
		Fa: "خطای دسترسی",
	},
	ErrInvalidPassword: {
		En: "invalid password",
		Fa: "رمز عبور فعلی اشتباه است",
	},
}
