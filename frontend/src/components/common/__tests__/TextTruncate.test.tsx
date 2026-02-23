import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import TextTruncate from '../TextTruncate';

const longText = 'これは非常に長いテキストです。省略表示のテストに使用します。このテキストは十分に長い必要があります。さらに長くするために追加のテキストを入れます。まだまだ続きます。';

describe('TextTruncate', () => {
  it('テキストが表示される', () => {
    render(<TextTruncate text="短いテキスト" />);

    expect(screen.getByText('短いテキスト')).toBeInTheDocument();
  });

  it('文字数制限でテキストが省略される', () => {
    render(<TextTruncate text={longText} maxLength={20} />);

    expect(screen.getByText(/これは非常に長いテキストです。省略/)).toBeInTheDocument();
    expect(screen.getByText('...')).toBeInTheDocument();
  });

  it('展開ボタンが表示される', () => {
    render(<TextTruncate text={longText} maxLength={20} />);

    expect(screen.getByText('もっと見る')).toBeInTheDocument();
  });

  it('展開ボタンクリックで全文表示', async () => {
    const user = userEvent.setup();
    render(<TextTruncate text={longText} maxLength={20} />);

    await user.click(screen.getByText('もっと見る'));

    expect(screen.getByText(longText)).toBeInTheDocument();
  });

  it('折りたたみボタンが表示される', async () => {
    const user = userEvent.setup();
    render(<TextTruncate text={longText} maxLength={20} />);

    await user.click(screen.getByText('もっと見る'));

    expect(screen.getByText('閉じる')).toBeInTheDocument();
  });

  it('折りたたみボタンクリックで省略表示に戻る', async () => {
    const user = userEvent.setup();
    render(<TextTruncate text={longText} maxLength={20} />);

    await user.click(screen.getByText('もっと見る'));
    await user.click(screen.getByText('閉じる'));

    expect(screen.getByText('もっと見る')).toBeInTheDocument();
  });

  it('短いテキストでは展開ボタンが表示されない', () => {
    render(<TextTruncate text="短い" maxLength={100} />);

    expect(screen.queryByText('もっと見る')).not.toBeInTheDocument();
  });

  it('カスタム展開ラベルが使用される', () => {
    render(<TextTruncate text={longText} maxLength={20} expandLabel="続きを読む" />);

    expect(screen.getByText('続きを読む')).toBeInTheDocument();
  });

  it('カスタム折りたたみラベルが使用される', async () => {
    const user = userEvent.setup();
    render(<TextTruncate text={longText} maxLength={20} collapseLabel="少なく表示" />);

    await user.click(screen.getByText('もっと見る'));

    expect(screen.getByText('少なく表示')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<TextTruncate text="テスト" className="custom-class" />);

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
