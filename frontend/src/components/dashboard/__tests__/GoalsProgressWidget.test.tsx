import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import GoalsProgressWidget from '../GoalsProgressWidget';
import type { LearningGoal } from '../../../api/goals';

const makeGoal = (overrides: Partial<LearningGoal> = {}): LearningGoal => ({
  id: 1,
  user_id: 1,
  title: 'TypeScript習得',
  description: '',
  category: 'language',
  target_date: null,
  progress: 60,
  status: 'active',
  created_at: '2026-02-19T00:00:00Z',
  updated_at: '2026-02-19T00:00:00Z',
  completed_at: null,
  ...overrides,
});

const renderWidget = (props: Partial<React.ComponentProps<typeof GoalsProgressWidget>> = {}) => {
  const defaultProps = {
    activeGoals: [] as LearningGoal[],
    completedGoals: [] as LearningGoal[],
    avgProgress: 0,
    loading: false,
    ...props,
  };
  return render(
    <MemoryRouter>
      <GoalsProgressWidget {...defaultProps} />
    </MemoryRouter>
  );
};

describe('GoalsProgressWidget', () => {
  it('ローディング中はスケルトンが表示される', () => {
    renderWidget({ loading: true });
    const skeletons = document.querySelectorAll('.animate-pulse');
    expect(skeletons.length).toBeGreaterThanOrEqual(2);
  });

  it('アクティブゴールが0件の場合は空メッセージが表示される', () => {
    renderWidget({ activeGoals: [] });
    expect(screen.getByText('進行中の目標がありません')).toBeInTheDocument();
    expect(screen.getByText('目標を作成する')).toBeInTheDocument();
  });

  it('ヘッダーに「学習目標の進捗」が表示される', () => {
    renderWidget();
    expect(screen.getByText('学習目標の進捗')).toBeInTheDocument();
  });

  it('「すべて見る」リンクが/goalsに遷移する', () => {
    renderWidget();
    const link = screen.getByText('すべて見る');
    expect(link.closest('a')).toHaveAttribute('href', '/goals');
  });

  it('統計行にアクティブ・完了・平均進捗が表示される', () => {
    const active = [makeGoal({ id: 1, progress: 40 }), makeGoal({ id: 2, progress: 80 })];
    const completed = [makeGoal({ id: 3, status: 'completed' })];
    renderWidget({ activeGoals: active, completedGoals: completed, avgProgress: 60 });
    expect(screen.getByText('2')).toBeInTheDocument();
    expect(screen.getByText('1')).toBeInTheDocument();
    expect(screen.getByText('60%')).toBeInTheDocument();
  });

  it('ゴールタイトルとプログレスが表示される', () => {
    const goals = [makeGoal({ id: 1, title: 'React学習', progress: 75 })];
    renderWidget({ activeGoals: goals, avgProgress: 75 });
    expect(screen.getByText('React学習')).toBeInTheDocument();
    expect(screen.getAllByText('75%').length).toBeGreaterThanOrEqual(1);
  });

  it('プログレスバーのaria属性が正しい', () => {
    const goals = [makeGoal({ id: 1, title: 'Go入門', progress: 50 })];
    renderWidget({ activeGoals: goals, avgProgress: 50 });
    const progressbar = screen.getByRole('progressbar');
    expect(progressbar).toHaveAttribute('aria-valuenow', '50');
    expect(progressbar).toHaveAttribute('aria-label', 'Go入門: 50%');
  });

  it('進捗80%以上はグリーン', () => {
    const goals = [makeGoal({ id: 1, progress: 85 })];
    renderWidget({ activeGoals: goals, avgProgress: 85 });
    const bar = screen.getByRole('progressbar').firstChild as HTMLElement;
    expect(bar.className).toContain('bg-green-500');
  });

  it('進捗50-79%はブルー', () => {
    const goals = [makeGoal({ id: 1, progress: 55 })];
    renderWidget({ activeGoals: goals, avgProgress: 55 });
    const bar = screen.getByRole('progressbar').firstChild as HTMLElement;
    expect(bar.className).toContain('bg-blue-500');
  });

  it('進捗50%未満はオレンジ', () => {
    const goals = [makeGoal({ id: 1, progress: 30 })];
    renderWidget({ activeGoals: goals, avgProgress: 30 });
    const bar = screen.getByRole('progressbar').firstChild as HTMLElement;
    expect(bar.className).toContain('bg-orange-500');
  });

  it('3件以下のゴールは「もっと見る」リンクが表示されない', () => {
    const goals = [
      makeGoal({ id: 1, title: 'Goal 1' }),
      makeGoal({ id: 2, title: 'Goal 2' }),
      makeGoal({ id: 3, title: 'Goal 3' }),
    ];
    renderWidget({ activeGoals: goals, avgProgress: 60 });
    expect(screen.queryByText(/他.*件の目標/)).not.toBeInTheDocument();
  });

  it('4件以上のゴールがある場合は「もっと見る」リンクが表示される', () => {
    const goals = [
      makeGoal({ id: 1, title: 'Goal A' }),
      makeGoal({ id: 2, title: 'Goal B' }),
      makeGoal({ id: 3, title: 'Goal C' }),
      makeGoal({ id: 4, title: 'Goal D' }),
      makeGoal({ id: 5, title: 'Goal E' }),
    ];
    renderWidget({ activeGoals: goals, avgProgress: 60 });
    // 最初の3件のみ表示
    expect(screen.getByText('Goal A')).toBeInTheDocument();
    expect(screen.getByText('Goal B')).toBeInTheDocument();
    expect(screen.getByText('Goal C')).toBeInTheDocument();
    expect(screen.queryByText('Goal D')).not.toBeInTheDocument();
  });
});
