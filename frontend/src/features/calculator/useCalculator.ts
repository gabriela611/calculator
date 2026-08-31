import { useState, useCallback } from 'react';
import type { Operation, ErrorDetail } from '../../api/types';
import { calculate } from '../../api/calculatorApi';

export interface UseCalculatorReturn {
  display: string;
  storedValue: number | null;
  pendingOp: Operation | null;
  isLoading: boolean;
  error: ErrorDetail | null;
  inputDigit: (digit: string) => void;
  inputDecimal: () => void;
  clear: () => void;
  selectOperation: (op: Operation) => Promise<void>;
  executeEquals: () => Promise<void>;
}

export function useCalculator(): UseCalculatorReturn {
  const [display, setDisplay] = useState<string>('0');
  const [storedValue, setStoredValue] = useState<number | null>(null);
  const [pendingOp, setPendingOp] = useState<Operation | null>(null);
  const [justCalculated, setJustCalculated] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<ErrorDetail | null>(null);

  const inputDigit = useCallback((digit: string) => {
    setError(null);
    if (justCalculated) {
      setDisplay(digit);
      setJustCalculated(false);
      return;
    }

    setDisplay((prev) => (prev === '0' ? digit : prev + digit));
  }, [justCalculated]);

  const inputDecimal = useCallback(() => {
    setError(null);
    if (justCalculated) {
      setDisplay('0.');
      setJustCalculated(false);
      return;
    }

    if (!display.includes('.')) {
      setDisplay((prev) => prev + '.');
    }
  }, [display, justCalculated]);

  const clear = useCallback(() => {
    setDisplay('0');
    setStoredValue(null);
    setPendingOp(null);
    setJustCalculated(false);
    setIsLoading(false);
    setError(null);
  }, []);

  const runCalculation = async (
    operation: Operation,
    operands: number[]
  ): Promise<number | null> => {
    setIsLoading(true);
    setError(null);

    const res = await calculate({ operation, operands });
    setIsLoading(false);

    if (res.success) {
      setDisplay(String(res.data.result));
      setError(null);
      return res.data.result;
    } else {
      setError(res.error);
      setDisplay('Error');
      return null;
    }
  };

  const selectOperation = useCallback(async (op: Operation) => {
    const currentNum = parseFloat(display);
    if (isNaN(currentNum)) return;

    // Unary operation: sqrt triggers immediately on current value
    if (op === 'sqrt') {
      const result = await runCalculation('sqrt', [currentNum]);
      if (result !== null) {
        setJustCalculated(true);
      }
      return;
    }

    // Binary operations: chained operation resolution if pendingOp and storedValue exist
    if (pendingOp && storedValue !== null && !justCalculated) {
      const result = await runCalculation(pendingOp, [storedValue, currentNum]);
      if (result !== null) {
        setStoredValue(result);
        setPendingOp(op);
        setJustCalculated(true);
      }
      return;
    }

    // Normal binary operation registration
    setStoredValue(currentNum);
    setPendingOp(op);
    setJustCalculated(true);
  }, [display, pendingOp, storedValue, justCalculated]);

  const executeEquals = useCallback(async () => {
    if (!pendingOp || storedValue === null) return;

    const currentNum = parseFloat(display);
    if (isNaN(currentNum)) return;

    const result = await runCalculation(pendingOp, [storedValue, currentNum]);
    if (result !== null) {
      setStoredValue(null);
      setPendingOp(null);
      setJustCalculated(true);
    }
  }, [pendingOp, storedValue, display]);

  return {
    display,
    storedValue,
    pendingOp,
    isLoading,
    error,
    inputDigit,
    inputDecimal,
    clear,
    selectOperation,
    executeEquals,
  };
}
