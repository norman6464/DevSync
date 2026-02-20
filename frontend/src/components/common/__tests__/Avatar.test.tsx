import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import Avatar from '../Avatar';

describe('Avatar', () => {
  it('avatarUrlなしの場合イニシャルを表示する', () => {
    render(<Avatar name="Alice" />);
    expect(screen.getByText('A')).toBeInTheDocument();
  });

  it('名前の先頭文字を大文字で表示する', () => {
    render(<Avatar name="bob" />);
    expect(screen.getByText('B')).toBeInTheDocument();
  });

  it('avatarUrlありの場合img要素を表示する', () => {
    render(<Avatar name="Alice" avatarUrl="https://example.com/avatar.png" />);
    const img = screen.getByRole('img');
    expect(img).toHaveAttribute('src', 'https://example.com/avatar.png');
    expect(img).toHaveAttribute('alt', 'Alice');
  });

  it('デフォルトサイズはmdクラス', () => {
    render(<Avatar name="Alice" />);
    const el = screen.getByText('A');
    expect(el.className).toContain('w-10');
    expect(el.className).toContain('h-10');
  });

  it('size=smの場合smクラスが適用される', () => {
    render(<Avatar name="Alice" size="sm" />);
    const el = screen.getByText('A');
    expect(el.className).toContain('w-8');
    expect(el.className).toContain('h-8');
  });

  it('size=lgの場合lgクラスが適用される', () => {
    render(<Avatar name="Alice" size="lg" />);
    const el = screen.getByText('A');
    expect(el.className).toContain('w-16');
    expect(el.className).toContain('h-16');
  });

  it('isOnline未指定の場合インジケーターが表示されない', () => {
    const { container } = render(<Avatar name="Alice" />);
    expect(container.querySelector('[aria-hidden="true"]')).toBeNull();
  });

  it('isOnline=trueの場合緑色のインジケーターが表示される', () => {
    const { container } = render(<Avatar name="Alice" isOnline={true} />);
    const dot = container.querySelector('[aria-hidden="true"]');
    expect(dot).not.toBeNull();
    expect(dot!.className).toContain('bg-green-500');
  });

  it('isOnline=falseの場合グレーのインジケーターが表示される', () => {
    const { container } = render(<Avatar name="Alice" isOnline={false} />);
    const dot = container.querySelector('[aria-hidden="true"]');
    expect(dot).not.toBeNull();
    expect(dot!.className).toContain('bg-gray-500');
  });

  it('不正なavatarUrlの場合イニシャルにフォールバックする', () => {
    render(<Avatar name="Alice" avatarUrl="javascript:alert(1)" />);
    expect(screen.getByText('A')).toBeInTheDocument();
  });
});
