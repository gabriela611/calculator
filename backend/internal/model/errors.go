package model

import "fmt"

// Standard API error codes
const (
	ErrCodeInvalidOperation    = "INVALID_OPERATION"
	ErrCodeMissingOperation    = "MISSING_OPERATION"
	ErrCodeMissingOperands     = "MISSING_OPERANDS"
	ErrCodeInvalidOperand      = "INVALID_OPERAND"
	ErrCodeInvalidOperandCount = "INVALID_OPERAND_COUNT"
	ErrCodeInvalidJSON         = "INVALID_JSON"
	ErrCodeDivisionByZero      = "DIVISION_BY_ZERO"
	ErrCodeNegativeSquareRoot  = "NEGATIVE_SQUARE_ROOT"
	ErrCodeInternalError       = "INTERNAL_ERROR"
)

// AppError defines structured domain errors carrying HTTP status, error code, and message.
type AppError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewAppError creates a new AppError instance.
func NewAppError(statusCode int, code, message string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}
