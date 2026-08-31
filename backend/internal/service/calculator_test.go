package service_test

import (
	"net/http"
	"testing"

	"backend/internal/model"
	"backend/internal/service"
)

func TestCalculatorService_Calculate(t *testing.T) {
	svc := service.NewCalculatorService()

	tests := []struct {
		name           string
		request        model.CalculationRequest
		wantResult     float64
		wantErrCode    string
		wantStatusCode int
	}{
		// Valid calculations
		{
			name:       "Addition",
			request:    model.CalculationRequest{Operation: "add", Operands: []float64{10, 5}},
			wantResult: 15,
		},
		{
			name:       "Subtraction",
			request:    model.CalculationRequest{Operation: "subtract", Operands: []float64{10, 5}},
			wantResult: 5,
		},
		{
			name:       "Multiplication",
			request:    model.CalculationRequest{Operation: "multiply", Operands: []float64{10, 5}},
			wantResult: 50,
		},
		{
			name:       "Division",
			request:    model.CalculationRequest{Operation: "divide", Operands: []float64{10, 5}},
			wantResult: 2,
		},
		{
			name:       "Power",
			request:    model.CalculationRequest{Operation: "power", Operands: []float64{2, 3}},
			wantResult: 8,
		},
		{
			name:       "Square Root",
			request:    model.CalculationRequest{Operation: "sqrt", Operands: []float64{25}},
			wantResult: 5,
		},
		{
			name:       "Percentage",
			request:    model.CalculationRequest{Operation: "percentage", Operands: []float64{20, 150}},
			wantResult: 30,
		},

		// Validation errors
		{
			name:           "Missing operation",
			request:        model.CalculationRequest{Operation: "", Operands: []float64{10, 5}},
			wantErrCode:    model.ErrCodeMissingOperation,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "Missing operands",
			request:        model.CalculationRequest{Operation: "add", Operands: nil},
			wantErrCode:    model.ErrCodeMissingOperands,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "Invalid operation",
			request:        model.CalculationRequest{Operation: "modulo", Operands: []float64{10, 5}},
			wantErrCode:    model.ErrCodeInvalidOperation,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "Invalid operand count for add (too few)",
			request:        model.CalculationRequest{Operation: "add", Operands: []float64{10}},
			wantErrCode:    model.ErrCodeInvalidOperandCount,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "Invalid operand count for sqrt (too many)",
			request:        model.CalculationRequest{Operation: "sqrt", Operands: []float64{25, 5}},
			wantErrCode:    model.ErrCodeInvalidOperandCount,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "Division by zero",
			request:        model.CalculationRequest{Operation: "divide", Operands: []float64{10, 0}},
			wantErrCode:    model.ErrCodeDivisionByZero,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "Negative square root",
			request:        model.CalculationRequest{Operation: "sqrt", Operands: []float64{-25}},
			wantErrCode:    model.ErrCodeNegativeSquareRoot,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, appErr := svc.Calculate(tt.request)

			if tt.wantErrCode != "" {
				if appErr == nil {
					t.Fatalf("expected error %s, got nil response", tt.wantErrCode)
				}
				if appErr.Code != tt.wantErrCode {
					t.Errorf("expected error code %s, got %s", tt.wantErrCode, appErr.Code)
				}
				if appErr.StatusCode != tt.wantStatusCode {
					t.Errorf("expected HTTP status code %d, got %d", tt.wantStatusCode, appErr.StatusCode)
				}
			} else {
				if appErr != nil {
					t.Fatalf("unexpected error: %v", appErr)
				}
				if resp.Result != tt.wantResult {
					t.Errorf("expected result %v, got %v", tt.wantResult, resp.Result)
				}
			}
		})
	}
}
