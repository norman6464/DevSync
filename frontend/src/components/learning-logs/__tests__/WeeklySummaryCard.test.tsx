import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import WeeklySummaryCard from '../WeeklySummaryCard';
import type { StreakInfo } from '../../../types/learningLog';

const baseStreakInfo: StreakInfo = {
  current_streak: 7,
  longest_streak: 14,
  total_days: 30,
  last_log_date: '2026-02-19',
};

describe('WeeklySummaryCard', () => {
  it('週間学習時間が分単位で表示される（60分未満）', () => {
    render(
      <WeeklySummaryCard weeklyDuration={45} streakInfo={baseStreakInfo} logCount={10} />
    );
    expect(screen.getByText('45分')).toBeInTheDocument();
  });

  it('週間学習時間が時間+分で表示される（60分以上）', () => {
    render(
      <WeeklySummaryCard weeklyDuration={125} streakInfo={baseStreakInfo} logCount={10} />
    );
    expect(screen.getByText(/2.*5/)).toBeInTheDocument();
  });

  it('週間学習時間0分が表示される', () => {
    render(
      <WeeklySummaryCard weeklyDuration={0} streakInfo={baseStreakInfo} logCount={0} />
    );
    expect(screen.getByText('0分')).toBeInTheDocument();
  });

  it('連続日数が表示される', () => {
    render(
      <WeeklySummaryCard weeklyDuration={30} streakInfo={baseStreakInfo} logCount={5} />
    );
    expect(screen.getByText(/7/)).toBeInTheDocument();
  });

  it('streakInfoがnullの場合0日が表示される', () => {
    render(
      <WeeklySummaryCard weeklyDuration={30} streakInfo={null} logCount={5} />
    );
    expect(screen.getByText(/^0/)).toBeInTheDocument();
  });

  it('ログ数が表示される', () => {
    render(
      <WeeklySummaryCard weeklyDuration={30} streakInfo={baseStreakInfo} logCount={42} />
    );
    expect(screen.getByText('42')).toBeInTheDocument();
  });

  it('ラベルが表示される', () => {
    render(
      <WeeklySummaryCard weeklyDuration={30} streakInfo={baseStreakInfo} logCount={5} />
    );
    expect(screen.getByText('今週の学習時間')).toBeInTheDocument();
    expect(screen.getByText('連続学習日数')).toBeInTheDocument();
    expect(screen.getByText('ログ件数')).toBeInTheDocument();
  });
});
