import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import FeatureCard from '../FeatureCard';

describe('FeatureCard', () => {
  it('タイトルが表示される', () => {
    render(<FeatureCard title="高速ビルド" description="Viteによる高速ビルド" />);
    expect(screen.getByText('高速ビルド')).toBeInTheDocument();
  });

  it('説明が表示される', () => {
    render(<FeatureCard title="高速ビルド" description="Viteによる高速ビルド" />);
    expect(screen.getByText('Viteによる高速ビルド')).toBeInTheDocument();
  });

  it('アイコンが表示される', () => {
    render(<FeatureCard title="テスト" description="説明" icon="🚀" />);
    expect(screen.getByText('🚀')).toBeInTheDocument();
  });

  it('バッジが表示される', () => {
    render(<FeatureCard title="テスト" description="説明" badge="NEW" />);
    expect(screen.getByText('NEW')).toBeInTheDocument();
  });

  it('バッジがない場合は非表示', () => {
    render(<FeatureCard title="テスト" description="説明" />);
    expect(screen.queryByText('NEW')).not.toBeInTheDocument();
  });

  it('クリックでonClickが呼ばれる', async () => {
    const onClick = vi.fn();
    const user = userEvent.setup();
    render(<FeatureCard title="テスト" description="説明" onClick={onClick} />);
    await user.click(screen.getByText('テスト'));
    expect(onClick).toHaveBeenCalled();
  });

  it('クリック可能時にcursorがpointer', () => {
    const { container } = render(<FeatureCard title="テスト" description="説明" onClick={() => {}} />);
    expect(container.querySelector('.cursor-pointer')).toBeInTheDocument();
  });

  it('disabled時にopacityが下がる', () => {
    const { container } = render(<FeatureCard title="テスト" description="説明" disabled />);
    expect(container.querySelector('.opacity-50')).toBeInTheDocument();
  });

  it('ホバーエフェクトが適用される', () => {
    const { container } = render(<FeatureCard title="テスト" description="説明" />);
    expect(container.querySelector('.hover\\:border-gray-600')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<FeatureCard title="テスト" description="説明" className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
