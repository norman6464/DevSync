import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import GoalFormModal from '../GoalFormModal';

const defaultProps = {
  isOpen: true,
  isEditing: false,
  saving: false,
  title: '',
  setTitle: vi.fn(),
  description: '',
  setDescription: vi.fn(),
  category: 'programming' as const,
  setCategory: vi.fn(),
  targetDate: '',
  setTargetDate: vi.fn(),
  onSubmit: vi.fn(),
  onCancel: vi.fn(),
};

describe('GoalFormModal', () => {
  it('新規作成モードでタイトルを表示する', () => {
    render(<GoalFormModal {...defaultProps} />);
    expect(screen.getByText('目標を追加')).toBeInTheDocument();
  });

  it('編集モードでタイトルを表示する', () => {
    render(<GoalFormModal {...defaultProps} isEditing={true} />);
    expect(screen.getByText('目標を編集')).toBeInTheDocument();
  });

  it('タイトル入力フィールドを表示する', () => {
    render(<GoalFormModal {...defaultProps} />);
    expect(screen.getByLabelText('タイトル')).toBeInTheDocument();
  });

  it('説明入力フィールドを表示する', () => {
    render(<GoalFormModal {...defaultProps} />);
    expect(screen.getByLabelText('学習の進捗を追跡し、目標を設定しましょう')).toBeInTheDocument();
  });

  it('カテゴリ選択を表示する', () => {
    render(<GoalFormModal {...defaultProps} />);
    expect(screen.getByLabelText('カテゴリ')).toBeInTheDocument();
  });

  it('タイトルが空の場合は送信ボタンが無効', () => {
    render(<GoalFormModal {...defaultProps} title="" />);
    expect(screen.getByText('作成')).toBeDisabled();
  });

  it('タイトルがある場合は送信ボタンが有効', () => {
    render(<GoalFormModal {...defaultProps} title="テスト目標" />);
    expect(screen.getByText('作成')).not.toBeDisabled();
  });

  it('saving中はボタンテキストが変わる', () => {
    render(<GoalFormModal {...defaultProps} saving={true} title="テスト" />);
    expect(screen.getByText('読み込み中...')).toBeInTheDocument();
  });

  it('編集モードでは保存ボタンを表示する', () => {
    render(<GoalFormModal {...defaultProps} isEditing={true} title="テスト" />);
    expect(screen.getByText('保存')).toBeInTheDocument();
  });

  it('キャンセルボタンを表示する', () => {
    render(<GoalFormModal {...defaultProps} />);
    expect(screen.getByText('キャンセル')).toBeInTheDocument();
  });

  it('isOpenがfalseの場合はモーダルを表示しない', () => {
    render(<GoalFormModal {...defaultProps} isOpen={false} />);
    expect(screen.queryByText('目標を追加')).not.toBeInTheDocument();
  });

  it('説明の文字数カウンターを表示する', () => {
    render(<GoalFormModal {...defaultProps} description="テスト説明" />);
    expect(screen.getByText('5/2000')).toBeInTheDocument();
  });
});
