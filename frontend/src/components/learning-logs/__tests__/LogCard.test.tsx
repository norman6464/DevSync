import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import LogCard from '../LogCard';
import type { LearningLog } from '../../../types/learningLog';

const baseLog: LearningLog = {
  id: 1,
  user_id: 1,
  title: 'TypeScript学習',
  content: 'ジェネリクスについて学んだ',
  category: 'coding',
  duration: 45,
  source: 'manual',
  is_favorite: false,
  created_at: '2026-02-19T10:00:00Z',
  updated_at: '2026-02-19T10:00:00Z',
};

const defaultProps = {
  log: baseLog,
  onEdit: vi.fn(),
  onDelete: vi.fn(),
  onToggleFavorite: vi.fn(),
};

describe('LogCard', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('タイトルが表示される', () => {
    render(<LogCard {...defaultProps} />);
    expect(screen.getByText('TypeScript学習')).toBeInTheDocument();
  });

  it('内容が表示される', () => {
    render(<LogCard {...defaultProps} />);
    expect(screen.getByText('ジェネリクスについて学んだ')).toBeInTheDocument();
  });

  it('カテゴリラベルが表示される', () => {
    render(<LogCard {...defaultProps} />);
    expect(screen.getByText('コーディング')).toBeInTheDocument();
  });

  it('学習時間が表示される', () => {
    render(<LogCard {...defaultProps} />);
    expect(screen.getByText('45分')).toBeInTheDocument();
  });

  it('学習時間0の場合は時間が表示されない', () => {
    const noTimeLog = { ...baseLog, duration: 0 };
    render(<LogCard {...defaultProps} log={noTimeLog} />);
    expect(screen.queryByText(/分$/)).not.toBeInTheDocument();
  });

  it('readingカテゴリが正しく表示される', () => {
    const readingLog = { ...baseLog, category: 'reading' as const };
    render(<LogCard {...defaultProps} log={readingLog} />);
    expect(screen.getByText('読書')).toBeInTheDocument();
  });

  it('お気に入りボタンクリックでonToggleFavoriteが呼ばれる', () => {
    render(<LogCard {...defaultProps} />);
    fireEvent.click(screen.getByTitle('お気に入り切替'));
    expect(defaultProps.onToggleFavorite).toHaveBeenCalledWith(1);
  });

  it('編集ボタンクリックでonEditが呼ばれる', () => {
    render(<LogCard {...defaultProps} />);
    fireEvent.click(screen.getByTitle('編集'));
    expect(defaultProps.onEdit).toHaveBeenCalledWith(baseLog);
  });

  it('削除ボタンクリックでonDeleteが呼ばれる', () => {
    render(<LogCard {...defaultProps} />);
    fireEvent.click(screen.getByTitle('削除'));
    expect(defaultProps.onDelete).toHaveBeenCalledWith(1);
  });

  it('日付が表示される', () => {
    render(<LogCard {...defaultProps} />);
    const dateEl = screen.getByText(/2026/);
    expect(dateEl).toBeInTheDocument();
  });
});
