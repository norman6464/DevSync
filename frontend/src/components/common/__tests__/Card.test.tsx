import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import Card from '../Card';
import { Book } from 'lucide-react';

describe('Card', () => {
  it('カードが表示される', () => {
    const { container } = render(<Card>カードコンテンツ</Card>);

    const card = container.querySelector('.bg-gray-900');
    expect(card).toBeInTheDocument();
  });

  it('子要素が表示される', () => {
    render(<Card>カードコンテンツ</Card>);

    expect(screen.getByText('カードコンテンツ')).toBeInTheDocument();
  });

  it('ヘッダーが表示される', () => {
    render(
      <Card>
        <Card.Header>ヘッダー</Card.Header>
        <Card.Body>ボディ</Card.Body>
      </Card>
    );

    expect(screen.getByText('ヘッダー')).toBeInTheDocument();
  });

  it('ボディが表示される', () => {
    render(
      <Card>
        <Card.Body>ボディコンテンツ</Card.Body>
      </Card>
    );

    expect(screen.getByText('ボディコンテンツ')).toBeInTheDocument();
  });

  it('フッターが表示される', () => {
    render(
      <Card>
        <Card.Body>ボディ</Card.Body>
        <Card.Footer>フッター</Card.Footer>
      </Card>
    );

    expect(screen.getByText('フッター')).toBeInTheDocument();
  });

  it('タイトルが表示される', () => {
    render(
      <Card>
        <Card.Header>
          <Card.Title>カードタイトル</Card.Title>
        </Card.Header>
      </Card>
    );

    expect(screen.getByText('カードタイトル')).toBeInTheDocument();
  });

  it('アイコン付きタイトルが表示される', () => {
    const { container } = render(
      <Card>
        <Card.Header>
          <Card.Title icon={<Book />}>タイトル</Card.Title>
        </Card.Header>
      </Card>
    );

    expect(screen.getByText('タイトル')).toBeInTheDocument();
    const icon = container.querySelector('svg');
    expect(icon).toBeInTheDocument();
  });

  it('クリック可能なカードが表示される', () => {
    const handleClick = vi.fn();

    render(<Card onClick={handleClick}>クリック可能</Card>);

    const card = screen.getByText('クリック可能');
    fireEvent.click(card);

    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('ホバーエフェクトがある', () => {
    const { container } = render(<Card hoverable>ホバー可能</Card>);

    const card = container.querySelector('.hover\\:border-gray-700');
    expect(card).toBeInTheDocument();
  });

  it('ボーダーが表示される', () => {
    const { container } = render(<Card>ボーダー</Card>);

    const card = container.querySelector('.border');
    expect(card).toBeInTheDocument();
  });

  it('角が丸くなっている', () => {
    const { container } = render(<Card>丸角</Card>);

    const card = container.querySelector('.rounded-lg');
    expect(card).toBeInTheDocument();
  });

  it('パディングが設定されている', () => {
    const { container } = render(<Card>パディング</Card>);

    const card = container.querySelector('.p-6');
    expect(card).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<Card className="custom-class">カスタム</Card>);

    const card = container.querySelector('.custom-class');
    expect(card).toBeInTheDocument();
  });

  it('ヘッダーとボディが分離されている', () => {
    render(
      <Card>
        <Card.Header>ヘッダー</Card.Header>
        <Card.Body>ボディ</Card.Body>
      </Card>
    );

    expect(screen.getByText('ヘッダー')).toBeInTheDocument();
    expect(screen.getByText('ボディ')).toBeInTheDocument();
  });

  it('フッターにボーダーがある', () => {
    render(
      <Card>
        <Card.Body>ボディ</Card.Body>
        <Card.Footer>フッター</Card.Footer>
      </Card>
    );

    const footer = screen.getByText('フッター');
    expect(footer).toHaveClass('border-t');
  });

  it('タイトルが太字で表示される', () => {
    render(
      <Card>
        <Card.Header>
          <Card.Title>太字タイトル</Card.Title>
        </Card.Header>
      </Card>
    );

    const title = screen.getByText('太字タイトル');
    expect(title).toHaveClass('font-semibold');
  });

  it('クリック可能なカードにカーソルポインターがある', () => {
    const { container } = render(<Card onClick={() => {}}>クリック</Card>);

    const card = container.querySelector('.cursor-pointer');
    expect(card).toBeInTheDocument();
  });

  it('複数のカードが独立して表示される', () => {
    render(
      <>
        <Card>カード1</Card>
        <Card>カード2</Card>
        <Card>カード3</Card>
      </>
    );

    expect(screen.getByText('カード1')).toBeInTheDocument();
    expect(screen.getByText('カード2')).toBeInTheDocument();
    expect(screen.getByText('カード3')).toBeInTheDocument();
  });
});
