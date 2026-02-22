import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Toggle from '../Toggle';

describe('Toggle', () => {
  it('トグルが表示される', () => {
    render(<Toggle checked={false} onChange={() => {}} />);

    expect(screen.getByRole('switch')).toBeInTheDocument();
  });

  it('ONの状態が表示される', () => {
    const { container } = render(<Toggle checked={true} onChange={() => {}} />);

    expect(container.querySelector('.bg-blue-600')).toBeInTheDocument();
  });

  it('OFFの状態が表示される', () => {
    const { container } = render(<Toggle checked={false} onChange={() => {}} />);

    expect(container.querySelector('.bg-gray-600')).toBeInTheDocument();
  });

  it('クリックで状態が切り替わる', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<Toggle checked={false} onChange={onChange} />);

    await user.click(screen.getByRole('switch'));

    expect(onChange).toHaveBeenCalledWith(true);
  });

  it('ON状態でクリックするとOFFになる', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<Toggle checked={true} onChange={onChange} />);

    await user.click(screen.getByRole('switch'));

    expect(onChange).toHaveBeenCalledWith(false);
  });

  it('ラベルが表示される', () => {
    render(<Toggle checked={false} onChange={() => {}} label="通知" />);

    expect(screen.getByText('通知')).toBeInTheDocument();
  });

  it('smサイズが適用される', () => {
    const { container } = render(<Toggle checked={false} onChange={() => {}} size="sm" />);

    expect(container.querySelector('.w-8')).toBeInTheDocument();
  });

  it('lgサイズが適用される', () => {
    const { container } = render(<Toggle checked={false} onChange={() => {}} size="lg" />);

    expect(container.querySelector('.w-14')).toBeInTheDocument();
  });

  it('無効状態が適用される', () => {
    render(<Toggle checked={false} onChange={() => {}} disabled />);

    expect(screen.getByRole('switch')).toBeDisabled();
  });

  it('無効状態ではクリックしても変化しない', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<Toggle checked={false} onChange={onChange} disabled />);

    await user.click(screen.getByRole('switch'));

    expect(onChange).not.toHaveBeenCalled();
  });

  it('aria-checkedが正しく設定される', () => {
    render(<Toggle checked={true} onChange={() => {}} />);

    expect(screen.getByRole('switch')).toHaveAttribute('aria-checked', 'true');
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<Toggle checked={false} onChange={() => {}} className="custom-class" />);

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
