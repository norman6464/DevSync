import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import ScrollToTop from '../ScrollToTop';

describe('ScrollToTop', () => {
  it('初期状態で非表示', () => {
    render(<ScrollToTop />);
    expect(screen.queryByLabelText('トップへ戻る')).not.toBeInTheDocument();
  });

  it('スクロール後に表示される', () => {
    render(<ScrollToTop threshold={100} />);
    Object.defineProperty(window, 'scrollY', { value: 200, writable: true });
    fireEvent.scroll(window);
    expect(screen.getByLabelText('トップへ戻る')).toBeInTheDocument();
  });

  it('スクロール戻りで非表示に戻る', () => {
    render(<ScrollToTop threshold={100} />);
    Object.defineProperty(window, 'scrollY', { value: 200, writable: true });
    fireEvent.scroll(window);
    Object.defineProperty(window, 'scrollY', { value: 50, writable: true });
    fireEvent.scroll(window);
    expect(screen.queryByLabelText('トップへ戻る')).not.toBeInTheDocument();
  });

  it('クリックでwindow.scrollToが呼ばれる', async () => {
    const scrollToMock = vi.fn();
    window.scrollTo = scrollToMock;
    const user = userEvent.setup();
    render(<ScrollToTop threshold={100} />);
    Object.defineProperty(window, 'scrollY', { value: 200, writable: true });
    fireEvent.scroll(window);
    await user.click(screen.getByLabelText('トップへ戻る'));
    expect(scrollToMock).toHaveBeenCalledWith({ top: 0, behavior: 'smooth' });
  });

  it('デフォルトしきい値は300', () => {
    render(<ScrollToTop />);
    Object.defineProperty(window, 'scrollY', { value: 200, writable: true });
    fireEvent.scroll(window);
    expect(screen.queryByLabelText('トップへ戻る')).not.toBeInTheDocument();
  });

  it('カスタムしきい値が適用される', () => {
    render(<ScrollToTop threshold={50} />);
    Object.defineProperty(window, 'scrollY', { value: 60, writable: true });
    fireEvent.scroll(window);
    expect(screen.getByLabelText('トップへ戻る')).toBeInTheDocument();
  });

  it('ボタンにアイコンが表示される', () => {
    render(<ScrollToTop threshold={0} />);
    Object.defineProperty(window, 'scrollY', { value: 10, writable: true });
    fireEvent.scroll(window);
    const btn = screen.getByLabelText('トップへ戻る');
    expect(btn.querySelector('svg')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    render(<ScrollToTop threshold={0} className="custom-class" />);
    Object.defineProperty(window, 'scrollY', { value: 10, writable: true });
    fireEvent.scroll(window);
    expect(screen.getByLabelText('トップへ戻る').closest('.custom-class')).toBeInTheDocument();
  });
});
