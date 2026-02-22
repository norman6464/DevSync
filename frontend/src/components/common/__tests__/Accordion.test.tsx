import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Accordion from '../Accordion';

const items = [
  { id: '1', title: 'セクション1', content: 'コンテンツ1' },
  { id: '2', title: 'セクション2', content: 'コンテンツ2' },
  { id: '3', title: 'セクション3', content: 'コンテンツ3' },
];

describe('Accordion', () => {
  it('全てのセクションタイトルが表示される', () => {
    render(<Accordion items={items} />);

    expect(screen.getByText('セクション1')).toBeInTheDocument();
    expect(screen.getByText('セクション2')).toBeInTheDocument();
    expect(screen.getByText('セクション3')).toBeInTheDocument();
  });

  it('初期状態でコンテンツが非表示', () => {
    render(<Accordion items={items} />);

    expect(screen.queryByText('コンテンツ1')).not.toBeInTheDocument();
    expect(screen.queryByText('コンテンツ2')).not.toBeInTheDocument();
  });

  it('クリックでセクションが展開される', async () => {
    const user = userEvent.setup();
    render(<Accordion items={items} />);

    await user.click(screen.getByText('セクション1'));

    expect(screen.getByText('コンテンツ1')).toBeInTheDocument();
  });

  it('展開中のセクションをクリックで閉じる', async () => {
    const user = userEvent.setup();
    render(<Accordion items={items} />);

    await user.click(screen.getByText('セクション1'));
    expect(screen.getByText('コンテンツ1')).toBeInTheDocument();

    await user.click(screen.getByText('セクション1'));
    expect(screen.queryByText('コンテンツ1')).not.toBeInTheDocument();
  });

  it('単一展開モードで他のセクションが閉じる', async () => {
    const user = userEvent.setup();
    render(<Accordion items={items} single />);

    await user.click(screen.getByText('セクション1'));
    expect(screen.getByText('コンテンツ1')).toBeInTheDocument();

    await user.click(screen.getByText('セクション2'));
    expect(screen.queryByText('コンテンツ1')).not.toBeInTheDocument();
    expect(screen.getByText('コンテンツ2')).toBeInTheDocument();
  });

  it('複数展開モードで複数セクションが開く', async () => {
    const user = userEvent.setup();
    render(<Accordion items={items} />);

    await user.click(screen.getByText('セクション1'));
    await user.click(screen.getByText('セクション2'));

    expect(screen.getByText('コンテンツ1')).toBeInTheDocument();
    expect(screen.getByText('コンテンツ2')).toBeInTheDocument();
  });

  it('デフォルト展開状態が設定できる', () => {
    render(<Accordion items={items} defaultOpenIds={['1', '3']} />);

    expect(screen.getByText('コンテンツ1')).toBeInTheDocument();
    expect(screen.queryByText('コンテンツ2')).not.toBeInTheDocument();
    expect(screen.getByText('コンテンツ3')).toBeInTheDocument();
  });

  it('onChangeコールバックが呼ばれる', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<Accordion items={items} onChange={onChange} />);

    await user.click(screen.getByText('セクション1'));

    expect(onChange).toHaveBeenCalledWith(['1']);
  });

  it('展開アイコンが表示される', () => {
    const { container } = render(<Accordion items={items} />);

    const icons = container.querySelectorAll('.lucide-chevron-down');
    expect(icons.length).toBe(3);
  });

  it('展開時にアイコンが回転する', async () => {
    const user = userEvent.setup();
    const { container } = render(<Accordion items={items} />);

    await user.click(screen.getByText('セクション1'));

    const rotatedIcon = container.querySelector('.rotate-180');
    expect(rotatedIcon).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<Accordion items={items} className="custom-class" />);

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('ReactNodeをコンテンツに使用できる', async () => {
    const user = userEvent.setup();
    const richItems = [
      { id: '1', title: 'リッチ', content: <span data-testid="rich">リッチコンテンツ</span> },
    ];
    render(<Accordion items={richItems} />);

    await user.click(screen.getByText('リッチ'));

    expect(screen.getByTestId('rich')).toBeInTheDocument();
  });
});
