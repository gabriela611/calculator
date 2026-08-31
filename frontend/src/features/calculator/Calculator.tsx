import React, { useEffect } from 'react';
import { useCalculator } from './useCalculator';
import type { Operation } from '../../api/types';
import './Calculator.css';

const OP_SYMBOLS: Record<Operation, string> = {
  add: '+',
  subtract: '-',
  multiply: '×',
  divide: '÷',
  power: '^',
  sqrt: '√',
  percentage: '%',
};

export const Calculator: React.FC = () => {
  const {
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
  } = useCalculator();

  // Keyboard shortcut listener
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key >= '0' && e.key <= '9') {
        inputDigit(e.key);
      } else if (e.key === '.' || e.key === ',') {
        inputDecimal();
      } else if (e.key === '+') {
        selectOperation('add');
      } else if (e.key === '-') {
        selectOperation('subtract');
      } else if (e.key === '*') {
        selectOperation('multiply');
      } else if (e.key === '/') {
        e.preventDefault();
        selectOperation('divide');
      } else if (e.key === '^') {
        selectOperation('power');
      } else if (e.key === '%') {
        selectOperation('percentage');
      } else if (e.key === 'Enter' || e.key === '=') {
        e.preventDefault();
        executeEquals();
      } else if (e.key === 'Escape' || e.key === 'c' || e.key === 'C') {
        clear();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [inputDigit, inputDecimal, selectOperation, executeEquals, clear]);

  // Expression preview string (e.g., "12 +" or "5 ^")
  const getExpressionPreview = (): string => {
    if (storedValue !== null && pendingOp) {
      return `${storedValue} ${OP_SYMBOLS[pendingOp] || ''}`;
    }
    return '';
  };

  return (
    <div className="calculator-card" role="region" aria-label="Calculator">
      {/* Top Display Section */}
      <div className="calculator-display-container">
        <div className="expression-preview" aria-live="polite">
          {getExpressionPreview()}
        </div>
        <div className="primary-display" aria-live="polite" tabIndex={0}>
          {display}
        </div>
        {isLoading && (
          <div className="loading-bar-container">
            <div className="loading-bar-fill"></div>
          </div>
        )}
        {error && (
          <div className="error-toast" role="alert">
            {error.message}
          </div>
        )}
      </div>

      {/* Keypad Grid - 4x5 Layout (Google Calc style) */}
      <div className="keypad-grid">
        {/* Row 1: Function keys & Division */}
        <button
          type="button"
          className="calc-btn btn-function"
          onClick={clear}
          aria-label="Clear display"
        >
          C
        </button>
        <button
          type="button"
          className={`calc-btn btn-function ${pendingOp === 'power' ? 'active-op' : ''}`}
          onClick={() => selectOperation('power')}
          aria-label="Power operation"
        >
          ^
        </button>
        <button
          type="button"
          className={`calc-btn btn-function ${pendingOp === 'percentage' ? 'active-op' : ''}`}
          onClick={() => selectOperation('percentage')}
          aria-label="Percentage operation"
        >
          %
        </button>
        <button
          type="button"
          className={`calc-btn btn-operator ${pendingOp === 'divide' ? 'active-op' : ''}`}
          onClick={() => selectOperation('divide')}
          aria-label="Divide"
        >
          ÷
        </button>

        {/* Row 2: 7, 8, 9, Multiply */}
        <button
          type="button"
          className="calc-btn btn-number"
          onClick={() => inputDigit('7')}
        >
          7
        </button>
        <button
          type="button"
          className="calc-btn btn-number"
          onClick={() => inputDigit('8')}
        >
          8
        </button>
        <button
          type="button"
          className="calc-btn btn-number"
          onClick={() => inputDigit('9')}
        >
          9
        </button>
        <button
          type="button"
          className={`calc-btn btn-operator ${pendingOp === 'multiply' ? 'active-op' : ''}`}
          onClick={() => selectOperation('multiply')}
          aria-label="Multiply"
        >
          ×
        </button>

        {/* Row 3: 4, 5, 6, Subtract */}
        <button
          type="button"
          className="calc-btn btn-number"
          onClick={() => inputDigit('4')}
        >
          4
        </button>
        <button
          type="button"
          className="calc-btn btn-number"
          onClick={() => inputDigit('5')}
        >
          5
        </button>
        <button
          type="button"
          className="calc-btn btn-number"
          onClick={() => inputDigit('6')}
        >
          6
        </button>
        <button
          type="button"
          className={`calc-btn btn-operator ${pendingOp === 'subtract' ? 'active-op' : ''}`}
          onClick={() => selectOperation('subtract')}
          aria-label="Subtract"
        >
          −
        </button>

        {/* Row 4: 1, 2, 3, Add */}
        <button
          type="button"
          className="calc-btn btn-number"
          onClick={() => inputDigit('1')}
        >
          1
        </button>
        <button
          type="button"
          className="calc-btn btn-number"
          onClick={() => inputDigit('2')}
        >
          2
        </button>
        <button
          type="button"
          className="calc-btn btn-number"
          onClick={() => inputDigit('3')}
        >
          3
        </button>
        <button
          type="button"
          className={`calc-btn btn-operator ${pendingOp === 'add' ? 'active-op' : ''}`}
          onClick={() => selectOperation('add')}
          aria-label="Add"
        >
          +
        </button>

        {/* Row 5: Square Root, 0, Decimal, Equals */}
        <button
          type="button"
          className="calc-btn btn-function"
          onClick={() => selectOperation('sqrt')}
          aria-label="Square root"
        >
          √
        </button>
        <button
          type="button"
          className="calc-btn btn-number"
          onClick={() => inputDigit('0')}
        >
          0
        </button>
        <button
          type="button"
          className="calc-btn btn-number"
          onClick={inputDecimal}
          aria-label="Decimal point"
        >
          .
        </button>
        <button
          type="button"
          className="calc-btn btn-equals"
          onClick={executeEquals}
          aria-label="Equals"
        >
          =
        </button>
      </div>
    </div>
  );
};
