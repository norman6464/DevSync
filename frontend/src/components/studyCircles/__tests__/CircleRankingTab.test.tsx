import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import CircleRankingTab from '../CircleRankingTab';
import type { CircleMemberStreak } from '../../../types/studyCircle';

const mockStreaks: CircleMemberStreak[] = [
  { user_id: 1, user_name: 'Alice', avatar_url: '', current_streak: 10, total_checkins: 30 },
  { user_id: 2, user_name: 'Bob', avatar_url: '', current_streak: 5, total_checkins: 15 },
  { user_id: 3, user_name: 'Charlie', avatar_url: '', current_streak: 3, total_checkins: 8 },
];

describe('CircleRankingTab', () => {
  it('ストリークランキングのタイトルが表示される', () => {
    render(<CircleRankingTab streaks={mockStreaks} />);
    expect(screen.getByText('ストリークランキング')).toBeInTheDocument();
  });

  it('メンバー名が表示される', () => {
    render(<CircleRankingTab streaks={mockStreaks} />);
    expect(screen.getByText('Alice')).toBeInTheDocument();
    expect(screen.getByText('Bob')).toBeInTheDocument();
    expect(screen.getByText('Charlie')).toBeInTheDocument();
  });

  it('連続日数が表示される', () => {
    render(<CircleRankingTab streaks={mockStreaks} />);
    expect(screen.getByText('10')).toBeInTheDocument();
    expect(screen.getByText('5')).toBeInTheDocument();
  });

  it('順位番号が表示される', () => {
    render(<CircleRankingTab streaks={mockStreaks} />);
    const rankings = screen.getAllByText(/^[123]$/);
    expect(rankings.length).toBeGreaterThanOrEqual(3);
  });

  it('ストリークが空の場合は空状態メッセージが表示される', () => {
    render(<CircleRankingTab streaks={[]} />);
    expect(screen.getByText('データがまだありません')).toBeInTheDocument();
  });
});
