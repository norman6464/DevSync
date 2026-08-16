import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/react';
import LoadingSpinner from '../LoadingSpinner';

describe('LoadingSpinner', () => {
  it('デフォルトサイズでレンダリングされる', () => {
    const { container } = render(<LoadingSpinner />);
    const spinner = container.querySelector('.animate-spin');
    expect(spinner).toBeInTheDocument();
    expect(spinner).toHaveClass('w-8', 'h-8');
  });

  it('小サイズ（sm）でレンダリングされる', () => {
    const { container } = render(<LoadingSpinner size="sm" />);
    const spinner = container.querySelector('.animate-spin');
    expect(spinner).toHaveClass('w-5', 'h-5');
  });

  it('大サイズ（lg）でレンダリングされる', () => {
    const { container } = render(<LoadingSpinner size="lg" />);
    const spinner = container.querySelector('.animate-spin');
    expect(spinner).toHaveClass('w-12', 'h-12');
  });

  it('プライマリカラーでレンダリングされる', () => {
    const { container } = render(<LoadingSpinner color="primary" />);
    const spinner = container.querySelector('.animate-spin');
    expect(spinner).toHaveClass('border-blue-500');
  });

  it('セカンダリカラーでレンダリングされる', () => {
    const { container } = render(<LoadingSpinner color="secondary" />);
    const spinner = container.querySelector('.animate-spin');
    expect(spinner).toHaveClass('border-gray-400');
  });

  it('ホワイトカラーでレンダリングされる', () => {
    const { container } = render(<LoadingSpinner color="white" />);
    const spinner = container.querySelector('.animate-spin');
    expect(spinner).toHaveClass('border-white');
  });

  it('カスタムclassNameが適用される', () => {
    const { container } = render(<LoadingSpinner className="my-custom-class" />);
    const wrapper = container.firstChild;
    expect(wrapper).toHaveClass('my-custom-class');
  });
});
