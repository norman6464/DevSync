import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import CountdownTimer from '../CountdownTimer';

describe('CountdownTimer', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date('2026-01-01T00:00:00Z'));
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('残り日数が表示される', () => {
    render(<CountdownTimer targetDate="2026-01-11T00:00:00Z" />);

    expect(screen.getByText('10')).toBeInTheDocument();
    expect(screen.getByText('日')).toBeInTheDocument();
  });

  it('残り時間が表示される', () => {
    render(<CountdownTimer targetDate="2026-01-01T05:00:00Z" />);

    expect(screen.getByText('5')).toBeInTheDocument();
    expect(screen.getByText('時間')).toBeInTheDocument();
  });

  it('残り分が表示される', () => {
    render(<CountdownTimer targetDate="2026-01-01T00:30:00Z" />);

    expect(screen.getByText('30')).toBeInTheDocument();
    expect(screen.getByText('分')).toBeInTheDocument();
  });

  it('残り秒が表示される', () => {
    render(<CountdownTimer targetDate="2026-01-01T00:00:45Z" />);

    expect(screen.getByText('45')).toBeInTheDocument();
    expect(screen.getByText('秒')).toBeInTheDocument();
  });

  it('期限切れ時にメッセージが表示される', () => {
    render(<CountdownTimer targetDate="2025-12-31T00:00:00Z" expiredMessage="期限切れ" />);

    expect(screen.getByText('期限切れ')).toBeInTheDocument();
  });

  it('デフォルトの期限切れメッセージ', () => {
    render(<CountdownTimer targetDate="2025-12-31T00:00:00Z" />);

    expect(screen.getByText('期限終了')).toBeInTheDocument();
  });

  it('ラベルが表示される', () => {
    render(<CountdownTimer targetDate="2026-01-02T00:00:00Z" label="締め切りまで" />);

    expect(screen.getByText('締め切りまで')).toBeInTheDocument();
  });

  it('コンパクトモードで表示される', () => {
    const { container } = render(
      <CountdownTimer targetDate="2026-01-02T00:00:00Z" compact />
    );

    expect(container.querySelector('.text-sm')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(
      <CountdownTimer targetDate="2026-01-02T00:00:00Z" className="custom-class" />
    );

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('onExpireコールバックが呼ばれる', () => {
    const onExpire = vi.fn();
    render(<CountdownTimer targetDate="2025-12-31T00:00:00Z" onExpire={onExpire} />);

    expect(onExpire).toHaveBeenCalledTimes(1);
  });
});
