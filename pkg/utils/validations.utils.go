package utils

type ValidationError struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Meta    []map[string]any `json:"meta"`
}

type ServiceValidationErrorMeta struct {
	// Attribute equals to struct's member title
	Attribute string `json:"attribute"`
	Error     error  `json:"error"`
	Message   string `json:"message"`
}

func (meta ServiceValidationErrorMeta) Map() map[string]any {
	return map[string]any{
		"attribute": meta.Attribute,
		"message":   meta.Message,
		"error":     meta.Error.Error(),
	}
}
