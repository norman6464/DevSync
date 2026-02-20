import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import QuickStatsWidget from '../QuickStatsWidget';

const renderWidget = (props: { activeCount: number; completedCount: number }) =>
  render(
    <MemoryRouter>
      <QuickStatsWidget {...props} />
    </MemoryRouter>
  );

describe('QuickStatsWidget', () => {
  it('ヘッダーが表示される', () => {
    renderWidget({ activeCount: 0, completedCount: 0 });
    expect(screen.getByText('クイック統計')).toBeInTheDocument();
  });

  it('完了目標数が表示される', () => {
    renderWidget({ activeCount: 2, completedCount: 5 });
    expect(screen.getByText('5')).toBeInTheDocument();
    expect(screen.getByText('達成した目標')).toBeInTheDocument();
  });

  it('進行中目標数が表示される', () => {
    renderWidget({ activeCount: 3, completedCount: 1 });
    expect(screen.getByText('3')).toBeInTheDocument();
    expect(screen.getByText('進行中の目標')).toBeInTheDocument();
  });

  it('/goalsへのリンクが2つ表示される', () => {
    renderWidget({ activeCount: 0, completedCount: 0 });
    const links = screen.getAllByRole('link');
    expect(links).toHaveLength(2);
    links.forEach((link) => {
      expect(link).toHaveAttribute('href', '/goals');
    });
  });

  it('カウントが0の場合も正しく表示される', () => {
    renderWidget({ activeCount: 0, completedCount: 0 });
    const zeros = screen.getAllByText('0');
    expect(zeros).toHaveLength(2);
  });

  it('大きな数値も表示される', () => {
    renderWidget({ activeCount: 100, completedCount: 250 });
    expect(screen.getByText('100')).toBeInTheDocument();
    expect(screen.getByText('250')).toBeInTheDocument();
  });
});
