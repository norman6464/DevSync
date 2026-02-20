import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import RoadmapCard from '../RoadmapCard';
import { type Roadmap } from '../../../api/roadmaps';

const makeRoadmap = (overrides: Partial<Roadmap> = {}): Roadmap => ({
  id: 1,
  user_id: 1,
  title: 'React学習ロードマップ',
  description: 'Reactの基礎から応用まで',
  category: 'framework',
  is_public: false,
  is_template: false,
  step_count: 10,
  completed_step_count: 3,
  progress: 30,
  status: 'active',
  created_at: '2025-01-01',
  updated_at: '2025-01-10',
  completed_at: null,
  ...overrides,
});

const defaultProps = {
  roadmap: makeRoadmap(),
  onView: vi.fn(),
  onEdit: vi.fn(),
  onDelete: vi.fn(),
};

const renderCard = (props = {}) =>
  render(<RoadmapCard {...defaultProps} {...props} />);

describe('RoadmapCard', () => {
  it('タイトルが表示される', () => {
    renderCard();
    expect(screen.getByText('React学習ロードマップ')).toBeInTheDocument();
  });

  it('説明が表示される', () => {
    renderCard();
    expect(screen.getByText('Reactの基礎から応用まで')).toBeInTheDocument();
  });

  it('進捗パーセンテージが表示される', () => {
    renderCard();
    expect(screen.getByText('30%')).toBeInTheDocument();
  });

  it('ステップ数が表示される', () => {
    renderCard();
    expect(screen.getByText(/3 \/ 10/)).toBeInTheDocument();
  });

  it('公開バッジが表示される（is_public=true）', () => {
    renderCard({ roadmap: makeRoadmap({ is_public: true }) });
    expect(screen.getByText('公開')).toBeInTheDocument();
  });

  it('完了バッジが表示される（status=completed）', () => {
    renderCard({ roadmap: makeRoadmap({ status: 'completed' }) });
    expect(screen.getByText('完了')).toBeInTheDocument();
  });

  it('カードクリックでonViewが呼ばれる', () => {
    const onView = vi.fn();
    renderCard({ onView });
    fireEvent.click(screen.getByText('React学習ロードマップ').closest('div[class*="cursor-pointer"]')!);
    expect(onView).toHaveBeenCalled();
  });

  it('編集ボタンクリックでonEditが呼ばれる', () => {
    const onEdit = vi.fn();
    renderCard({ onEdit });
    const buttons = screen.getAllByRole('button');
    fireEvent.click(buttons[0]);
    expect(onEdit).toHaveBeenCalledWith(expect.objectContaining({ id: 1 }));
  });

  it('削除ボタンクリックでonDeleteが呼ばれる', () => {
    const onDelete = vi.fn();
    renderCard({ onDelete });
    const buttons = screen.getAllByRole('button');
    fireEvent.click(buttons[1]);
    expect(onDelete).toHaveBeenCalledWith(1);
  });

  it('説明がない場合はpタグが非表示', () => {
    renderCard({ roadmap: makeRoadmap({ description: '' }) });
    expect(screen.queryByText('Reactの基礎から応用まで')).not.toBeInTheDocument();
  });

  it('公開でない場合は公開バッジが非表示', () => {
    renderCard();
    expect(screen.queryByText('公開')).not.toBeInTheDocument();
  });
});
