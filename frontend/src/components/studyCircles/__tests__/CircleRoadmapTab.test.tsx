import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import CircleRoadmapTab from '../CircleRoadmapTab';
import type { StudyCircle, StudyCircleMemberProgress } from '../../../types/studyCircle';

const mockCircle: StudyCircle = {
  id: 1,
  name: 'React学習会',
  topic: 'React',
  description: '',
  owner_id: 1,
  max_members: 10,
  status: 'active',
  steps: [
    { id: 1, circle_id: 1, title: 'JSX基礎', description: 'JSXの書き方を学ぶ', order_index: 0, resource_url: 'https://example.com', created_at: '', updated_at: '' },
    { id: 2, circle_id: 1, title: 'Hooks入門', description: '', order_index: 1, resource_url: '', created_at: '', updated_at: '' },
  ],
  members: [
    { id: 1, circle_id: 1, user_id: 1, user: { id: 1, name: 'Alice', avatar_url: '' }, role: 'owner', joined_at: '' },
  ],
  created_at: '',
  updated_at: '',
};

const mockProgress: StudyCircleMemberProgress[] = [
  { id: 1, circle_id: 1, step_id: 1, user_id: 1, is_completed: true, completed_at: '2026-02-18T10:00:00Z' },
];

const mockUser = { id: 1, name: 'Alice', username: 'alice', email: 'a@e.com', bio: '', avatar_url: '', github_username: '', created_at: '', updated_at: '' };

describe('CircleRoadmapTab', () => {
  const defaultProps = {
    circle: mockCircle,
    progress: mockProgress,
    currentUser: mockUser,
    isOwner: true,
    saving: false,
    onCreateStep: vi.fn(),
    onDeleteStep: vi.fn(),
    onToggleProgress: vi.fn(),
  };

  it('ステップタイトルが表示される', () => {
    render(<CircleRoadmapTab {...defaultProps} />);
    expect(screen.getByText('JSX基礎')).toBeInTheDocument();
    expect(screen.getByText('Hooks入門')).toBeInTheDocument();
  });

  it('ステップの説明が表示される', () => {
    render(<CircleRoadmapTab {...defaultProps} />);
    expect(screen.getByText('JSXの書き方を学ぶ')).toBeInTheDocument();
  });

  it('完了済みステップに取り消し線が適用される', () => {
    render(<CircleRoadmapTab {...defaultProps} />);
    const completedStep = screen.getByText('JSX基礎');
    expect(completedStep).toHaveClass('line-through');
  });

  it('リソースURLリンクが表示される', () => {
    render(<CircleRoadmapTab {...defaultProps} />);
    expect(screen.getByText('参考URL')).toBeInTheDocument();
  });

  it('オーナーにステップ追加ボタンが表示される', () => {
    render(<CircleRoadmapTab {...defaultProps} />);
    expect(screen.getByText('ステップを追加')).toBeInTheDocument();
  });

  it('非オーナーにステップ追加ボタンが非表示', () => {
    render(<CircleRoadmapTab {...defaultProps} isOwner={false} />);
    expect(screen.queryByText('ステップを追加')).not.toBeInTheDocument();
  });

  it('ステップ追加ボタンクリックでフォームが表示される', () => {
    render(<CircleRoadmapTab {...defaultProps} />);
    fireEvent.click(screen.getByText('ステップを追加'));
    expect(screen.getByPlaceholderText('ステップ')).toBeInTheDocument();
  });

  it('ステップがない場合は空状態メッセージが表示される', () => {
    const emptyCircle = { ...mockCircle, steps: [] };
    render(<CircleRoadmapTab {...defaultProps} circle={emptyCircle} />);
    expect(screen.getByText('ステップがまだありません')).toBeInTheDocument();
  });
});
