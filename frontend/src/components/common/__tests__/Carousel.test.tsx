import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Carousel from '../Carousel';

const items = ['スライド1', 'スライド2', 'スライド3'];

describe('Carousel', () => {
  it('最初のスライドが表示される', () => {
    render(<Carousel items={items} renderItem={(item) => <div>{item}</div>} />);
    expect(screen.getByText('スライド1')).toBeInTheDocument();
  });

  it('次へボタンで次のスライドが表示される', async () => {
    const user = userEvent.setup();
    render(<Carousel items={items} renderItem={(item) => <div>{item}</div>} />);
    await user.click(screen.getByLabelText('次へ'));
    expect(screen.getByText('スライド2')).toBeInTheDocument();
  });

  it('前へボタンで前のスライドが表示される', async () => {
    const user = userEvent.setup();
    render(<Carousel items={items} renderItem={(item) => <div>{item}</div>} />);
    await user.click(screen.getByLabelText('次へ'));
    await user.click(screen.getByLabelText('前へ'));
    expect(screen.getByText('スライド1')).toBeInTheDocument();
  });

  it('最初のスライドで前へボタンが無効（ループなし）', () => {
    render(<Carousel items={items} renderItem={(item) => <div>{item}</div>} />);
    expect(screen.getByLabelText('前へ')).toBeDisabled();
  });

  it('最後のスライドで次へボタンが無効（ループなし）', async () => {
    const user = userEvent.setup();
    render(<Carousel items={items} renderItem={(item) => <div>{item}</div>} />);
    await user.click(screen.getByLabelText('次へ'));
    await user.click(screen.getByLabelText('次へ'));
    expect(screen.getByLabelText('次へ')).toBeDisabled();
  });

  it('ループ有効時に最後から最初へ戻る', async () => {
    const user = userEvent.setup();
    render(<Carousel items={items} renderItem={(item) => <div>{item}</div>} loop />);
    await user.click(screen.getByLabelText('次へ'));
    await user.click(screen.getByLabelText('次へ'));
    await user.click(screen.getByLabelText('次へ'));
    expect(screen.getByText('スライド1')).toBeInTheDocument();
  });

  it('ドットインジケーターが表示される', () => {
    render(<Carousel items={items} renderItem={(item) => <div>{item}</div>} showDots />);
    const dots = screen.getAllByTestId('carousel-dot');
    expect(dots.length).toBe(3);
  });

  it('ドットクリックでスライドが切り替わる', async () => {
    const user = userEvent.setup();
    render(<Carousel items={items} renderItem={(item) => <div>{item}</div>} showDots />);
    await user.click(screen.getAllByTestId('carousel-dot')[2]);
    expect(screen.getByText('スライド3')).toBeInTheDocument();
  });

  it('アクティブなドットにスタイルが適用される', () => {
    render(<Carousel items={items} renderItem={(item) => <div>{item}</div>} showDots />);
    const dots = screen.getAllByTestId('carousel-dot');
    expect(dots[0].className).toContain('bg-blue-500');
  });

  it('onChangeコールバックが呼ばれる', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<Carousel items={items} renderItem={(item) => <div>{item}</div>} onChange={onChange} />);
    await user.click(screen.getByLabelText('次へ'));
    expect(onChange).toHaveBeenCalledWith(1);
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<Carousel items={items} renderItem={(item) => <div>{item}</div>} className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
