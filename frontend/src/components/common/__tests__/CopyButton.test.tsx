import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent, act } from '@testing-library/react';
import CopyButton from '../CopyButton';

describe('CopyButton', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    Object.defineProperty(navigator, 'clipboard', {
      value: {
        writeText: vi.fn().mockResolvedValue(undefined),
      },
      writable: true,
      configurable: true,
    });
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('コピーボタンが表示される', () => {
    render(<CopyButton text="コピーするテキスト" />);

    expect(screen.getByRole('button')).toBeInTheDocument();
  });

  it('コピーアイコンが表示される', () => {
    const { container } = render(<CopyButton text="テスト" />);

    expect(container.querySelector('.lucide-copy')).toBeInTheDocument();
  });

  it('クリックでクリップボードにコピーされる', async () => {
    render(<CopyButton text="コピーテキスト" />);

    await act(async () => {
      fireEvent.click(screen.getByRole('button'));
    });

    expect(navigator.clipboard.writeText).toHaveBeenCalledWith('コピーテキスト');
  });

  it('コピー成功後にチェックアイコンが表示される', async () => {
    const { container } = render(<CopyButton text="テスト" />);

    await act(async () => {
      fireEvent.click(screen.getByRole('button'));
    });

    expect(container.querySelector('.lucide-check')).toBeInTheDocument();
  });

  it('一定時間後に元のアイコンに戻る', async () => {
    const { container } = render(<CopyButton text="テスト" />);

    await act(async () => {
      fireEvent.click(screen.getByRole('button'));
    });
    expect(container.querySelector('.lucide-check')).toBeInTheDocument();

    act(() => {
      vi.advanceTimersByTime(2000);
    });

    expect(container.querySelector('.lucide-copy')).toBeInTheDocument();
  });

  it('ラベルが表示される', () => {
    render(<CopyButton text="テスト" label="コピー" />);

    expect(screen.getByText('コピー')).toBeInTheDocument();
  });

  it('コピー成功後にラベルが変わる', async () => {
    render(<CopyButton text="テスト" label="コピー" successLabel="コピー済み" />);

    await act(async () => {
      fireEvent.click(screen.getByRole('button'));
    });

    expect(screen.getByText('コピー済み')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<CopyButton text="テスト" className="custom-class" />);

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('onCopyコールバックが呼ばれる', async () => {
    const onCopy = vi.fn();
    render(<CopyButton text="テスト" onCopy={onCopy} />);

    await act(async () => {
      fireEvent.click(screen.getByRole('button'));
    });

    expect(onCopy).toHaveBeenCalledTimes(1);
  });
});
