import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import DatePicker from '../DatePicker';

describe('DatePicker', () => {
  it('入力フィールドが表示される', () => {
    render(<DatePicker value="" onChange={() => {}} />);

    expect(screen.getByRole('textbox')).toBeInTheDocument();
  });

  it('プレースホルダーが表示される', () => {
    render(<DatePicker value="" onChange={() => {}} placeholder="日付を選択" />);

    expect(screen.getByPlaceholderText('日付を選択')).toBeInTheDocument();
  });

  it('値が表示される', () => {
    render(<DatePicker value="2026-01-15" onChange={() => {}} />);

    expect(screen.getByDisplayValue('2026-01-15')).toBeInTheDocument();
  });

  it('カレンダーアイコンが表示される', () => {
    const { container } = render(<DatePicker value="" onChange={() => {}} />);

    expect(container.querySelector('.lucide-calendar')).toBeInTheDocument();
  });

  it('クリックでカレンダーが開く', async () => {
    const user = userEvent.setup();
    render(<DatePicker value="" onChange={() => {}} />);

    await user.click(screen.getByRole('textbox'));

    expect(screen.getByText('日')).toBeInTheDocument();
    expect(screen.getByText('月')).toBeInTheDocument();
  });

  it('日付クリックで値が変わる', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<DatePicker value="2026-01-15" onChange={onChange} />);

    await user.click(screen.getByRole('textbox'));

    const day10 = screen.getByText('10');
    await user.click(day10);

    expect(onChange).toHaveBeenCalled();
  });

  it('前月ボタンで月が切り替わる', async () => {
    const user = userEvent.setup();
    render(<DatePicker value="2026-02-15" onChange={() => {}} />);

    await user.click(screen.getByRole('textbox'));

    const prevButton = screen.getByLabelText('前月');
    await user.click(prevButton);

    expect(screen.getByText(/1月/)).toBeInTheDocument();
  });

  it('次月ボタンで月が切り替わる', async () => {
    const user = userEvent.setup();
    render(<DatePicker value="2026-01-15" onChange={() => {}} />);

    await user.click(screen.getByRole('textbox'));

    const nextButton = screen.getByLabelText('次月');
    await user.click(nextButton);

    expect(screen.getByText(/2月/)).toBeInTheDocument();
  });

  it('無効状態が適用される', () => {
    render(<DatePicker value="" onChange={() => {}} disabled />);

    expect(screen.getByRole('textbox')).toBeDisabled();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<DatePicker value="" onChange={() => {}} className="custom-class" />);

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
