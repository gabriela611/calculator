package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/handler"
	"backend/internal/model"
	"backend/internal/service"
)

func TestCalculatorHandler_Scenarios(t *testing.T) {
	svc := service.NewCalculatorService()
	router := handler.NewRouter(svc)

	tests := []struct {
		name           string
		method         string
		payload        string
		wantStatusCode int
		wantErrorCode  string
		wantResult     *float64
	}{
		// 200 OK
		{
			name:           "Valid Add",
			method:         http.MethodPost,
			payload:        `{"operation":"add","operands":[10,5]}`,
			wantStatusCode: http.StatusOK,
			wantResult:     floatPtr(15),
		},
		{
			name:           "Valid Sqrt",
			method:         http.MethodPost,
			payload:        `{"operation":"sqrt","operands":[25]}`,
			wantStatusCode: http.StatusOK,
			wantResult:     floatPtr(5),
		},

		// 400 Bad Request
		{
			name:           "Invalid Operation",
			method:         http.MethodPost,
			payload:        `{"operation":"modulo","operands":[10,5]}`,
			wantStatusCode: http.StatusBadRequest,
			wantErrorCode:  model.ErrCodeInvalidOperation,
		},
		{
			name:           "Missing Operation",
			method:         http.MethodPost,
			payload:        `{"operands":[10,5]}`,
			wantStatusCode: http.StatusBadRequest,
			wantErrorCode:  model.ErrCodeMissingOperation,
		},
		{
			name:           "Missing Operands",
			method:         http.MethodPost,
			payload:        `{"operation":"add"}`,
			wantStatusCode: http.StatusBadRequest,
			wantErrorCode:  model.ErrCodeMissingOperands,
		},
		{
			name:           "Invalid Operand (non-numeric)",
			method:         http.MethodPost,
			payload:        `{"operation":"add","operands":[10,"hello"]}`,
			wantStatusCode: http.StatusBadRequest,
			wantErrorCode:  model.ErrCodeInvalidOperand,
		},
		{
			name:           "Invalid Operand Count",
			method:         http.MethodPost,
			payload:        `{"operation":"sqrt","operands":[25,5]}`,
			wantStatusCode: http.StatusBadRequest,
			wantErrorCode:  model.ErrCodeInvalidOperandCount,
		},
		{
			name:           "Invalid JSON (malformed)",
			method:         http.MethodPost,
			payload:        `{"operation":"add","operands":[10,5`,
			wantStatusCode: http.StatusBadRequest,
			wantErrorCode:  model.ErrCodeInvalidJSON,
		},

		// 422 Unprocessable Entity
		{
			name:           "Division by zero",
			method:         http.MethodPost,
			payload:        `{"operation":"divide","operands":[10,0]}`,
			wantStatusCode: http.StatusUnprocessableEntity,
			wantErrorCode:  model.ErrCodeDivisionByZero,
		},
		{
			name:           "Negative square root",
			method:         http.MethodPost,
			payload:        `{"operation":"sqrt","operands":[-25]}`,
			wantStatusCode: http.StatusUnprocessableEntity,
			wantErrorCode:  model.ErrCodeNegativeSquareRoot,
		},

		// Method Not Allowed
		{
			name:           "GET Method Not Allowed",
			method:         http.MethodGet,
			payload:        "",
			wantStatusCode: http.StatusMethodNotAllowed,
			wantErrorCode:  "METHOD_NOT_ALLOWED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/v1/calculate", bytes.NewBufferString(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Fatalf("expected HTTP status %d, got %d. Body: %s", tt.wantStatusCode, rec.Code, rec.Body.String())
			}

			if tt.wantErrorCode != "" {
				var errResp model.ErrorResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
					t.Fatalf("failed to unmarshal error response: %v", err)
				}
				if errResp.Error.Code != tt.wantErrorCode {
					t.Errorf("expected error code %s, got %s", tt.wantErrorCode, errResp.Error.Code)
				}
			}

			if tt.wantResult != nil {
				var succResp model.CalculationResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &succResp); err != nil {
					t.Fatalf("failed to unmarshal success response: %v", err)
				}
				if succResp.Result != *tt.wantResult {
					t.Errorf("expected result %v, got %v", *tt.wantResult, succResp.Result)
				}
			}
		})
	}
}

func TestCORS_Preflight(t *testing.T) {
	svc := service.NewCalculatorService()
	router := handler.NewRouter(svc)

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/calculate", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected CORS preflight status 204, got %d", rec.Code)
	}

	if origin := rec.Header().Get("Access-Control-Allow-Origin"); origin != "*" {
		t.Errorf("expected CORS origin '*', got %q", origin)
	}
}

func floatPtr(v float64) *float64 {
	return &v
}
