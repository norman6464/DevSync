import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import SliderInput from '../SliderInput';

describe('SliderInput', () => {
  it('スライダーが表示される', () => {
    render(<SliderInput value={50} onChange={() => {}} />);
    expect(screen.getByRole('slider')).toBeInTheDocument();
  });

  it('値が表示される', () => {
    render(<SliderInput value={50} onChange={() => {}} showValue />);
    expect(screen.getByText('50')).toBeInTheDocument();
  });

  it('値の変更でコールバックが呼ばれる', () => {
    const onChange = vi.fn();
    render(<SliderInput value={50} onChange={onChange} />);
    fireEvent.change(screen.getByRole('slider'), { target: { value: '75' } });
    expect(onChange).toHaveBeenCalledWith(75);
  });

  it('最小値が設定される', () => {
    render(<SliderInput value={50} onChange={() => {}} min={10} />);
    expect(screen.getByRole('slider')).toHaveAttribute('min', '10');
  });

  it('最大値が設定される', () => {
    render(<SliderInput value={50} onChange={() => {}} max={200} />);
    expect(screen.getByRole('slider')).toHaveAttribute('max', '200');
  });

  it('ステップが設定される', () => {
    render(<SliderInput value={50} onChange={() => {}} step={5} />);
    expect(screen.getByRole('slider')).toHaveAttribute('step', '5');
  });

  it('ラベルが表示される', () => {
    render(<SliderInput value={50} onChange={() => {}} label="音量" />);
    expect(screen.getByText('音量')).toBeInTheDocument();
  });

  it('無効状態が適用される', () => {
    render(<SliderInput value={50} onChange={() => {}} disabled />);
    expect(screen.getByRole('slider')).toBeDisabled();
  });

  it('サフィックスが表示される', () => {
    render(<SliderInput value={50} onChange={() => {}} showValue suffix="%" />);
    expect(screen.getByText('%')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<SliderInput value={50} onChange={() => {}} className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
