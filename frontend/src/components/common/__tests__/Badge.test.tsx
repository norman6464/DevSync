import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import Badge from '../Badge';

describe('Badge', () => {
  it('バッジが表示される', () => {
    render(<Badge>テスト</Badge>);
    expect(screen.getByText('テスト')).toBeInTheDocument();
  });

  it('初級バッジが緑色で表示される', () => {
    render(<Badge variant="beginner">初級</Badge>);

    const badge = screen.getByText('初級');
    expect(badge).toHaveClass('bg-green-500/20', 'text-green-400');
  });

  it('中級バッジが青色で表示される', () => {
    render(<Badge variant="intermediate">中級</Badge>);

    const badge = screen.getByText('中級');
    expect(badge).toHaveClass('bg-blue-500/20', 'text-blue-400');
  });

  it('上級バッジが紫色で表示される', () => {
    render(<Badge variant="advanced">上級</Badge>);

    const badge = screen.getByText('上級');
    expect(badge).toHaveClass('bg-purple-500/20', 'text-purple-400');
  });

  it('エキスパートバッジがオレンジ色で表示される', () => {
    render(<Badge variant="expert">エキスパート</Badge>);

    const badge = screen.getByText('エキスパート');
    expect(badge).toHaveClass('bg-orange-500/20', 'text-orange-400');
  });

  it('デフォルトバッジが灰色で表示される', () => {
    render(<Badge variant="default">デフォルト</Badge>);

    const badge = screen.getByText('デフォルト');
    expect(badge).toHaveClass('bg-gray-500/20', 'text-gray-400');
  });

  it('小サイズのバッジが表示される', () => {
    render(<Badge size="sm">小</Badge>);

    const badge = screen.getByText('小');
    expect(badge).toHaveClass('text-xs', 'px-2', 'py-0.5');
  });

  it('中サイズのバッジが表示される', () => {
    render(<Badge size="md">中</Badge>);

    const badge = screen.getByText('中');
    expect(badge).toHaveClass('text-sm', 'px-2.5', 'py-1');
  });

  it('大サイズのバッジが表示される', () => {
    render(<Badge size="lg">大</Badge>);

    const badge = screen.getByText('大');
    expect(badge).toHaveClass('text-base', 'px-3', 'py-1.5');
  });

  it('アイコン付きバッジが表示される', () => {
    render(
      <Badge icon={<span data-testid="icon">★</span>}>アイコン</Badge>
    );

    expect(screen.getByTestId('icon')).toBeInTheDocument();
    expect(screen.getByText('アイコン')).toBeInTheDocument();
  });

  it('角が丸くなっている', () => {
    render(<Badge>丸角</Badge>);

    const badge = screen.getByText('丸角');
    expect(badge).toHaveClass('rounded-full');
  });

  it('フォントが太字になっている', () => {
    render(<Badge>太字</Badge>);

    const badge = screen.getByText('太字');
    expect(badge).toHaveClass('font-medium');
  });

  it('カスタムクラス名が適用される', () => {
    render(<Badge className="custom-class">カスタム</Badge>);

    const badge = screen.getByText('カスタム');
    expect(badge).toHaveClass('custom-class');
  });

  it('デフォルトサイズは中サイズ', () => {
    render(<Badge>デフォルト</Badge>);

    const badge = screen.getByText('デフォルト');
    expect(badge).toHaveClass('text-sm', 'px-2.5', 'py-1');
  });

  it('デフォルトバリアントはdefault', () => {
    render(<Badge>バリアント</Badge>);

    const badge = screen.getByText('バリアント');
    expect(badge).toHaveClass('bg-gray-500/20', 'text-gray-400');
  });

  it('複数のバッジが独立して表示される', () => {
    render(
      <>
        <Badge variant="beginner">初級</Badge>
        <Badge variant="intermediate">中級</Badge>
        <Badge variant="advanced">上級</Badge>
      </>
    );

    expect(screen.getByText('初級')).toBeInTheDocument();
    expect(screen.getByText('中級')).toBeInTheDocument();
    expect(screen.getByText('上級')).toBeInTheDocument();
  });
});
