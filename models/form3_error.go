package models

type AppError struct {
	Error   error
	Message string
	Code    int
}

func NewAppError(err error, message string, code int) AppError {
	return AppError{
		Error:   err,
		Message: message,
		Code:    code,
	}
}
