package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"backend/internal/model"
	"backend/internal/service"
)

// CalculatorHandler handles HTTP requests for calculations.
type CalculatorHandler struct {
	svc service.CalculatorService
}

// NewCalculatorHandler initializes a new CalculatorHandler.
func NewCalculatorHandler(svc service.CalculatorService) *CalculatorHandler {
	return &CalculatorHandler{svc: svc}
}

// Calculate Handles POST /api/v1/calculate requests.
func (h *CalculatorHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only POST requests are allowed")
		return
	}

	var req model.CalculationRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // strict parsing

	if err := decoder.Decode(&req); err != nil {
		var typeErr *json.UnmarshalTypeError
		var syntaxErr *json.SyntaxError

		if errors.As(err, &typeErr) {
			writeError(w, http.StatusBadRequest, model.ErrCodeInvalidOperand, "All operands must be numeric values")
			return
		}

		if errors.As(err, &syntaxErr) || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
			writeError(w, http.StatusBadRequest, model.ErrCodeInvalidJSON, "Request body contains malformed JSON")
			return
		}

		// Fallback for malformed JSON syntax or unknown structure errors
		writeError(w, http.StatusBadRequest, model.ErrCodeInvalidJSON, "Request body contains malformed JSON")
		return
	}

	res, appErr := h.svc.Calculate(req)
	if appErr != nil {
		writeError(w, appErr.StatusCode, appErr.Code, appErr.Message)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func writeError(w http.ResponseWriter, statusCode int, code, message string) {
	w.WriteHeader(statusCode)
	errResp := model.ErrorResponse{
		Error: model.ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
	_ = json.NewEncoder(w).Encode(errResp)
}
