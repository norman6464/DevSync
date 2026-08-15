import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import Avatar from '../Avatar';

describe('Avatar', () => {
  it('画像が表示される', () => {
    render(<Avatar avatarUrl="/avatar.jpg" alt="ユーザー" />);
    const img = screen.getByRole('img');
    expect(img).toHaveAttribute('src', '/avatar.jpg');
    expect(img).toHaveAttribute('alt', 'ユーザー');
  });

  // 呼び出し側は一貫して avatarUrl で渡す。名前がずれると画像が黙って消える。
  it('画像 URL があればイニシャルではなく画像を出す', () => {
    render(<Avatar avatarUrl="/avatar.jpg" name="田中太郎" />);
    expect(screen.getByRole('img')).toHaveAttribute('src', '/avatar.jpg');
    expect(screen.queryByText('田')).not.toBeInTheDocument();
  });

  it('alt を省略すると名前が代替テキストになる', () => {
    render(<Avatar avatarUrl="/avatar.jpg" name="田中太郎" />);
    expect(screen.getByRole('img')).toHaveAttribute('alt', '田中太郎');
  });

  // alt="" は「装飾画像なので読み上げない」という指定。名前で上書きしてはいけない。
  it('空の alt を渡したら名前で上書きしない', () => {
    const { container } = render(<Avatar avatarUrl="/avatar.jpg" alt="" name="田中太郎" />);
    expect(container.querySelector('img')).toHaveAttribute('alt', '');
  });

  it('xsサイズが適用される', () => {
    const { container } = render(<Avatar name="T" size="xs" />);
    expect(container.querySelector('.w-6')).toBeInTheDocument();
  });

  it('xsサイズでもアイコンのフォールバックが潰れない', () => {
    const { container } = render(<Avatar size="xs" />);
    expect(container.querySelector('.lucide-user')).toHaveClass('w-3', 'h-3');
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
