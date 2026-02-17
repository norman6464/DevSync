import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { I18nextProvider } from 'react-i18next';
import i18n from '../../../i18n';
import ConfirmDialog from '../ConfirmDialog';

const renderWithI18n = (component: React.ReactElement) => {
  return render(
    <I18nextProvider i18n={i18n}>{component}</I18nextProvider>
  );
};

describe('ConfirmDialog', () => {
  const defaultProps = {
    isOpen: true,
    title: 'テストタイトル',
    message: 'テストメッセージ',
    onConfirm: vi.fn(),
    onCancel: vi.fn(),
  };

  it('isOpen=trueのときダイアログが表示される', () => {
    renderWithI18n(<ConfirmDialog {...defaultProps} />);
    expect(screen.getByText('テストタイトル')).toBeInTheDocument();
    expect(screen.getByText('テストメッセージ')).toBeInTheDocument();
  });

  it('isOpen=falseのときダイアログが表示されない', () => {
    renderWithI18n(<ConfirmDialog {...defaultProps} isOpen={false} />);
    expect(screen.queryByText('テストタイトル')).not.toBeInTheDocument();
  });

  it('確認ボタンをクリックするとonConfirmが呼ばれる', async () => {
    const user = userEvent.setup();
    const onConfirm = vi.fn();
    renderWithI18n(<ConfirmDialog {...defaultProps} onConfirm={onConfirm} />);
    await user.click(screen.getByRole('button', { name: '確認' }));
    expect(onConfirm).toHaveBeenCalledTimes(1);
  });

  it('キャンセルボタンをクリックするとonCancelが呼ばれる', async () => {
    const user = userEvent.setup();
    const onCancel = vi.fn();
    renderWithI18n(<ConfirmDialog {...defaultProps} onCancel={onCancel} />);
    await user.click(screen.getByRole('button', { name: 'キャンセル' }));
    expect(onCancel).toHaveBeenCalledTimes(1);
  });

  it('Escapeキーでonancelが呼ばれる', async () => {
    const user = userEvent.setup();
    const onCancel = vi.fn();
    renderWithI18n(<ConfirmDialog {...defaultProps} onCancel={onCancel} />);
    await user.keyboard('{Escape}');
    expect(onCancel).toHaveBeenCalledTimes(1);
  });

  it('variant="danger"のとき確認ボタンが赤色スタイルになる', () => {
    renderWithI18n(<ConfirmDialog {...defaultProps} variant="danger" />);
    const confirmBtn = screen.getByRole('button', { name: '確認' });
    expect(confirmBtn).toHaveClass('bg-red-600');
  });

  it('variant="warning"のとき確認ボタンが黄色スタイルになる', () => {
    renderWithI18n(<ConfirmDialog {...defaultProps} variant="warning" />);
    const confirmBtn = screen.getByRole('button', { name: '確認' });
    expect(confirmBtn).toHaveClass('bg-yellow-600');
  });

  it('カスタムconfirmText/cancelTextが表示される', () => {
    renderWithI18n(
      <ConfirmDialog
        {...defaultProps}
        confirmText="削除する"
        cancelText="やめる"
      />
    );
    expect(screen.getByRole('button', { name: '削除する' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'やめる' })).toBeInTheDocument();
  });

  it('オーバーレイをクリックするとonCancelが呼ばれる', async () => {
    const user = userEvent.setup();
    const onCancel = vi.fn();
    renderWithI18n(<ConfirmDialog {...defaultProps} onCancel={onCancel} />);
    const overlay = screen.getByTestId('confirm-dialog-overlay');
    await user.click(overlay);
    expect(onCancel).toHaveBeenCalledTimes(1);
  });

  it('ダイアログ本体をクリックしてもonCancelは呼ばれない', async () => {
    const user = userEvent.setup();
    const onCancel = vi.fn();
    renderWithI18n(<ConfirmDialog {...defaultProps} onCancel={onCancel} />);
    const dialog = screen.getByRole('dialog');
    await user.click(dialog);
    expect(onCancel).not.toHaveBeenCalled();
  });
});
