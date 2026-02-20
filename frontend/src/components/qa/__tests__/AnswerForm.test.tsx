import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AnswerForm from '../AnswerForm';

describe('AnswerForm', () => {
  const onSubmit = vi.fn().mockResolvedValue(true);
  const onCancel = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('テキストエリアが表示される', () => {
    render(<AnswerForm onSubmit={onSubmit} />);
    expect(screen.getByRole('textbox')).toBeInTheDocument();
  });

  it('テキスト入力が可能', () => {
    render(<AnswerForm onSubmit={onSubmit} />);
    const textarea = screen.getByRole('textbox');
    fireEvent.change(textarea, { target: { value: 'テスト回答' } });
    expect(textarea).toHaveValue('テスト回答');
  });

  it('空入力時は送信ボタンがdisabled', () => {
    render(<AnswerForm onSubmit={onSubmit} />);
    expect(screen.getByText('回答を投稿')).toBeDisabled();
  });

  it('入力時は送信ボタンがenabled', () => {
    render(<AnswerForm onSubmit={onSubmit} />);
    fireEvent.change(screen.getByRole('textbox'), { target: { value: '回答内容' } });
    expect(screen.getByText('回答を投稿')).not.toBeDisabled();
  });

  it('フォーム送信でonSubmitが呼ばれる', async () => {
    render(<AnswerForm onSubmit={onSubmit} />);
    fireEvent.change(screen.getByRole('textbox'), { target: { value: '回答内容' } });
    fireEvent.submit(screen.getByRole('textbox').closest('form')!);
    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith('回答内容');
    });
  });

  it('送信成功後にテキストエリアがクリアされる', async () => {
    render(<AnswerForm onSubmit={onSubmit} />);
    const textarea = screen.getByRole('textbox');
    fireEvent.change(textarea, { target: { value: '送信テスト' } });
    fireEvent.submit(textarea.closest('form')!);
    await waitFor(() => {
      expect(textarea).toHaveValue('');
    });
  });

  it('Ctrl+Enterで送信される', async () => {
    render(<AnswerForm onSubmit={onSubmit} />);
    const textarea = screen.getByRole('textbox');
    fireEvent.change(textarea, { target: { value: 'ショートカット送信' } });
    fireEvent.keyDown(textarea, { key: 'Enter', ctrlKey: true });
    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith('ショートカット送信');
    });
  });

  it('Ctrl+Enterで空入力時は送信されない', () => {
    render(<AnswerForm onSubmit={onSubmit} />);
    fireEvent.keyDown(screen.getByRole('textbox'), { key: 'Enter', ctrlKey: true });
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('onCancelが渡された場合キャンセルボタンが表示される', () => {
    render(<AnswerForm onSubmit={onSubmit} onCancel={onCancel} />);
    const cancelBtn = screen.getByText('キャンセル');
    expect(cancelBtn).toBeInTheDocument();
    fireEvent.click(cancelBtn);
    expect(onCancel).toHaveBeenCalledOnce();
  });

  it('isEditモードでは更新ボタンテキストが表示される', () => {
    render(<AnswerForm onSubmit={onSubmit} isEdit initialBody="既存回答" />);
    expect(screen.getByText('回答を更新')).toBeInTheDocument();
    expect(screen.getByRole('textbox')).toHaveValue('既存回答');
  });

  it('ショートカットヒントが表示される', () => {
    render(<AnswerForm onSubmit={onSubmit} />);
    expect(screen.getByText(/Ctrl\+Enter/)).toBeInTheDocument();
  });
});
