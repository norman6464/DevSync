import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import ProfileCard from '../ProfileCard';

describe('ProfileCard', () => {
  it('名前が表示される', () => {
    render(<ProfileCard name="田中太郎" />);
    expect(screen.getByText('田中太郎')).toBeInTheDocument();
  });

  it('肩書きが表示される', () => {
    render(<ProfileCard name="田中太郎" title="フロントエンドエンジニア" />);
    expect(screen.getByText('フロントエンドエンジニア')).toBeInTheDocument();
  });

  it('アバター画像が表示される', () => {
    render(<ProfileCard name="田中太郎" avatarUrl="/avatar.jpg" />);
    const img = screen.getByRole('img');
    expect(img).toHaveAttribute('src', '/avatar.jpg');
  });

  it('アバターがない場合イニシャルが表示される', () => {
    render(<ProfileCard name="田中太郎" />);
    expect(screen.getByText('田')).toBeInTheDocument();
  });

  it('スタッツが表示される', () => {
    const stats = [
      { label: '投稿', value: 42 },
      { label: 'フォロワー', value: 128 },
    ];
    render(<ProfileCard name="田中" stats={stats} />);
    expect(screen.getByText('42')).toBeInTheDocument();
    expect(screen.getByText('128')).toBeInTheDocument();
    expect(screen.getByText('投稿')).toBeInTheDocument();
  });

  it('フォローボタンが表示される', async () => {
    const onFollow = vi.fn();
    const user = userEvent.setup();
    render(<ProfileCard name="田中" onFollow={onFollow} />);
    await user.click(screen.getByText('フォロー'));
    expect(onFollow).toHaveBeenCalled();
  });

  it('フォロー済みの場合ボタンテキストが変わる', () => {
    render(<ProfileCard name="田中" onFollow={() => {}} isFollowing />);
    expect(screen.getByText('フォロー中')).toBeInTheDocument();
  });

  it('自己紹介が表示される', () => {
    render(<ProfileCard name="田中" bio="React大好きエンジニアです" />);
    expect(screen.getByText('React大好きエンジニアです')).toBeInTheDocument();
  });

  it('フォローボタンがない場合非表示', () => {
    render(<ProfileCard name="田中" />);
    expect(screen.queryByText('フォロー')).not.toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<ProfileCard name="田中" className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
