import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import GoalCard from '../GoalCard';
import type { LearningGoal } from '../../../api/goals';

const baseGoal: LearningGoal = {
  id: 1,
  user_id: 1,
  title: 'TypeScriptを学ぶ',
  description: 'TypeScriptの基礎をマスターする',
  category: 'language',
  target_date: null,
  progress: 50,
  status: 'active',
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
  completed_at: null,
};

const defaultProps = {
  goal: baseGoal,
  onEdit: vi.fn(),
  onDelete: vi.fn(),
  onDuplicate: vi.fn(),
  onProgressChange: vi.fn(),
  onStatusChange: vi.fn(),
};

describe('GoalCard', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('タイトルと説明が表示される', () => {
    render(<GoalCard {...defaultProps} />);
    expect(screen.getByText('TypeScriptを学ぶ')).toBeInTheDocument();
    expect(screen.getByText('TypeScriptの基礎をマスターする')).toBeInTheDocument();
  });

  it('進捗が表示される', () => {
    render(<GoalCard {...defaultProps} />);
    expect(screen.getByText('50%')).toBeInTheDocument();
    expect(screen.getByText('進捗')).toBeInTheDocument();
  });

  it('ステータスバッジが表示される（active）', () => {
    render(<GoalCard {...defaultProps} />);
    expect(screen.getByText('進行中')).toBeInTheDocument();
  });

  it('ステータスバッジが表示される（completed）', () => {
    const completedGoal = { ...baseGoal, status: 'completed' as const, progress: 100 };
    render(<GoalCard {...defaultProps} goal={completedGoal} />);
    expect(screen.getByText('完了')).toBeInTheDocument();
  });

  it('ステータスバッジが表示される（paused）', () => {
    const pausedGoal = { ...baseGoal, status: 'paused' as const };
    render(<GoalCard {...defaultProps} goal={pausedGoal} />);
    expect(screen.getByText('一時停止')).toBeInTheDocument();
  });

  it('カテゴリラベルが表示される', () => {
    render(<GoalCard {...defaultProps} />);
    expect(screen.getByText('言語')).toBeInTheDocument();
  });

  it('activeの場合に一時停止ボタンが表示される', () => {
    render(<GoalCard {...defaultProps} />);
    expect(screen.getByLabelText('一時停止')).toBeInTheDocument();
  });

  it('pausedの場合に再開ボタンが表示される', () => {
    const pausedGoal = { ...baseGoal, status: 'paused' as const };
    render(<GoalCard {...defaultProps} goal={pausedGoal} />);
    expect(screen.getByLabelText('再開')).toBeInTheDocument();
  });

  it('一時停止ボタンクリックでonStatusChangeが呼ばれる', () => {
    render(<GoalCard {...defaultProps} />);
    fireEvent.click(screen.getByLabelText('一時停止'));
    expect(defaultProps.onStatusChange).toHaveBeenCalledWith(baseGoal, 'paused');
  });

  it('編集ボタンクリックでonEditが呼ばれる', () => {
    render(<GoalCard {...defaultProps} />);
    fireEvent.click(screen.getByLabelText('編集'));
    expect(defaultProps.onEdit).toHaveBeenCalledWith(baseGoal);
  });

  it('削除ボタンクリックでonDeleteが呼ばれる', () => {
    render(<GoalCard {...defaultProps} />);
    fireEvent.click(screen.getByLabelText('削除'));
    expect(defaultProps.onDelete).toHaveBeenCalledWith(1);
  });

  it('複製ボタンクリックでonDuplicateが呼ばれる', () => {
    render(<GoalCard {...defaultProps} />);
    fireEvent.click(screen.getByLabelText('複製'));
    expect(defaultProps.onDuplicate).toHaveBeenCalledWith(1);
  });

  it('activeの場合に進捗スライダーが表示される', () => {
    render(<GoalCard {...defaultProps} />);
    const slider = screen.getByRole('slider');
    expect(slider).toBeInTheDocument();
    expect(slider).toHaveValue('50');
  });

  it('completedの場合に進捗スライダーが非表示', () => {
    const completedGoal = { ...baseGoal, status: 'completed' as const, progress: 100 };
    render(<GoalCard {...defaultProps} goal={completedGoal} />);
    expect(screen.queryByRole('slider')).not.toBeInTheDocument();
  });

  it('目標日が表示される', () => {
    const goalWithDate = { ...baseGoal, target_date: '2026-06-01T00:00:00Z' };
    render(<GoalCard {...defaultProps} goal={goalWithDate} />);
    expect(screen.getByText(/目標日/)).toBeInTheDocument();
  });
});
