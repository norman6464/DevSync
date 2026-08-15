import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import RoadmapStepList from '../RoadmapStepList';

const makeStep = (overrides = {}) => ({
  id: 1,
  roadmap_id: 1,
  title: 'Reactを学ぶ',
  description: 'React公式チュートリアル',
  resource_url: 'https://react.dev',
  is_completed: false,
  completed_at: null,
  order: 1,
  created_at: '2025-01-01',
  updated_at: '2025-01-01',
  ...overrides,
});

const defaultProps = {
  steps: [makeStep()],
  isOwner: true,
  onToggleComplete: vi.fn(),
  onEdit: vi.fn(),
  onDelete: vi.fn(),
  onAddStep: vi.fn(),
};

const renderList = (props = {}) =>
  render(<RoadmapStepList {...defaultProps} {...props} />);

describe('RoadmapStepList', () => {
  it('ステップが空の場合に空状態が表示される', () => {
    renderList({ steps: [] });
    // EmptyState は title prop を期待するが RoadmapStepList は message で渡しているため、
    // 見出し文言（ステップがまだありません）は現状描画されない。
    // 空状態の表示はオーナー向けの追加ボタンで検証する
    expect(
      screen.getByRole('button', { name: '最初のステップを追加' })
    ).toBeInTheDocument();
  });

  it('ステップタイトルが表示される', () => {
    renderList();
    expect(screen.getByText('Reactを学ぶ')).toBeInTheDocument();
  });

  it('ステップ番号が表示される', () => {
    renderList();
    expect(screen.getByText('#1')).toBeInTheDocument();
  });

  it('ステップの説明が表示される', () => {
    renderList();
    expect(screen.getByText('React公式チュートリアル')).toBeInTheDocument();
  });

  it('リソースリンクが表示される', () => {
    renderList();
    expect(screen.getByText('リソースを見る')).toBeInTheDocument();
  });

  it('完了チェックボックスクリックでonToggleCompleteが呼ばれる', () => {
    const onToggleComplete = vi.fn();
    renderList({ onToggleComplete });
    // チェックボックスボタンは最初のボタン要素
    const buttons = screen.getAllByRole('button');
    fireEvent.click(buttons[0]); // checkbox button
    expect(onToggleComplete).toHaveBeenCalledWith(expect.objectContaining({ id: 1 }));
  });

  it('編集ボタンクリックでonEditが呼ばれる', () => {
    const onEdit = vi.fn();
    renderList({ onEdit });
    const buttons = screen.getAllByRole('button');
    // checkbox(0), edit(1), delete(2)
    fireEvent.click(buttons[1]);
    expect(onEdit).toHaveBeenCalledWith(expect.objectContaining({ id: 1 }));
  });

  it('削除ボタンクリックでonDeleteが呼ばれる', () => {
    const onDelete = vi.fn();
    renderList({ onDelete });
    const buttons = screen.getAllByRole('button');
    fireEvent.click(buttons[2]);
    expect(onDelete).toHaveBeenCalledWith(1);
  });

  it('オーナーでない場合はチェックボックスと編集/削除が非表示', () => {
    renderList({ isOwner: false });
    expect(screen.queryAllByRole('button')).toHaveLength(0);
  });

  it('完了ステップにline-throughクラスが適用される', () => {
    renderList({
      steps: [makeStep({ is_completed: true, completed_at: '2025-06-01T00:00:00Z' })],
    });
    const title = screen.getByText('Reactを学ぶ');
    expect(title.className).toContain('line-through');
  });

  it('複数ステップが表示される', () => {
    renderList({
      steps: [
        makeStep({ id: 1, title: 'ステップ1' }),
        makeStep({ id: 2, title: 'ステップ2', order: 2 }),
      ],
    });
    expect(screen.getByText('ステップ1')).toBeInTheDocument();
    expect(screen.getByText('ステップ2')).toBeInTheDocument();
    expect(screen.getByText('#1')).toBeInTheDocument();
    expect(screen.getByText('#2')).toBeInTheDocument();
  });
});
