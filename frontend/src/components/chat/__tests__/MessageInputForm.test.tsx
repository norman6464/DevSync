import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import MessageInputForm from '../MessageInputForm';

const defaultProps = {
  value: '',
  onChange: vi.fn(),
  onSubmit: vi.fn(),
  placeholder: 'メッセージを入力...',
};

describe('MessageInputForm', () => {
  it('入力フィールドを表示する', () => {
    render(<MessageInputForm {...defaultProps} />);
    expect(screen.getByPlaceholderText('メッセージを入力...')).toBeInTheDocument();
  });

  it('送信ボタンを表示する', () => {
    render(<MessageInputForm {...defaultProps} />);
    expect(screen.getByText('送信')).toBeInTheDocument();
  });

  it('値が空の場合は送信ボタンが無効', () => {
    render(<MessageInputForm {...defaultProps} value="" />);
    expect(screen.getByText('送信').closest('button')).toBeDisabled();
  });

  it('値がある場合は送信ボタンが有効', () => {
    render(<MessageInputForm {...defaultProps} value="テスト" />);
    expect(screen.getByText('送信').closest('button')).not.toBeDisabled();
  });

  it('空白のみの場合は送信ボタンが無効', () => {
    render(<MessageInputForm {...defaultProps} value="   " />);
    expect(screen.getByText('送信').closest('button')).toBeDisabled();
  });

  it('入力変更でonChangeが呼ばれる', () => {
    const onChange = vi.fn();
    render(<MessageInputForm {...defaultProps} onChange={onChange} />);
    fireEvent.change(screen.getByPlaceholderText('メッセージを入力...'), {
      target: { value: 'Hello' },
    });
    expect(onChange).toHaveBeenCalledWith('Hello');
  });

  it('フォーム送信でonSubmitが呼ばれる', () => {
    const onSubmit = vi.fn((e) => e.preventDefault());
    render(<MessageInputForm {...defaultProps} value="テスト" onSubmit={onSubmit} />);
    fireEvent.submit(screen.getByText('送信').closest('form')!);
    expect(onSubmit).toHaveBeenCalled();
  });

  it('プレースホルダーが設定される', () => {
    render(<MessageInputForm {...defaultProps} placeholder="カスタムプレースホルダー" />);
    expect(screen.getByPlaceholderText('カスタムプレースホルダー')).toBeInTheDocument();
  });
});
