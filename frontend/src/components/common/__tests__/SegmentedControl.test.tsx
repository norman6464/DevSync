import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import SegmentedControl from '../SegmentedControl';

const options = [
  { value: 'day', label: '日' },
  { value: 'week', label: '週' },
  { value: 'month', label: '月' },
];

describe('SegmentedControl', () => {
  it('すべてのオプションが表示される', () => {
    render(<SegmentedControl options={options} value="day" onChange={() => {}} />);
    expect(screen.getByText('日')).toBeInTheDocument();
    expect(screen.getByText('週')).toBeInTheDocument();
    expect(screen.getByText('月')).toBeInTheDocument();
  });

  it('選択中のオプションがハイライトされる', () => {
    render(<SegmentedControl options={options} value="week" onChange={() => {}} />);
    const weekBtn = screen.getByText('週').closest('button');
    expect(weekBtn?.className).toContain('bg-blue-600');
  });

  it('非選択のオプションはハイライトされない', () => {
    render(<SegmentedControl options={options} value="week" onChange={() => {}} />);
    const dayBtn = screen.getByText('日').closest('button');
    expect(dayBtn?.className).not.toContain('bg-blue-600');
  });

  it('クリックでonChangeが呼ばれる', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<SegmentedControl options={options} value="day" onChange={onChange} />);
    await user.click(screen.getByText('月'));
    expect(onChange).toHaveBeenCalledWith('month');
  });

  it('disabled状態で全ボタンが無効になる', () => {
    render(<SegmentedControl options={options} value="day" onChange={() => {}} disabled />);
    const buttons = screen.getAllByRole('button');
    buttons.forEach((btn) => expect(btn).toBeDisabled());
  });

  it('smサイズが適用される', () => {
    const { container } = render(<SegmentedControl options={options} value="day" onChange={() => {}} size="sm" />);
    expect(container.querySelector('.text-xs')).toBeInTheDocument();
  });

  it('lgサイズが適用される', () => {
    const { container } = render(<SegmentedControl options={options} value="day" onChange={() => {}} size="lg" />);
    expect(container.querySelector('.text-base')).toBeInTheDocument();
  });

  it('mdサイズがデフォルト', () => {
    const { container } = render(<SegmentedControl options={options} value="day" onChange={() => {}} />);
    expect(container.querySelector('.text-sm')).toBeInTheDocument();
  });

  it('ラベルが表示される', () => {
    render(<SegmentedControl options={options} value="day" onChange={() => {}} label="期間" />);
    expect(screen.getByText('期間')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<SegmentedControl options={options} value="day" onChange={() => {}} className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
