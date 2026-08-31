export type Operation = 'add' | 'subtract' | 'multiply' | 'divide' | 'power' | 'sqrt' | 'percentage';

export interface CalculationRequest {
  operation: Operation | string;
  operands: number[];
}

export interface CalculationResponse {
  operation: string;
  operands: number[];
  result: number;
}

export interface ErrorDetail {
  code: string;
  message: string;
}

export interface ErrorResponse {
  error: ErrorDetail;
}

export type ApiResult<T> =
  | { success: true; data: T }
  | { success: false; error: ErrorDetail; statusCode: number };
