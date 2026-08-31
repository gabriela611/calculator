import type { CalculationRequest, CalculationResponse, ErrorResponse, ApiResult } from './types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

/**
 * Sends a calculation request to the backend REST API.
 */
export async function calculate(request: CalculationRequest): Promise<ApiResult<CalculationResponse>> {
  try {
    const response = await fetch(`${API_BASE_URL}/api/v1/calculate`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });

    const data = await response.json();

    if (!response.ok) {
      const errorData = data as ErrorResponse;
      return {
        success: false,
        error: errorData.error || { code: 'UNKNOWN_ERROR', message: 'An unknown error occurred' },
        statusCode: response.status,
      };
    }

    return {
      success: true,
      data: data as CalculationResponse,
    };
  } catch (err) {
    return {
      success: false,
      error: {
        code: 'NETWORK_ERROR',
        message: err instanceof Error ? err.message : 'Network connection error',
      },
      statusCode: 0,
    };
  }
}
