import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/react';
import InlineLoader from '../InlineLoader';

describe('InlineLoader', () => {
  it('デフォルトでレンダリングされる', () => {
    const { container } = render(<InlineLoader />);
    const spinner = container.querySelector('.animate-spin');
    expect(spinner).toBeInTheDocument();
  });

  it('小サイズでレンダリングされる', () => {
    const { container } = render(<InlineLoader size="sm" />);
    const spinner = container.querySelector('.animate-spin');
    expect(spinner).toHaveClass('w-4', 'h-4');
  });

  it('中サイズでレンダリングされる', () => {
    const { container } = render(<InlineLoader size="md" />);
    const spinner = container.querySelector('.animate-spin');
    expect(spinner).toHaveClass('w-5', 'h-5');
  });

  it('大サイズでレンダリングされる', () => {
    const { container } = render(<InlineLoader size="lg" />);
    const spinner = container.querySelector('.animate-spin');
    expect(spinner).toHaveClass('w-6', 'h-6');
  });

  it('カスタムclassNameが適用される', () => {
    const { container } = render(<InlineLoader className="ml-2" />);
    const wrapper = container.firstChild;
    expect(wrapper).toHaveClass('ml-2');
  });
});
