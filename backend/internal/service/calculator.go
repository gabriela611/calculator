package service

import (
	"fmt"
	"math"
	"net/http"

	"backend/internal/model"
)

// Supported operations
const (
	OpAdd        = "add"
	OpSubtract   = "subtract"
	OpMultiply   = "multiply"
	OpDivide     = "divide"
	OpPower      = "power"
	OpSqrt       = "sqrt"
	OpPercentage = "percentage"
)

var requiredOperandCounts = map[string]int{
	OpAdd:        2,
	OpSubtract:   2,
	OpMultiply:   2,
	OpDivide:     2,
	OpPower:      2,
	OpSqrt:       1,
	OpPercentage: 2,
}

// CalculatorService defines business logic interface for calculations.
type CalculatorService interface {
	Calculate(req model.CalculationRequest) (*model.CalculationResponse, *model.AppError)
}

type calculatorService struct{}

// NewCalculatorService constructs a new CalculatorService.
func NewCalculatorService() CalculatorService {
	return &calculatorService{}
}

// Calculate performs mathematical operations after validating operation and operand count rules.
func (s *calculatorService) Calculate(req model.CalculationRequest) (*model.CalculationResponse, *model.AppError) {
	if req.Operation == "" {
		return nil, model.NewAppError(http.StatusBadRequest, model.ErrCodeMissingOperation, "The operation field is required")
	}

	if req.Operands == nil {
		return nil, model.NewAppError(http.StatusBadRequest, model.ErrCodeMissingOperands, "The operands field is required")
	}

	expectedCount, exists := requiredOperandCounts[req.Operation]
	if !exists {
		return nil, model.NewAppError(http.StatusBadRequest, model.ErrCodeInvalidOperation, "Unsupported operation")
	}

	if len(req.Operands) != expectedCount {
		msg := fmt.Sprintf("Operation '%s' requires exactly %d operand", req.Operation, expectedCount)
		if expectedCount > 1 {
			msg += "s"
		}
		return nil, model.NewAppError(http.StatusBadRequest, model.ErrCodeInvalidOperandCount, msg)
	}

	var result float64

	switch req.Operation {
	case OpAdd:
		result = req.Operands[0] + req.Operands[1]

	case OpSubtract:
		result = req.Operands[0] - req.Operands[1]

	case OpMultiply:
		result = req.Operands[0] * req.Operands[1]

	case OpDivide:
		if req.Operands[1] == 0 {
			return nil, model.NewAppError(http.StatusUnprocessableEntity, model.ErrCodeDivisionByZero, "Cannot divide by zero")
		}
		result = req.Operands[0] / req.Operands[1]

	case OpPower:
		result = math.Pow(req.Operands[0], req.Operands[1])

	case OpSqrt:
		if req.Operands[0] < 0 {
			return nil, model.NewAppError(http.StatusUnprocessableEntity, model.ErrCodeNegativeSquareRoot, "Cannot calculate the square root of a negative number")
		}
		result = math.Sqrt(req.Operands[0])

	case OpPercentage:
		result = (req.Operands[0] / 100.0) * req.Operands[1]
	}

	return &model.CalculationResponse{
		Operation: req.Operation,
		Operands:  req.Operands,
		Result:    result,
	}, nil
}
