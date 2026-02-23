import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import Timeline from '../Timeline';

const events = [
  { id: '1', title: 'プロジェクト開始', date: '2026-01-01', description: '初回ミーティング' },
  { id: '2', title: '設計完了', date: '2026-01-15' },
  { id: '3', title: 'リリース', date: '2026-02-01', active: true },
];

describe('Timeline', () => {
  it('全てのイベントタイトルが表示される', () => {
    render(<Timeline events={events} />);

    expect(screen.getByText('プロジェクト開始')).toBeInTheDocument();
    expect(screen.getByText('設計完了')).toBeInTheDocument();
    expect(screen.getByText('リリース')).toBeInTheDocument();
  });

  it('日付が表示される', () => {
    render(<Timeline events={events} />);

    expect(screen.getByText('2026-01-01')).toBeInTheDocument();
    expect(screen.getByText('2026-01-15')).toBeInTheDocument();
  });

  it('説明文が表示される', () => {
    render(<Timeline events={events} />);

    expect(screen.getByText('初回ミーティング')).toBeInTheDocument();
  });

  it('アクティブなイベントがハイライトされる', () => {
    const { container } = render(<Timeline events={events} />);

    const activeDots = container.querySelectorAll('.bg-blue-500');
    expect(activeDots.length).toBeGreaterThanOrEqual(1);
  });

  it('非アクティブなイベントはグレー', () => {
    const { container } = render(<Timeline events={events} />);

    const grayDots = container.querySelectorAll('.bg-gray-600');
    expect(grayDots.length).toBe(2);
  });

  it('接続線が表示される', () => {
    const { container } = render(<Timeline events={events} />);

    const lines = container.querySelectorAll('.w-0\\.5');
    expect(lines.length).toBe(2);
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<Timeline events={events} className="custom-class" />);

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('単一イベントでも表示される', () => {
    render(<Timeline events={[{ id: '1', title: '唯一', date: '2026-01-01' }]} />);

    expect(screen.getByText('唯一')).toBeInTheDocument();
  });

  it('ReactNodeをコンテンツに使用できる', () => {
    const richEvents = [
      { id: '1', title: 'リッチ', date: '2026-01-01', content: <span data-testid="rich">カスタム</span> },
    ];
    render(<Timeline events={richEvents} />);

    expect(screen.getByTestId('rich')).toBeInTheDocument();
  });

  it('ドットが全イベントに表示される', () => {
    const { container } = render(<Timeline events={events} />);

    const dots = container.querySelectorAll('.rounded-full');
    expect(dots.length).toBe(3);
  });
});
