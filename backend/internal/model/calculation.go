package model

// CalculationRequest represents the payload for a calculation request.
type CalculationRequest struct {
	Operation string    `json:"operation"`
	Operands  []float64 `json:"operands"`
}

// CalculationResponse represents the successful response of a calculation.
type CalculationResponse struct {
	Operation string    `json:"operation"`
	Operands  []float64 `json:"operands"`
	Result    float64   `json:"result"`
}

// ErrorDetail contains the specific error code and human-readable message.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse represents the standardized error response payload.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}
