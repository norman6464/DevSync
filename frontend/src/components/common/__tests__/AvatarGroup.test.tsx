import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import AvatarGroup from '../AvatarGroup';

const users = [
  { name: 'Alice', image: 'https://example.com/alice.jpg' },
  { name: 'Bob', image: 'https://example.com/bob.jpg' },
  { name: 'Charlie', image: 'https://example.com/charlie.jpg' },
  { name: 'Diana', image: 'https://example.com/diana.jpg' },
  { name: 'Eve', image: 'https://example.com/eve.jpg' },
];

describe('AvatarGroup', () => {
  it('全てのアバターが表示される', () => {
    render(<AvatarGroup users={users.slice(0, 3)} />);

    const images = screen.getAllByRole('img');
    expect(images.length).toBe(3);
  });

  it('最大表示数を超えると+N表示', () => {
    render(<AvatarGroup users={users} max={3} />);

    expect(screen.getByText('+2')).toBeInTheDocument();
  });

  it('最大表示数以下では+N非表示', () => {
    render(<AvatarGroup users={users.slice(0, 3)} max={5} />);

    expect(screen.queryByText(/\+/)).not.toBeInTheDocument();
  });

  it('画像のalt属性にユーザー名が設定される', () => {
    render(<AvatarGroup users={users.slice(0, 2)} />);

    expect(screen.getByAltText('Alice')).toBeInTheDocument();
    expect(screen.getByAltText('Bob')).toBeInTheDocument();
  });

  it('画像がない場合はイニシャルが表示される', () => {
    const usersNoImage = [
      { name: 'Alice' },
      { name: 'Bob' },
    ];
    render(<AvatarGroup users={usersNoImage} />);

    expect(screen.getByText('A')).toBeInTheDocument();
    expect(screen.getByText('B')).toBeInTheDocument();
  });

  it('smサイズが適用される', () => {
    const { container } = render(<AvatarGroup users={users.slice(0, 1)} size="sm" />);

    expect(container.querySelector('.w-8')).toBeInTheDocument();
  });

  it('mdサイズが適用される', () => {
    const { container } = render(<AvatarGroup users={users.slice(0, 1)} size="md" />);

    expect(container.querySelector('.w-10')).toBeInTheDocument();
  });

  it('lgサイズが適用される', () => {
    const { container } = render(<AvatarGroup users={users.slice(0, 1)} size="lg" />);

    expect(container.querySelector('.w-12')).toBeInTheDocument();
  });

  it('アバターが重なって表示される', () => {
    const { container } = render(<AvatarGroup users={users.slice(0, 3)} />);

    const overlapping = container.querySelectorAll('.-ml-2');
    expect(overlapping.length).toBe(2);
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<AvatarGroup users={users.slice(0, 1)} className="custom-class" />);

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('空の配列で何も表示されない', () => {
    const { container } = render(<AvatarGroup users={[]} />);

    expect(container.querySelector('img')).not.toBeInTheDocument();
  });
});
