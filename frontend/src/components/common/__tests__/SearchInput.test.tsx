import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import SearchInput from '../SearchInput';

describe('SearchInput', () => {
  it('検索入力欄が表示される', () => {
    render(<SearchInput value="" onChange={() => {}} />);

    expect(screen.getByRole('textbox')).toBeInTheDocument();
  });

  it('検索アイコンが表示される', () => {
    const { container } = render(<SearchInput value="" onChange={() => {}} />);

    expect(container.querySelector('.lucide-search')).toBeInTheDocument();
  });

  it('プレースホルダーが表示される', () => {
    render(<SearchInput value="" onChange={() => {}} placeholder="検索..." />);

    expect(screen.getByPlaceholderText('検索...')).toBeInTheDocument();
  });

  it('値が入力される', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<SearchInput value="" onChange={onChange} />);

    await user.type(screen.getByRole('textbox'), 'a');

    expect(onChange).toHaveBeenCalledWith('a');
  });

  it('クリアボタンで値がリセットされる', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<SearchInput value="テスト" onChange={onChange} />);

    const clearButton = screen.getByRole('button');
    await user.click(clearButton);

    expect(onChange).toHaveBeenCalledWith('');
  });

  it('値が空の場合はクリアボタンが表示されない', () => {
    render(<SearchInput value="" onChange={() => {}} />);

    expect(screen.queryByRole('button')).not.toBeInTheDocument();
  });

  it('ローディング状態が表示される', () => {
    const { container } = render(<SearchInput value="" onChange={() => {}} loading />);

    expect(container.querySelector('.animate-spin')).toBeInTheDocument();
  });

  it('ローディング中は検索アイコンが非表示', () => {
    const { container } = render(<SearchInput value="" onChange={() => {}} loading />);

    expect(container.querySelector('.lucide-search')).not.toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<SearchInput value="" onChange={() => {}} className="custom-class" />);

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('無効状態が適用される', () => {
    render(<SearchInput value="" onChange={() => {}} disabled />);

    expect(screen.getByRole('textbox')).toBeDisabled();
  });

  it('showButton で検索ボタンが表示される', () => {
    render(<SearchInput value="" onChange={() => {}} onSearch={() => {}} showButton />);

    expect(screen.getByRole('button', { name: '検索' })).toBeInTheDocument();
  });

  it('showButton が無ければ検索ボタンを表示しない', () => {
    render(<SearchInput value="" onChange={() => {}} />);

    expect(screen.queryByRole('button', { name: '検索' })).not.toBeInTheDocument();
  });

  it('検索ボタンクリックで onSearch が呼ばれる', () => {
    const onSearch = vi.fn();
    render(<SearchInput value="react" onChange={() => {}} onSearch={onSearch} showButton />);

    fireEvent.click(screen.getByRole('button', { name: '検索' }));

    expect(onSearch).toHaveBeenCalledTimes(1);
  });

  it('Enter キーで onSearch が呼ばれる', () => {
    const onSearch = vi.fn();
    render(<SearchInput value="react" onChange={() => {}} onSearch={onSearch} />);

    fireEvent.keyDown(screen.getByRole('textbox'), { key: 'Enter' });

    expect(onSearch).toHaveBeenCalledTimes(1);
  });

  it('onSearch が無ければ Enter で何も起きない', () => {
    render(<SearchInput value="react" onChange={() => {}} />);

    fireEvent.keyDown(screen.getByRole('textbox'), { key: 'Enter' });
    // クラッシュしないことの確認
    expect(screen.getByRole('textbox')).toBeInTheDocument();
  });
});
