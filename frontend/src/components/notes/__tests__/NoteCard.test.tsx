import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import NoteCard from '../NoteCard';
import type { Note } from '../../../api/notes';

const makeNote = (overrides: Partial<Note> = {}): Note => ({
  id: 1,
  user_id: 10,
  title: 'テストノート',
  content: 'テスト内容です。これはサンプルテキストです。',
  tags: 'React,TypeScript',
  is_favorite: false,
  is_archived: false,
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-02-01T00:00:00Z',
  ...overrides,
});

const defaultProps = {
  note: makeNote(),
  onToggleFavorite: vi.fn(),
  onEdit: vi.fn(),
  onDelete: vi.fn(),
};

describe('NoteCard', () => {
  it('ノートタイトルを表示する', () => {
    render(<NoteCard {...defaultProps} />);
    expect(screen.getByText('テストノート')).toBeInTheDocument();
  });

  it('ノート内容を表示する', () => {
    render(<NoteCard {...defaultProps} />);
    expect(screen.getByText('テスト内容です。これはサンプルテキストです。')).toBeInTheDocument();
  });

  it('タグを表示する', () => {
    render(<NoteCard {...defaultProps} />);
    expect(screen.getByText('React')).toBeInTheDocument();
    expect(screen.getByText('TypeScript')).toBeInTheDocument();
  });

  it('タグが空の場合はタグセクションを表示しない', () => {
    render(<NoteCard {...defaultProps} note={makeNote({ tags: '' })} />);
    expect(screen.queryByText('React')).not.toBeInTheDocument();
  });

  it('文字数を表示する', () => {
    render(<NoteCard {...defaultProps} />);
    expect(screen.getByText(/22/)).toBeInTheDocument();
  });

  it('お気に入りボタンをクリックするとonToggleFavoriteが呼ばれる', () => {
    const onToggleFavorite = vi.fn();
    render(<NoteCard {...defaultProps} onToggleFavorite={onToggleFavorite} />);
    fireEvent.click(screen.getByLabelText('お気に入り'));
    expect(onToggleFavorite).toHaveBeenCalledWith(1);
  });

  it('編集ボタンをクリックするとonEditが呼ばれる', () => {
    const onEdit = vi.fn();
    const note = makeNote();
    render(<NoteCard {...defaultProps} note={note} onEdit={onEdit} />);
    fireEvent.click(screen.getByLabelText('編集'));
    expect(onEdit).toHaveBeenCalledWith(note);
  });

  it('削除ボタンをクリックするとonDeleteが呼ばれる', () => {
    const onDelete = vi.fn();
    render(<NoteCard {...defaultProps} onDelete={onDelete} />);
    fireEvent.click(screen.getByLabelText('削除'));
    expect(onDelete).toHaveBeenCalledWith(1);
  });

  it('お気に入りノートにはお気に入りアイコンが表示される', () => {
    render(<NoteCard {...defaultProps} note={makeNote({ is_favorite: true })} />);
    const title = screen.getByText('テストノート');
    const starIcon = title.parentElement?.querySelector('.fill-yellow-500');
    expect(starIcon).toBeInTheDocument();
  });
});
