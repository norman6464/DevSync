import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import ProfileGoalsSection from '../ProfileGoalsSection';

const makeGoal = (overrides = {}) => ({
  id: 1,
  title: 'TypeScriptを学ぶ',
  description: '型安全なコードを書けるようになる',
  category: 'language',
  status: 'active',
  progress: 60,
  ...overrides,
});

const defaultProps = {
  goals: [makeGoal()],
  goalStats: { active_goals: 2, completed_goals: 1 },
  isOwnProfile: false,
};

const renderSection = (props = {}) =>
  render(
    <MemoryRouter>
      <ProfileGoalsSection {...defaultProps} {...props} />
    </MemoryRouter>
  );

describe('ProfileGoalsSection', () => {
  it('目標が空の場合何も表示しない', () => {
    const { container } = renderSection({ goals: [] });
    expect(container.innerHTML).toBe('');
  });

  it('目標タイトルが表示される', () => {
    renderSection();
    expect(screen.getByText('TypeScriptを学ぶ')).toBeInTheDocument();
  });

  it('目標の説明が表示される', () => {
    renderSection();
    expect(screen.getByText('型安全なコードを書けるようになる')).toBeInTheDocument();
  });

  it('進捗パーセントが表示される', () => {
    renderSection();
    expect(screen.getByText('60%')).toBeInTheDocument();
  });

  it('progressbarのaria属性が正しい', () => {
    renderSection();
    const progressbar = screen.getByRole('progressbar');
    expect(progressbar).toHaveAttribute('aria-valuenow', '60');
    expect(progressbar).toHaveAttribute('aria-valuemin', '0');
    expect(progressbar).toHaveAttribute('aria-valuemax', '100');
  });

  it('isOwnProfile=trueで目標管理リンクが表示される', () => {
    renderSection({ isOwnProfile: true });
    expect(screen.getByText('目標を管理')).toBeInTheDocument();
  });

  it('isOwnProfile=falseで目標管理リンクが非表示', () => {
    renderSection({ isOwnProfile: false });
    expect(screen.queryByText('目標を管理')).not.toBeInTheDocument();
  });

  it('最大4件まで表示される', () => {
    const goals = Array.from({ length: 6 }, (_, i) =>
      makeGoal({ id: i + 1, title: `目標${i + 1}` })
    );
    renderSection({ goals });
    expect(screen.getByText('目標1')).toBeInTheDocument();
    expect(screen.getByText('目標4')).toBeInTheDocument();
    expect(screen.queryByText('目標5')).not.toBeInTheDocument();
  });
});
