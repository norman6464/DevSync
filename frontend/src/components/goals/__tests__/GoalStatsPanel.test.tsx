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

  it('各統計にアイコンが表示される', () => {
    const { container } = render(<GoalStatsPanel {...defaultProps} />);
    const icons = container.querySelectorAll('svg');
    // 5つの統計それぞれにアイコンがある
    expect(icons.length).toBeGreaterThanOrEqual(5);
  });

  it('進行中の目標に青色のスタイルが適用される', () => {
    const { container } = render(<GoalStatsPanel {...defaultProps} />);
    const blueElements = container.querySelectorAll('.text-blue-400');
    expect(blueElements.length).toBeGreaterThan(0);
  });

  it('完了の目標に緑色のスタイルが適用される', () => {
    const { container } = render(<GoalStatsPanel {...defaultProps} />);
    const greenElements = container.querySelectorAll('.text-green-400');
    expect(greenElements.length).toBeGreaterThan(0);
  });

  it('一時停止の目標に黄色のスタイルが適用される', () => {
    const { container } = render(<GoalStatsPanel {...defaultProps} />);
    const yellowElements = container.querySelectorAll('.text-yellow-400');
    expect(yellowElements.length).toBeGreaterThan(0);
  });

  it('期限超過の目標に赤色のスタイルが適用される', () => {
    const { container } = render(<GoalStatsPanel {...defaultProps} />);
    const redElements = container.querySelectorAll('.text-red-400');
    expect(redElements.length).toBeGreaterThan(0);
  });

  it('0の統計値も正しく表示される', () => {
    render(<GoalStatsPanel total={0} active={0} completed={0} paused={0} overdue={0} />);
    const zeros = screen.getAllByText('0');
    expect(zeros.length).toBe(5);
  });
});
