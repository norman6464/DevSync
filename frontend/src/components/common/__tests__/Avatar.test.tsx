import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import Avatar from '../Avatar';

describe('Avatar', () => {
  it('画像が表示される', () => {
    render(<Avatar src="/avatar.jpg" alt="ユーザー" />);
    const img = screen.getByRole('img');
    expect(img).toHaveAttribute('src', '/avatar.jpg');
    expect(img).toHaveAttribute('alt', 'ユーザー');
  });

  it('画像がない場合イニシャルが表示される', () => {
    render(<Avatar name="田中太郎" />);
    expect(screen.getByText('田')).toBeInTheDocument();
  });

  it('名前がない場合デフォルトアイコンが表示される', () => {
    const { container } = render(<Avatar />);
    expect(container.querySelector('.lucide-user')).toBeInTheDocument();
  });

  it('smサイズが適用される', () => {
    const { container } = render(<Avatar name="T" size="sm" />);
    expect(container.querySelector('.w-8')).toBeInTheDocument();
  });

  it('mdサイズが適用される（デフォルト）', () => {
    const { container } = render(<Avatar name="T" />);
    expect(container.querySelector('.w-10')).toBeInTheDocument();
  });

  it('lgサイズが適用される', () => {
    const { container } = render(<Avatar name="T" size="lg" />);
    expect(container.querySelector('.w-12')).toBeInTheDocument();
  });

  it('xlサイズが適用される', () => {
    const { container } = render(<Avatar name="T" size="xl" />);
    expect(container.querySelector('.w-16')).toBeInTheDocument();
  });

  it('オンラインステータスが表示される', () => {
    const { container } = render(<Avatar name="T" online />);
    expect(container.querySelector('.bg-green-500')).toBeInTheDocument();
  });

  it('オフラインステータスが表示される', () => {
    const { container } = render(<Avatar name="T" online={false} />);
    expect(container.querySelector('.bg-gray-500')).toBeInTheDocument();
  });

  it('角丸がデフォルト', () => {
    const { container } = render(<Avatar name="T" />);
    expect(container.querySelector('.rounded-full')).toBeInTheDocument();
  });

  it('四角にできる', () => {
    const { container } = render(<Avatar name="T" rounded={false} />);
    expect(container.querySelector('.rounded-lg')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<Avatar name="T" className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
