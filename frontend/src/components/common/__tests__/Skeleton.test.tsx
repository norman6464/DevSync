import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/react';
import Skeleton from '../Skeleton';

describe('Skeleton', () => {
  it('スケルトンが表示される', () => {
    const { container } = render(<Skeleton />);

    const skeleton = container.querySelector('.animate-pulse');
    expect(skeleton).toBeInTheDocument();
  });

  it('テキスト形状が表示される', () => {
    const { container } = render(<Skeleton variant="text" />);

    const skeleton = container.querySelector('.h-4');
    expect(skeleton).toBeInTheDocument();
  });

  it('サークル形状が表示される', () => {
    const { container } = render(<Skeleton variant="circle" />);

    const skeleton = container.querySelector('.rounded-full');
    expect(skeleton).toBeInTheDocument();
  });

  it('レクタングル形状が表示される', () => {
    const { container } = render(<Skeleton variant="rectangle" />);

    const skeleton = container.querySelector('.rounded');
    expect(skeleton).toBeInTheDocument();
  });

  it('背景色がある', () => {
    const { container } = render(<Skeleton />);

    const skeleton = container.querySelector('.bg-gray-800');
    expect(skeleton).toBeInTheDocument();
  });

  it('幅が指定できる', () => {
    const { container } = render(<Skeleton width="200px" />);

    const skeleton = container.querySelector('.animate-pulse');
    expect(skeleton).toHaveStyle({ width: '200px' });
  });

  it('高さが指定できる', () => {
    const { container } = render(<Skeleton height="100px" />);

    const skeleton = container.querySelector('.animate-pulse');
    expect(skeleton).toHaveStyle({ height: '100px' });
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<Skeleton className="custom-class" />);

    const skeleton = container.querySelector('.custom-class');
    expect(skeleton).toBeInTheDocument();
  });

  it('デフォルトはテキスト形状', () => {
    const { container } = render(<Skeleton />);

    const skeleton = container.querySelector('.h-4');
    expect(skeleton).toBeInTheDocument();
  });

  it('複数のスケルトンが表示される', () => {
    const { container } = render(
      <>
        <Skeleton />
        <Skeleton />
        <Skeleton />
      </>
    );

    const skeletons = container.querySelectorAll('.animate-pulse');
    expect(skeletons.length).toBe(3);
  });
});
