import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import LogFormModal from '../LogFormModal';
import type { LearningLog } from '../../../types/learningLog';

const defaultProps = {
  isOpen: true,
  editingLog: null as LearningLog | null,
  title: '',
  setTitle: vi.fn(),
  content: '',
  setContent: vi.fn(),
  category: 'coding' as const,
  setCategory: vi.fn(),
  duration: '',
  setDuration: vi.fn(),
  saving: false,
  onSubmit: vi.fn(),
  onClose: vi.fn(),
};

const makeLearningLog = (): LearningLog => ({
  id: 1,
  user_id: 1,
  title: 'テストログ',
  content: 'テスト内容',
  category: 'coding',
  duration: 60,
  source: 'manual',
  is_favorite: false,
  created_at: '2025-01-01',
  updated_at: '2025-01-01',
});

describe('LogFormModal', () => {
  it('新規作成時のタイトルを表示する', () => {
    render(<LogFormModal {...defaultProps} />);
    const elements = screen.getAllByText('ログを追加');
    expect(elements.length).toBeGreaterThanOrEqual(1);
    expect(elements[0].tagName).toBe('H2');
  });

  it('編集時のタイトルを表示する', () => {
    render(<LogFormModal {...defaultProps} editingLog={makeLearningLog()} />);
    expect(screen.getByText('ログを編集')).toBeInTheDocument();
  });

  it('タイトル入力フィールドを表示する', () => {
    render(<LogFormModal {...defaultProps} />);
    expect(screen.getByLabelText('タイトル')).toBeInTheDocument();
  });

  it('内容テキストエリアを表示する', () => {
    render(<LogFormModal {...defaultProps} />);
    expect(screen.getByLabelText('内容')).toBeInTheDocument();
  });

  it('カテゴリ選択を表示する', () => {
    render(<LogFormModal {...defaultProps} />);
    expect(screen.getByLabelText('カテゴリ')).toBeInTheDocument();
  });

  it('5つのカテゴリ選択肢を表示する', () => {
    render(<LogFormModal {...defaultProps} />);
    const select = screen.getByLabelText('カテゴリ');
    const options = select.querySelectorAll('option');
    expect(options).toHaveLength(5);
  });

  it('学習時間入力フィールドを表示する', () => {
    render(<LogFormModal {...defaultProps} />);
    expect(screen.getByLabelText('学習時間')).toBeInTheDocument();
  });

  it('キャンセルボタンクリック時にonCloseが呼ばれる', () => {
    const onClose = vi.fn();
    render(<LogFormModal {...defaultProps} onClose={onClose} />);
    fireEvent.click(screen.getByText('キャンセル'));
    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('タイトルとコンテンツが空の場合に送信ボタンが無効になる', () => {
    render(<LogFormModal {...defaultProps} title="" content="" />);
    const buttons = screen.getAllByRole('button');
    const submitButton = buttons.find(b => b.getAttribute('type') === 'submit');
    expect(submitButton).toBeDisabled();
  });

  it('タイトルとコンテンツがある場合に送信ボタンが有効になる', () => {
    render(<LogFormModal {...defaultProps} title="テスト" content="内容あり" />);
    const buttons = screen.getAllByRole('button');
    const submitButton = buttons.find(b => b.getAttribute('type') === 'submit');
    expect(submitButton).not.toBeDisabled();
  });

  it('saving中は送信ボタンが無効になる', () => {
    render(<LogFormModal {...defaultProps} title="テスト" content="内容あり" saving={true} />);
    const buttons = screen.getAllByRole('button');
    const submitButton = buttons.find(b => b.getAttribute('type') === 'submit');
    expect(submitButton).toBeDisabled();
  });

  it('内容の文字数カウンターを表示する', () => {
    render(<LogFormModal {...defaultProps} content="Hello" />);
    expect(screen.getByText('5/5000')).toBeInTheDocument();
  });

  it('isOpen=falseの場合モーダルが表示されない', () => {
    render(<LogFormModal {...defaultProps} isOpen={false} />);
    expect(screen.queryByText('ログを追加')).not.toBeInTheDocument();
  });
});
