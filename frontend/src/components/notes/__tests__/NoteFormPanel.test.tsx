import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import NoteFormPanel from '../NoteFormPanel';

const defaultProps = {
  editingNote: false,
  title: '',
  setTitle: vi.fn(),
  content: '',
  setContent: vi.fn(),
  tags: '',
  setTags: vi.fn(),
  saving: false,
  onSubmit: vi.fn(),
  onCancel: vi.fn(),
};

describe('NoteFormPanel', () => {
  it('新規作成タイトルを表示する', () => {
    render(<NoteFormPanel {...defaultProps} />);
    expect(screen.getByText('新規ノート')).toBeInTheDocument();
  });

  it('編集タイトルを表示する', () => {
    render(<NoteFormPanel {...defaultProps} editingNote={true} />);
    expect(screen.getByText('ノート編集')).toBeInTheDocument();
  });

  it('タイトル・内容・タグのフィールドを表示する', () => {
    render(<NoteFormPanel {...defaultProps} />);
    expect(screen.getByLabelText('タイトル')).toBeInTheDocument();
    expect(screen.getByLabelText('内容（マークダウン対応）')).toBeInTheDocument();
    expect(screen.getByLabelText('タグ（カンマ区切り）')).toBeInTheDocument();
  });

  it('タイトル入力でsetTitleが呼ばれる', () => {
    const setTitle = vi.fn();
    render(<NoteFormPanel {...defaultProps} setTitle={setTitle} />);
    fireEvent.change(screen.getByLabelText('タイトル'), { target: { value: 'テスト' } });
    expect(setTitle).toHaveBeenCalledWith('テスト');
  });

  it('内容入力でsetContentが呼ばれる', () => {
    const setContent = vi.fn();
    render(<NoteFormPanel {...defaultProps} setContent={setContent} />);
    fireEvent.change(screen.getByLabelText('内容（マークダウン対応）'), { target: { value: 'テスト内容' } });
    expect(setContent).toHaveBeenCalledWith('テスト内容');
  });

  it('タグ入力でsetTagsが呼ばれる', () => {
    const setTags = vi.fn();
    render(<NoteFormPanel {...defaultProps} setTags={setTags} />);
    fireEvent.change(screen.getByLabelText('タグ（カンマ区切り）'), { target: { value: 'React, Go' } });
    expect(setTags).toHaveBeenCalledWith('React, Go');
  });

  it('文字数カウンターを表示する', () => {
    render(<NoteFormPanel {...defaultProps} content="Hello" />);
    expect(screen.getByText('5/10000')).toBeInTheDocument();
  });

  it('キャンセルボタンでonCancelが呼ばれる', () => {
    const onCancel = vi.fn();
    render(<NoteFormPanel {...defaultProps} onCancel={onCancel} />);
    fireEvent.click(screen.getByText('キャンセル'));
    expect(onCancel).toHaveBeenCalled();
  });

  it('保存中は保存ボタンが無効になる', () => {
    render(<NoteFormPanel {...defaultProps} saving={true} />);
    const saveBtn = screen.getByText('保存中...');
    expect(saveBtn.closest('button')).toBeDisabled();
  });

  it('保存ボタンのテキストが通常時は「保存」', () => {
    render(<NoteFormPanel {...defaultProps} />);
    expect(screen.getByText('保存')).toBeInTheDocument();
  });

  it('タイトルにmaxLength=200が設定されている', () => {
    render(<NoteFormPanel {...defaultProps} />);
    expect(screen.getByLabelText('タイトル')).toHaveAttribute('maxLength', '200');
  });

  it('内容にmaxLength=10000が設定されている', () => {
    render(<NoteFormPanel {...defaultProps} />);
    expect(screen.getByLabelText('内容（マークダウン対応）')).toHaveAttribute('maxLength', '10000');
  });
});
