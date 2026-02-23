import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import CommandPalette from '../CommandPalette';

const commands = [
  { id: '1', label: 'ダッシュボード', category: 'ページ' },
  { id: '2', label: 'プロフィール', category: 'ページ' },
  { id: '3', label: 'ダークモード切替', category: '設定' },
  { id: '4', label: 'ログアウト', category: 'アクション' },
];

describe('CommandPalette', () => {
  it('開いた状態でコマンドが表示される', () => {
    render(<CommandPalette open commands={commands} onSelect={() => {}} onClose={() => {}} />);
    expect(screen.getByText('ダッシュボード')).toBeInTheDocument();
    expect(screen.getByText('ログアウト')).toBeInTheDocument();
  });

  it('閉じた状態で非表示', () => {
    render(<CommandPalette open={false} commands={commands} onSelect={() => {}} onClose={() => {}} />);
    expect(screen.queryByText('ダッシュボード')).not.toBeInTheDocument();
  });

  it('検索でフィルタリングされる', async () => {
    const user = userEvent.setup();
    render(<CommandPalette open commands={commands} onSelect={() => {}} onClose={() => {}} />);
    await user.type(screen.getByPlaceholderText('コマンドを検索...'), 'ダッシュ');
    expect(screen.getByText('ダッシュボード')).toBeInTheDocument();
    expect(screen.queryByText('ログアウト')).not.toBeInTheDocument();
  });

  it('コマンドクリックでonSelectが呼ばれる', async () => {
    const onSelect = vi.fn();
    const user = userEvent.setup();
    render(<CommandPalette open commands={commands} onSelect={onSelect} onClose={() => {}} />);
    await user.click(screen.getByText('ダッシュボード'));
    expect(onSelect).toHaveBeenCalledWith(commands[0]);
  });

  it('オーバーレイクリックでonCloseが呼ばれる', async () => {
    const onClose = vi.fn();
    const user = userEvent.setup();
    render(<CommandPalette open commands={commands} onSelect={() => {}} onClose={onClose} />);
    await user.click(screen.getByTestId('palette-overlay'));
    expect(onClose).toHaveBeenCalled();
  });

  it('カテゴリが表示される', () => {
    render(<CommandPalette open commands={commands} onSelect={() => {}} onClose={() => {}} />);
    expect(screen.getByText('ページ')).toBeInTheDocument();
    expect(screen.getByText('設定')).toBeInTheDocument();
    expect(screen.getByText('アクション')).toBeInTheDocument();
  });

  it('検索結果が空の場合メッセージが表示される', async () => {
    const user = userEvent.setup();
    render(<CommandPalette open commands={commands} onSelect={() => {}} onClose={() => {}} />);
    await user.type(screen.getByPlaceholderText('コマンドを検索...'), 'zzzzz');
    expect(screen.getByText('コマンドが見つかりません')).toBeInTheDocument();
  });

  it('検索フィールドが自動フォーカスされる', () => {
    render(<CommandPalette open commands={commands} onSelect={() => {}} onClose={() => {}} />);
    expect(screen.getByPlaceholderText('コマンドを検索...')).toHaveFocus();
  });

  it('アイコン付きコマンドが表示される', () => {
    const cmds = [{ id: '1', label: 'テスト', category: 'カテゴリ', icon: '🔍' }];
    render(<CommandPalette open commands={cmds} onSelect={() => {}} onClose={() => {}} />);
    expect(screen.getByText('🔍')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    render(<CommandPalette open commands={commands} onSelect={() => {}} onClose={() => {}} className="custom-class" />);
    expect(document.querySelector('.custom-class')).toBeInTheDocument();
  });
});
