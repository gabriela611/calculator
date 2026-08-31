import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { Calculator } from './Calculator';
import * as calculatorApi from '../../api/calculatorApi';

vi.mock('../../api/calculatorApi');

describe('Calculator Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders correctly with initial display 0', () => {
    const { container } = render(<Calculator />);
    const display = container.querySelector('.primary-display');
    expect(display).toHaveTextContent('0');
  });

  it('handles chained binary operation flow (5 + 3 + 2 = 10)', async () => {
    vi.mocked(calculatorApi.calculate)
      .mockResolvedValueOnce({
        success: true,
        data: { operation: 'add', operands: [5, 3], result: 8 },
      })
      .mockResolvedValueOnce({
        success: true,
        data: { operation: 'add', operands: [8, 2], result: 10 },
      });

    render(<Calculator />);

    fireEvent.click(screen.getByRole('button', { name: '5' }));
    fireEvent.click(screen.getByRole('button', { name: 'Add' }));
    fireEvent.click(screen.getByRole('button', { name: '3' }));
    fireEvent.click(screen.getByRole('button', { name: 'Add' })); // Triggers first calc: 5 + 3 = 8

    await waitFor(() => {
      expect(calculatorApi.calculate).toHaveBeenLastCalledWith({
        operation: 'add',
        operands: [5, 3],
      });
    });

    fireEvent.click(screen.getByRole('button', { name: '2' }));
    fireEvent.click(screen.getByRole('button', { name: 'Equals' })); // Triggers second calc: 8 + 2 = 10

    await waitFor(() => {
      expect(calculatorApi.calculate).toHaveBeenLastCalledWith({
        operation: 'add',
        operands: [8, 2],
      });
    });
  });

  it('handles unary sqrt operation flow (sqrt(16) = 4)', async () => {
    vi.mocked(calculatorApi.calculate).mockResolvedValueOnce({
      success: true,
      data: { operation: 'sqrt', operands: [16], result: 4 },
    });

    render(<Calculator />);

    fireEvent.click(screen.getByRole('button', { name: '1' }));
    fireEvent.click(screen.getByRole('button', { name: '6' }));
    fireEvent.click(screen.getByRole('button', { name: 'Square root' }));

    await waitFor(() => {
      expect(calculatorApi.calculate).toHaveBeenCalledWith({
        operation: 'sqrt',
        operands: [16],
      });
    });
  });

  it('handles division by zero error gracefully', async () => {
    vi.mocked(calculatorApi.calculate).mockResolvedValueOnce({
      success: false,
      error: { code: 'DIVISION_BY_ZERO', message: 'Cannot divide by zero' },
      statusCode: 422,
    });

    const { container } = render(<Calculator />);
    const display = container.querySelector('.primary-display');

    fireEvent.click(screen.getByRole('button', { name: '5' }));
    fireEvent.click(screen.getByRole('button', { name: 'Divide' }));
    fireEvent.click(screen.getByRole('button', { name: '0' }));
    fireEvent.click(screen.getByRole('button', { name: 'Equals' }));

    await waitFor(() => {
      expect(screen.getByText('Cannot divide by zero')).toBeInTheDocument();
      expect(display).toHaveTextContent('Error');
    });
  });

  it('clears state when C button is pressed', async () => {
    const { container } = render(<Calculator />);
    const display = container.querySelector('.primary-display');

    fireEvent.click(screen.getByRole('button', { name: '9' }));
    fireEvent.click(screen.getByRole('button', { name: '9' }));
    expect(display).toHaveTextContent('99');

    fireEvent.click(screen.getByRole('button', { name: 'Clear display' }));
    expect(display).toHaveTextContent('0');
  });
});
