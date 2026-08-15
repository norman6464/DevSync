import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import ConfirmDialog from '../ConfirmDialog';
import { useConfirm } from '../../../hooks/useConfirm';

/** useConfirm と ConfirmDialog を実際の使い方どおりに結線した確認画面。 */
function Harness({
  options,
  onResult,
}: {
  options: Parameters<ReturnType<typeof useConfirm>['confirm']>[0];
  onResult?: (value: boolean) => void;
}) {
  const { confirm, dialogProps } = useConfirm();

  return (
    <>
      <button type="button" onClick={() => confirm(options).then((v) => onResult?.(v))}>
        開く
      </button>
      <ConfirmDialog {...dialogProps} />
    </>
  );
}

const open = async (ui: React.ReactElement) => {
  const user = userEvent.setup();
  render(ui);
  await user.click(screen.getByRole('button', { name: '開く' }));
  return user;
};

describe('useConfirm と ConfirmDialog の結線', () => {
  // フックが渡す名前とコンポーネントが受ける名前がずれると、
  // 指定した文言が黙って捨てられ既定のラベルになる。
  it('指定した確認ボタンの文言が表示される', async () => {
    await open(
      <Harness options={{ title: '削除の確認', message: '本当に削除しますか？', confirmText: '削除する' }} />,
    );

    expect(screen.getByRole('button', { name: '削除する' })).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: '確認' })).not.toBeInTheDocument();
  });

  it('指定した取消ボタンの文言が表示される', async () => {
    await open(
      <Harness options={{ title: '削除の確認', message: '本当に削除しますか？', cancelText: 'やめる' }} />,
    );

    expect(screen.getByRole('button', { name: 'やめる' })).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'キャンセル' })).not.toBeInTheDocument();
  });

  it('文言を指定しなければ既定のラベルになる', async () => {
    await open(<Harness options={{ title: '確認', message: 'よろしいですか？' }} />);

    expect(screen.getByRole('button', { name: '確認' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'キャンセル' })).toBeInTheDocument();
  });

  it('タイトルと本文がダイアログの名前と説明になる', async () => {
    await open(<Harness options={{ title: '削除の確認', message: '本当に削除しますか？' }} />);

    const dialog = screen.getByRole('alertdialog', { name: '削除の確認' });
    expect(dialog).toHaveAccessibleDescription('本当に削除しますか？');
  });

  it('確認すると true で解決しダイアログが閉じる', async () => {
    const results: boolean[] = [];
    const user = await open(
      <Harness
        options={{ title: '削除の確認', message: '本当に削除しますか？', confirmText: '削除する' }}
        onResult={(v) => results.push(v)}
      />,
    );

    await user.click(screen.getByRole('button', { name: '削除する' }));

    expect(results).toEqual([true]);
    expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument();
  });

  it('取り消すと false で解決しダイアログが閉じる', async () => {
    const results: boolean[] = [];
    const user = await open(
      <Harness
        options={{ title: '削除の確認', message: '本当に削除しますか？', cancelText: 'やめる' }}
        onResult={(v) => results.push(v)}
      />,
    );

    await user.click(screen.getByRole('button', { name: 'やめる' }));

    expect(results).toEqual([false]);
    expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument();
  });
});
