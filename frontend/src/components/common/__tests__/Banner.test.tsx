import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Banner from '../Banner';

describe('Banner', () => {
  it('メッセージが表示される', () => {
    render(<Banner message="お知らせです" />);
    expect(screen.getByText('お知らせです')).toBeInTheDocument();
  });

  it('infoバリアントのスタイルが適用される', () => {
    const { container } = render(<Banner message="情報" variant="info" />);
    expect(container.querySelector('.bg-blue-900\\/30')).toBeInTheDocument();
  });

  it('successバリアントのスタイルが適用される', () => {
    const { container } = render(<Banner message="成功" variant="success" />);
    expect(container.querySelector('.bg-green-900\\/30')).toBeInTheDocument();
  });

  it('warningバリアントのスタイルが適用される', () => {
    const { container } = render(<Banner message="警告" variant="warning" />);
    expect(container.querySelector('.bg-yellow-900\\/30')).toBeInTheDocument();
  });

  it('errorバリアントのスタイルが適用される', () => {
    const { container } = render(<Banner message="エラー" variant="error" />);
    expect(container.querySelector('.bg-red-900\\/30')).toBeInTheDocument();
  });

  it('閉じるボタンでonCloseが呼ばれる', async () => {
    const onClose = vi.fn();
    const user = userEvent.setup();
    render(<Banner message="テスト" onClose={onClose} />);
    await user.click(screen.getByLabelText('閉じる'));
    expect(onClose).toHaveBeenCalled();
  });

  it('onCloseがない場合閉じるボタンが非表示', () => {
    render(<Banner message="テスト" />);
    expect(screen.queryByLabelText('閉じる')).not.toBeInTheDocument();
  });

  it('タイトルが表示される', () => {
    render(<Banner message="内容" title="重要なお知らせ" />);
    expect(screen.getByText('重要なお知らせ')).toBeInTheDocument();
  });

  it('アクションボタンが表示される', async () => {
    const onAction = vi.fn();
    const user = userEvent.setup();
    render(<Banner message="内容" actionLabel="詳細を見る" onAction={onAction} />);
    await user.click(screen.getByText('詳細を見る'));
    expect(onAction).toHaveBeenCalled();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<Banner message="テスト" className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('デフォルトバリアントはinfo', () => {
    const { container } = render(<Banner message="テスト" />);
    expect(container.querySelector('.bg-blue-900\\/30')).toBeInTheDocument();
  });
});
