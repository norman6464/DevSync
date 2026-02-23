import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import PasswordInput from '../PasswordInput';

describe('PasswordInput', () => {
  it('パスワード入力が表示される', () => {
    render(<PasswordInput value="" onChange={() => {}} />);
    const input = screen.getByPlaceholderText('パスワード');
    expect(input).toHaveAttribute('type', 'password');
  });

  it('表示切り替えボタンが表示される', () => {
    const { container } = render(<PasswordInput value="" onChange={() => {}} />);
    expect(container.querySelector('.lucide-eye-off')).toBeInTheDocument();
  });

  it('クリックでパスワードが表示される', async () => {
    const user = userEvent.setup();
    render(<PasswordInput value="secret" onChange={() => {}} />);
    await user.click(screen.getByRole('button'));
    expect(screen.getByPlaceholderText('パスワード')).toHaveAttribute('type', 'text');
  });

  it('再クリックで非表示に戻る', async () => {
    const user = userEvent.setup();
    render(<PasswordInput value="secret" onChange={() => {}} />);
    await user.click(screen.getByRole('button'));
    await user.click(screen.getByRole('button'));
    expect(screen.getByPlaceholderText('パスワード')).toHaveAttribute('type', 'password');
  });

  it('値の変更でコールバックが呼ばれる', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<PasswordInput value="" onChange={onChange} />);
    await user.type(screen.getByPlaceholderText('パスワード'), 'a');
    expect(onChange).toHaveBeenCalled();
  });

  it('強度メーターが表示される', () => {
    const { container } = render(<PasswordInput value="Test1234!" onChange={() => {}} showStrength />);
    const bars = container.querySelectorAll('[data-testid="strength-bar"]');
    expect(bars.length).toBe(4);
  });

  it('弱いパスワードは赤', () => {
    const { container } = render(<PasswordInput value="ab" onChange={() => {}} showStrength />);
    expect(container.querySelector('.bg-red-500')).toBeInTheDocument();
  });

  it('ラベルが表示される', () => {
    render(<PasswordInput value="" onChange={() => {}} label="新しいパスワード" />);
    expect(screen.getByText('新しいパスワード')).toBeInTheDocument();
  });

  it('エラーメッセージが表示される', () => {
    render(<PasswordInput value="" onChange={() => {}} error="必須項目です" />);
    expect(screen.getByText('必須項目です')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<PasswordInput value="" onChange={() => {}} className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('無効状態が適用される', () => {
    render(<PasswordInput value="" onChange={() => {}} disabled />);
    expect(screen.getByPlaceholderText('パスワード')).toBeDisabled();
  });
});
