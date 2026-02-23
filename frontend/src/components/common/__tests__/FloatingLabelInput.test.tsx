import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import FloatingLabelInput from '../FloatingLabelInput';

describe('FloatingLabelInput', () => {
  it('入力フィールドが表示される', () => {
    render(<FloatingLabelInput value="" onChange={() => {}} label="メール" />);
    expect(screen.getByRole('textbox')).toBeInTheDocument();
  });

  it('ラベルが表示される', () => {
    render(<FloatingLabelInput value="" onChange={() => {}} label="メール" />);
    expect(screen.getByText('メール')).toBeInTheDocument();
  });

  it('値がある時ラベルが浮く', () => {
    render(<FloatingLabelInput value="test@example.com" onChange={() => {}} label="メール" />);
    const label = screen.getByText('メール');
    expect(label.className).toContain('text-xs');
  });

  it('値がない時ラベルが通常位置', () => {
    render(<FloatingLabelInput value="" onChange={() => {}} label="メール" />);
    const label = screen.getByText('メール');
    expect(label.className).toContain('text-sm');
  });

  it('値の変更でコールバックが呼ばれる', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<FloatingLabelInput value="" onChange={onChange} label="メール" />);
    await user.type(screen.getByRole('textbox'), 'a');
    expect(onChange).toHaveBeenCalled();
  });

  it('エラーメッセージが表示される', () => {
    render(<FloatingLabelInput value="" onChange={() => {}} label="メール" error="必須です" />);
    expect(screen.getByText('必須です')).toBeInTheDocument();
  });

  it('エラー時にボーダーが赤くなる', () => {
    render(<FloatingLabelInput value="" onChange={() => {}} label="メール" error="エラー" />);
    expect(screen.getByRole('textbox').className).toContain('border-red-500');
  });

  it('無効状態が適用される', () => {
    render(<FloatingLabelInput value="" onChange={() => {}} label="メール" disabled />);
    expect(screen.getByRole('textbox')).toBeDisabled();
  });

  it('type属性が設定される', () => {
    render(<FloatingLabelInput value="" onChange={() => {}} label="メール" type="email" />);
    expect(screen.getByRole('textbox')).toHaveAttribute('type', 'email');
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<FloatingLabelInput value="" onChange={() => {}} label="メール" className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
