import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import GoalStatsPanel from '../GoalStatsPanel';

describe('GoalStatsPanel', () => {
  const defaultProps = {
    total: 10,
    active: 5,
    completed: 3,
    paused: 1,
    overdue: 1,
  };

  it('合計目標数を表示する', () => {
    render(<GoalStatsPanel {...defaultProps} />);
    expect(screen.getByText('10')).toBeInTheDocument();
    expect(screen.getByText('合計目標')).toBeInTheDocument();
  });

  it('進行中の目標数を表示する', () => {
    render(<GoalStatsPanel {...defaultProps} />);
    expect(screen.getByText('5')).toBeInTheDocument();
    expect(screen.getByText('進行中')).toBeInTheDocument();
  });

  it('完了の目標数を表示する', () => {
    render(<GoalStatsPanel {...defaultProps} />);
    expect(screen.getByText('完了')).toBeInTheDocument();
  });

  it('一時停止の目標数を表示する', () => {
    render(<GoalStatsPanel {...defaultProps} />);
    expect(screen.getByText('一時停止')).toBeInTheDocument();
  });

  it('期限超過の目標数を表示する', () => {
    render(<GoalStatsPanel {...defaultProps} />);
    expect(screen.getByText('期限超過')).toBeInTheDocument();
  });

  it('5つの統計カードを表示する', () => {
    const { container } = render(<GoalStatsPanel {...defaultProps} />);
    const cards = container.querySelectorAll('.bg-gray-900');
    expect(cards.length).toBe(5);
  });
});
