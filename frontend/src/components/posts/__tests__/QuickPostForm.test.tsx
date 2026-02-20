import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import QuickPostForm from '../QuickPostForm';

// Mock authStore
vi.mock('../../../store/authStore', () => ({
  useAuthStore: (selector: (s: unknown) => unknown) =>
    selector({ user: { id: 1, name: 'TestUser', avatar_url: '' } }),
}));

// Mock useAutoSave
vi.mock('../../../hooks/useAutoSave', () => ({
  useAutoSave: () => ({
    saveStatus: 'idle',
    lastSaved: null,
    clearSaved: vi.fn(),
    getSaved: () => null,
  }),
}));

describe('QuickPostForm', () => {
  const onSubmit = vi.fn().mockResolvedValue(undefined);

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('テキストエリアが表示される', () => {
    render(<QuickPostForm onSubmit={onSubmit} />);
    expect(screen.getByRole('textbox')).toBeInTheDocument();
  });

  it('テキスト入力が可能', () => {
    render(<QuickPostForm onSubmit={onSubmit} />);
    const textarea = screen.getByRole('textbox');
    fireEvent.change(textarea, { target: { value: 'テスト投稿' } });
    expect(textarea).toHaveValue('テスト投稿');
  });

  it('フォーカス時に送信ボタンが表示される', () => {
    render(<QuickPostForm onSubmit={onSubmit} />);
    fireEvent.focus(screen.getByRole('textbox'));
    expect(screen.getByText('投稿')).toBeInTheDocument();
  });

  it('フォーカス時にショートカットヒントが表示される', () => {
    render(<QuickPostForm onSubmit={onSubmit} />);
    fireEvent.focus(screen.getByRole('textbox'));
    expect(screen.getByText(/Ctrl\+Enter/)).toBeInTheDocument();
  });

  it('フォーカス時に下書き保存ボタンが表示される', () => {
    render(<QuickPostForm onSubmit={onSubmit} />);
    fireEvent.focus(screen.getByRole('textbox'));
    expect(screen.getByText('下書き保存')).toBeInTheDocument();
  });

  it('Ctrl+Enterでコンテンツがあれば送信される', async () => {
    render(<QuickPostForm onSubmit={onSubmit} />);
    const textarea = screen.getByRole('textbox');
    fireEvent.change(textarea, { target: { value: 'テスト投稿内容' } });
    fireEvent.keyDown(textarea, { key: 'Enter', ctrlKey: true });
    expect(onSubmit).toHaveBeenCalledOnce();
  });

  it('Ctrl+Enterでコンテンツが空なら送信されない', () => {
    render(<QuickPostForm onSubmit={onSubmit} />);
    const textarea = screen.getByRole('textbox');
    fireEvent.keyDown(textarea, { key: 'Enter', ctrlKey: true });
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('Meta+Enterでも送信される（Mac対応）', async () => {
    render(<QuickPostForm onSubmit={onSubmit} />);
    const textarea = screen.getByRole('textbox');
    fireEvent.change(textarea, { target: { value: 'Mac投稿' } });
    fireEvent.keyDown(textarea, { key: 'Enter', metaKey: true });
    expect(onSubmit).toHaveBeenCalledOnce();
  });

  it('通常のEnterでは送信されない', () => {
    render(<QuickPostForm onSubmit={onSubmit} />);
    const textarea = screen.getByRole('textbox');
    fireEvent.change(textarea, { target: { value: 'テスト' } });
    fireEvent.keyDown(textarea, { key: 'Enter' });
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('maxLengthが5000に設定されている', () => {
    render(<QuickPostForm onSubmit={onSubmit} />);
    expect(screen.getByRole('textbox')).toHaveAttribute('maxLength', '5000');
  });
});
