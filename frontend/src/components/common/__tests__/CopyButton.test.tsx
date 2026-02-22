import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import CopyButton from '../CopyButton';

describe('CopyButton', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    Object.assign(navigator, {
      clipboard: {
        writeText: vi.fn().mockResolvedValue(undefined),
      },
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
    const user = userEvent.setup({ advanceTimers: vi.advanceTimersByTime });
    render(<CopyButton text="コピーテキスト" />);

    await user.click(screen.getByRole('button'));

    expect(navigator.clipboard.writeText).toHaveBeenCalledWith('コピーテキスト');
  });

  it('コピー成功後にチェックアイコンが表示される', async () => {
    const user = userEvent.setup({ advanceTimers: vi.advanceTimersByTime });
    const { container } = render(<CopyButton text="テスト" />);

    await user.click(screen.getByRole('button'));

    expect(container.querySelector('.lucide-check')).toBeInTheDocument();
  });

  it('一定時間後に元のアイコンに戻る', async () => {
    const user = userEvent.setup({ advanceTimers: vi.advanceTimersByTime });
    const { container } = render(<CopyButton text="テスト" />);

    await user.click(screen.getByRole('button'));
    expect(container.querySelector('.lucide-check')).toBeInTheDocument();

    vi.advanceTimersByTime(2000);

    expect(container.querySelector('.lucide-copy')).toBeInTheDocument();
  });

  it('ラベルが表示される', () => {
    render(<CopyButton text="テスト" label="コピー" />);

    expect(screen.getByText('コピー')).toBeInTheDocument();
  });

  it('コピー成功後にラベルが変わる', async () => {
    const user = userEvent.setup({ advanceTimers: vi.advanceTimersByTime });
    render(<CopyButton text="テスト" label="コピー" successLabel="コピー済み" />);

    await user.click(screen.getByRole('button'));

    expect(screen.getByText('コピー済み')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<CopyButton text="テスト" className="custom-class" />);

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('onCopyコールバックが呼ばれる', async () => {
    const onCopy = vi.fn();
    const user = userEvent.setup({ advanceTimers: vi.advanceTimersByTime });
    render(<CopyButton text="テスト" onCopy={onCopy} />);

    await user.click(screen.getByRole('button'));

    expect(onCopy).toHaveBeenCalledTimes(1);
  });
});
