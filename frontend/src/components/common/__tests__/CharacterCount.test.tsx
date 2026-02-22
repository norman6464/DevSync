import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import CharacterCount from '../CharacterCount';

describe('CharacterCount', () => {
  it('現在の文字数が表示される', () => {
    render(<CharacterCount current={50} max={200} />);

    expect(screen.getByText('50')).toBeInTheDocument();
  });

  it('最大文字数が表示される', () => {
    render(<CharacterCount current={50} max={200} />);

    expect(screen.getByText('/ 200')).toBeInTheDocument();
  });

  it('文字数超過時に警告色が適用される', () => {
    const { container } = render(<CharacterCount current={250} max={200} />);

    expect(container.querySelector('.text-red-400')).toBeInTheDocument();
  });

  it('80%以上で注意色が適用される', () => {
    const { container } = render(<CharacterCount current={170} max={200} />);

    expect(container.querySelector('.text-yellow-400')).toBeInTheDocument();
  });

  it('80%未満で通常色が適用される', () => {
    const { container } = render(<CharacterCount current={50} max={200} />);

    expect(container.querySelector('.text-gray-400')).toBeInTheDocument();
  });

  it('プログレスバーが表示される', () => {
    const { container } = render(<CharacterCount current={100} max={200} showBar />);

    const bar = container.querySelector('.h-1');
    expect(bar).toBeInTheDocument();
  });

  it('プログレスバーの幅が正しい', () => {
    const { container } = render(<CharacterCount current={100} max={200} showBar />);

    const progress = container.querySelector('[style*="width: 50%"]');
    expect(progress).toBeInTheDocument();
  });

  it('残り文字数モードが表示される', () => {
    render(<CharacterCount current={50} max={200} showRemaining />);

    expect(screen.getByText('150')).toBeInTheDocument();
    expect(screen.getByText('残り')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<CharacterCount current={50} max={200} className="custom-class" />);

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('文字数が0の場合も表示される', () => {
    render(<CharacterCount current={0} max={200} />);

    expect(screen.getByText('0')).toBeInTheDocument();
  });
});
