import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import ConfirmDialog from '../ConfirmDialog';

describe('ConfirmDialog', () => {
  it('タイトルが表示される', () => {
    render(<ConfirmDialog isOpen title="確認ダイアログ" message="削除しますか？" onConfirm={() => {}} onCancel={() => {}} />);

    expect(screen.getByText('確認ダイアログ')).toBeInTheDocument();
  });

  it('メッセージが表示される', () => {
    render(<ConfirmDialog isOpen title="確認ダイアログ" message="削除しますか？" onConfirm={() => {}} onCancel={() => {}} />);

    expect(screen.getByText('削除しますか？')).toBeInTheDocument();
  });

  it('確認ボタンが表示される', () => {
    render(<ConfirmDialog isOpen title="削除の確認" message="削除？" onConfirm={() => {}} onCancel={() => {}} confirmLabel="実行" />);

    expect(screen.getByText('実行')).toBeInTheDocument();
  });

  it('キャンセルボタンが表示される', () => {
    render(<ConfirmDialog isOpen title="確認" message="削除？" onConfirm={() => {}} onCancel={() => {}} />);

    expect(screen.getByText('キャンセル')).toBeInTheDocument();
  });

  it('確認ボタンクリックでコールバック呼ばれる', async () => {
    const onConfirm = vi.fn();
    const user = userEvent.setup();
    render(<ConfirmDialog isOpen title="確認" message="削除？" onConfirm={onConfirm} onCancel={() => {}} confirmLabel="削除する" />);

    await user.click(screen.getByText('削除する'));

    expect(onConfirm).toHaveBeenCalledTimes(1);
  });

  it('キャンセルボタンクリックでコールバック呼ばれる', async () => {
    const onCancel = vi.fn();
    const user = userEvent.setup();
    render(<ConfirmDialog isOpen title="確認" message="削除？" onConfirm={() => {}} onCancel={onCancel} />);

    await user.click(screen.getByText('キャンセル'));

    expect(onCancel).toHaveBeenCalledTimes(1);
  });

  it('isOpenがfalseの場合は非表示', () => {
    render(<ConfirmDialog isOpen={false} title="確認" message="削除？" onConfirm={() => {}} onCancel={() => {}} />);

    expect(screen.queryByText('削除？')).not.toBeInTheDocument();
  });

  it('dangerバリアントで赤いボタン', () => {
    const { container } = render(
      <ConfirmDialog isOpen title="確認" message="削除？" onConfirm={() => {}} onCancel={() => {}} variant="danger" />
    );

    expect(container.querySelector('.bg-red-600')).toBeInTheDocument();
  });

  it('ローディング状態が表示される', () => {
    const { container } = render(
      <ConfirmDialog isOpen title="確認" message="削除？" onConfirm={() => {}} onCancel={() => {}} loading />
    );

    expect(container.querySelector('.animate-spin')).toBeInTheDocument();
  });

  it('ローディング中はボタンが無効', () => {
    render(
      <ConfirmDialog isOpen title="確認" message="削除？" onConfirm={() => {}} onCancel={() => {}} loading confirmLabel="OK" />
    );

    const buttons = screen.getAllByRole('button');
    const confirmBtn = buttons.find(b => b.textContent?.includes('OK'));
    expect(confirmBtn).toBeDisabled();
  });

  it('カスタム確認ラベルが使用される', () => {
    render(
      <ConfirmDialog isOpen title="確認" message="削除？" onConfirm={() => {}} onCancel={() => {}} confirmLabel="削除する" />
    );

    expect(screen.getByText('削除する')).toBeInTheDocument();
  });
});
