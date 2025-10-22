package errorext

type ErrorKey string

func (e ErrorKey) String() string {
	return string(e)
}

func (e ErrorKey) GetMessage() ErrorMessage {
	if msg, ok := ErrorMessages[e]; ok {
		return msg
	}

	return ErrorMessages[ErrGeneralOccurrence]
}

var (
	ErrGeneralOccurrence           ErrorKey = "general_error_occurrence"
	ErrConnectingToPostgres        ErrorKey = "connecting_to_postgres_error"
	ErrGooseSetDialect             ErrorKey = "goose_set_dialect_error"
	ErrGooseUp                     ErrorKey = "goose_up_error"
	ErrGormOpen                    ErrorKey = "gorm_open_error"
	ErrValidation                  ErrorKey = "validation_error"
	ErrNotFound                    ErrorKey = "not_found_error"
	ErrIncorrectPassword           ErrorKey = "incorrect_password_error"
	EnvValueIsEmpty                ErrorKey = "env_value_is_empty"
	ErrFailedGenerateToken         ErrorKey = "failed_generate_token_error"
	ErrNotImplemented              ErrorKey = "not_implemented_error"
	ErrInvalidPassword             ErrorKey = "invalid_password_error"
	ErrUnAuthorized                ErrorKey = "authorization_error"
	ErrUserHasActiveSessionAlready ErrorKey = "user_has_active_session_already"
	ErrSessionIsAlreadyDone        ErrorKey = "session_is_already_done"
)
