import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/react';
import { Skeleton, PostCardSkeleton, UserCardSkeleton } from '../Skeleton';

describe('Skeleton', () => {
  it('デフォルトでレンダリングされる', () => {
    const { container } = render(<Skeleton />);
    const skeleton = container.firstChild;
    expect(skeleton).toHaveClass('animate-pulse', 'bg-gray-800', 'rounded');
  });

  it('カスタムclassNameが適用される', () => {
    const { container } = render(<Skeleton className="h-4 w-24" />);
    const skeleton = container.firstChild;
    expect(skeleton).toHaveClass('h-4', 'w-24');
  });
});

describe('PostCardSkeleton', () => {
  it('投稿カード用スケルトンがレンダリングされる', () => {
    const { container } = render(<PostCardSkeleton />);
    const card = container.querySelector('.bg-gray-900');
    expect(card).toBeInTheDocument();
    expect(card).toHaveClass('border', 'border-gray-800', 'rounded-md', 'p-5');
  });

  it('複数のスケルトン要素が含まれる', () => {
    const { container } = render(<PostCardSkeleton />);
    const skeletons = container.querySelectorAll('.animate-pulse');
    expect(skeletons.length).toBeGreaterThan(5);
  });
});

describe('UserCardSkeleton', () => {
  it('ユーザーカード用スケルトンがレンダリングされる', () => {
    const { container } = render(<UserCardSkeleton />);
    const card = container.querySelector('.bg-gray-900');
    expect(card).toBeInTheDocument();
    expect(card).toHaveClass('border', 'border-gray-800', 'rounded-md', 'p-4');
  });

  it('複数のスケルトン要素が含まれる', () => {
    const { container } = render(<UserCardSkeleton />);
    const skeletons = container.querySelectorAll('.animate-pulse');
    expect(skeletons.length).toBeGreaterThan(4);
  });
});
