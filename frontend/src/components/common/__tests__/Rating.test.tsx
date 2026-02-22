import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Rating from '../Rating';

describe('Rating', () => {
  it('デフォルトで5つの星が表示される', () => {
    const { container } = render(<Rating value={0} onChange={() => {}} />);

    const stars = container.querySelectorAll('.lucide-star');
    expect(stars.length).toBe(5);
  });

  it('指定した数の星が表示される', () => {
    const { container } = render(<Rating value={0} onChange={() => {}} max={3} />);

    const stars = container.querySelectorAll('.lucide-star');
    expect(stars.length).toBe(3);
  });

  it('値に応じて星が塗りつぶされる', () => {
    const { container } = render(<Rating value={3} onChange={() => {}} />);

    const filledStars = container.querySelectorAll('.text-yellow-400');
    expect(filledStars.length).toBe(3);
  });

  it('クリックで値が変わる', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<Rating value={0} onChange={onChange} />);

    const buttons = screen.getAllByRole('button');
    await user.click(buttons[2]);

    expect(onChange).toHaveBeenCalledWith(3);
  });

  it('読み取り専用モードではクリックできない', () => {
    render(<Rating value={3} readOnly />);

    const buttons = screen.queryAllByRole('button');
    expect(buttons.length).toBe(0);
  });

  it('ホバーでプレビューが表示される', async () => {
    const { container } = render(<Rating value={1} onChange={() => {}} />);
    const user = userEvent.setup();

    const buttons = screen.getAllByRole('button');
    await user.hover(buttons[3]);

    const highlightedStars = container.querySelectorAll('.text-yellow-400');
    expect(highlightedStars.length).toBe(4);
  });

  it('ホバー解除で元の値に戻る', async () => {
    const { container } = render(<Rating value={2} onChange={() => {}} />);
    const user = userEvent.setup();

    const buttons = screen.getAllByRole('button');
    await user.hover(buttons[4]);
    await user.unhover(buttons[4]);

    const filledStars = container.querySelectorAll('.text-yellow-400');
    expect(filledStars.length).toBe(2);
  });

  it('smサイズが適用される', () => {
    const { container } = render(<Rating value={0} onChange={() => {}} size="sm" />);

    const star = container.querySelector('.w-4');
    expect(star).toBeInTheDocument();
  });

  it('lgサイズが適用される', () => {
    const { container } = render(<Rating value={0} onChange={() => {}} size="lg" />);

    const star = container.querySelector('.w-8');
    expect(star).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<Rating value={0} onChange={() => {}} className="custom-class" />);

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('同じ値をクリックすると0にリセットされる', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<Rating value={3} onChange={onChange} />);

    const buttons = screen.getAllByRole('button');
    await user.click(buttons[2]);

    expect(onChange).toHaveBeenCalledWith(0);
  });
});
