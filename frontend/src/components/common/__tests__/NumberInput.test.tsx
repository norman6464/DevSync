import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import NumberInput from '../NumberInput';

describe('NumberInput', () => {
  it('数値入力が表示される', () => {
    render(<NumberInput value={5} onChange={() => {}} />);
    expect(screen.getByRole('spinbutton')).toBeInTheDocument();
  });

  it('値が表示される', () => {
    render(<NumberInput value={42} onChange={() => {}} />);
    expect(screen.getByDisplayValue('42')).toBeInTheDocument();
  });

  it('増加ボタンが表示される', () => {
    const { container } = render(<NumberInput value={5} onChange={() => {}} />);
    expect(container.querySelector('.lucide-plus')).toBeInTheDocument();
  });

  it('減少ボタンが表示される', () => {
    const { container } = render(<NumberInput value={5} onChange={() => {}} />);
    expect(container.querySelector('.lucide-minus')).toBeInTheDocument();
  });

  it('増加ボタンで値が増える', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<NumberInput value={5} onChange={onChange} />);
    await user.click(screen.getByLabelText('増加'));
    expect(onChange).toHaveBeenCalledWith(6);
  });

  it('減少ボタンで値が減る', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<NumberInput value={5} onChange={onChange} />);
    await user.click(screen.getByLabelText('減少'));
    expect(onChange).toHaveBeenCalledWith(4);
  });

  it('最小値以下にはならない', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<NumberInput value={0} onChange={onChange} min={0} />);
    await user.click(screen.getByLabelText('減少'));
    expect(onChange).not.toHaveBeenCalled();
  });

  it('最大値以上にはならない', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<NumberInput value={10} onChange={onChange} max={10} />);
    await user.click(screen.getByLabelText('増加'));
    expect(onChange).not.toHaveBeenCalled();
  });

  it('ステップが設定される', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<NumberInput value={10} onChange={onChange} step={5} />);
    await user.click(screen.getByLabelText('増加'));
    expect(onChange).toHaveBeenCalledWith(15);
  });

  it('無効状態が適用される', () => {
    render(<NumberInput value={5} onChange={() => {}} disabled />);
    expect(screen.getByRole('spinbutton')).toBeDisabled();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<NumberInput value={5} onChange={() => {}} className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
