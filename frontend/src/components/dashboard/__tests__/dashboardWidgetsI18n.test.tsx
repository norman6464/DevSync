import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import QuickStatsWidget from '../QuickStatsWidget';
import GoalsProgressWidget from '../GoalsProgressWidget';
import RecentNotificationsWidget from '../RecentNotificationsWidget';
import type { LearningGoal } from '../../../api/goals';
import type { Notification } from '../../../types/notification';

const wrap = (ui: React.ReactElement) => render(<MemoryRouter>{ui}</MemoryRouter>);

/** 翻訳が引けなかったとき i18next はキー文字列をそのまま返すので、それが画面に出ていないか見る。 */
function rawKeys(container: HTMLElement): string[] {
  const text = container.textContent ?? '';
  return [...text.matchAll(/\b[a-z][a-zA-Z0-9]*\.[a-zA-Z0-9_.]+\b/g)]
    .map((m) => m[0])
    .filter((candidate) => !candidate.endsWith('...'));
}

const makeGoal = (overrides: Partial<LearningGoal> = {}): LearningGoal =>
  ({
    id: 1,
    user_id: 1,
    title: '目標',
    status: 'active',
    progress: 40,
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
    ...overrides,
  }) as LearningGoal;

const makeNotification = (overrides: Partial<Notification> = {}): Notification =>
  ({
    id: 1,
    type: 'like',
    is_read: false,
    created_at: '2026-01-01T00:00:00Z',
    actor: { id: 2, name: 'テスト太郎', username: 'tester' },
    ...overrides,
  }) as Notification;

describe('ダッシュボードのウィジェット', () => {
  it('クイック統計が翻訳キーではなく日本語を表示する', () => {
    const { container } = wrap(<QuickStatsWidget activeCount={2} completedCount={5} />);

    expect(rawKeys(container)).toEqual([]);
    expect(container.textContent).toContain('クイック統計');
    expect(container.textContent).toContain('達成した目標');
    expect(container.textContent).toContain('進行中の目標');
  });

  it('学習ゴール進捗が翻訳キーではなく日本語を表示する', () => {
    const activeGoals = [1, 2, 3, 4, 5].map((n) => makeGoal({ id: n, title: `目標${n}` }));
    const { container } = wrap(
      <GoalsProgressWidget
        activeGoals={activeGoals}
        completedGoals={[makeGoal({ id: 9, status: 'completed', progress: 100 })]}
        avgProgress={40}
        loading={false}
      />,
    );

    expect(rawKeys(container)).toEqual([]);
    expect(container.textContent).toContain('学習ゴール進捗');
    // 4 件目以降は「他 N 件の目標」に畳まれる。{{count}} が欠けると数値が消える。
    expect(container.textContent).toContain('他 2 件の目標');
  });

  it('目標が無いときも翻訳キーを表示しない', () => {
    const { container } = wrap(
      <GoalsProgressWidget activeGoals={[]} completedGoals={[]} avgProgress={0} loading={false} />,
    );

    expect(rawKeys(container)).toEqual([]);
    expect(container.textContent).toContain('進行中の目標はありません');
    expect(container.textContent).toContain('目標を作成');
  });

  it('最近の通知が翻訳キーではなく日本語を表示する', () => {
    const { container } = wrap(
      <RecentNotificationsWidget notifications={[makeNotification()]} loading={false} />,
    );

    expect(rawKeys(container)).toEqual([]);
    expect(container.textContent).toContain('最近の通知');
  });

  it('通知が無いときも翻訳キーを表示しない', () => {
    const { container } = wrap(<RecentNotificationsWidget notifications={[]} loading={false} />);

    expect(rawKeys(container)).toEqual([]);
    expect(container.textContent).toContain('通知はありません');
  });
});
