import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import OTPInput from '../OTPInput';

describe('OTPInput', () => {
  it('指定桁数の入力フィールドが表示される', () => {
    render(<OTPInput length={6} value="" onChange={() => {}} />);
    const inputs = screen.getAllByRole('textbox');
    expect(inputs.length).toBe(6);
  });

  it('4桁の入力フィールドが表示される', () => {
    render(<OTPInput length={4} value="" onChange={() => {}} />);
    const inputs = screen.getAllByRole('textbox');
    expect(inputs.length).toBe(4);
  });

  it('値が各フィールドに表示される', () => {
    render(<OTPInput length={4} value="1234" onChange={() => {}} />);
    const inputs = screen.getAllByRole('textbox');
    expect(inputs[0]).toHaveValue('1');
    expect(inputs[1]).toHaveValue('2');
    expect(inputs[2]).toHaveValue('3');
    expect(inputs[3]).toHaveValue('4');
  });

  it('入力でonChangeが呼ばれる', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<OTPInput length={4} value="" onChange={onChange} />);
    const inputs = screen.getAllByRole('textbox');
    await user.type(inputs[0], '5');
    expect(onChange).toHaveBeenCalled();
  });

  it('ラベルが表示される', () => {
    render(<OTPInput length={4} value="" onChange={() => {}} label="認証コード" />);
    expect(screen.getByText('認証コード')).toBeInTheDocument();
  });

  it('エラーメッセージが表示される', () => {
    render(<OTPInput length={4} value="" onChange={() => {}} error="無効なコードです" />);
    expect(screen.getByText('無効なコードです')).toBeInTheDocument();
  });

  it('無効状態で入力が無効になる', () => {
    render(<OTPInput length={4} value="" onChange={() => {}} disabled />);
    const inputs = screen.getAllByRole('textbox');
    inputs.forEach((input) => expect(input).toBeDisabled());
  });

  it('エラー時にボーダーが赤くなる', () => {
    render(<OTPInput length={4} value="" onChange={() => {}} error="エラー" />);
    const inputs = screen.getAllByRole('textbox');
    expect(inputs[0].className).toContain('border-red-500');
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<OTPInput length={4} value="" onChange={() => {}} className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('数字のみ入力可能', () => {
    render(<OTPInput length={4} value="" onChange={() => {}} />);
    const inputs = screen.getAllByRole('textbox');
    expect(inputs[0]).toHaveAttribute('inputMode', 'numeric');
  });
});
