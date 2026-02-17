import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import SearchInput from '../SearchInput';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

describe('SearchInput', () => {
  it('入力フィールドが表示される', () => {
    render(<SearchInput value="" onChange={() => {}} placeholder="検索..." />);
    expect(screen.getByPlaceholderText('検索...')).toBeInTheDocument();
  });

  it('値の変更でonChangeが呼ばれる', () => {
    const onChange = vi.fn();
    render(<SearchInput value="" onChange={onChange} placeholder="検索..." />);
    fireEvent.change(screen.getByPlaceholderText('検索...'), { target: { value: 'test' } });
    expect(onChange).toHaveBeenCalledWith('test');
  });

  it('EnterキーでonSearchが呼ばれる', () => {
    const onSearch = vi.fn();
    render(<SearchInput value="query" onChange={() => {}} onSearch={onSearch} placeholder="検索..." />);
    fireEvent.keyDown(screen.getByPlaceholderText('検索...'), { key: 'Enter' });
    expect(onSearch).toHaveBeenCalledOnce();
  });

  it('onSearchが未指定の場合Enterキーでエラーにならない', () => {
    render(<SearchInput value="query" onChange={() => {}} placeholder="検索..." />);
    expect(() => {
      fireEvent.keyDown(screen.getByPlaceholderText('検索...'), { key: 'Enter' });
    }).not.toThrow();
  });

  it('showButton=trueで検索ボタンが表示される', () => {
    const onSearch = vi.fn();
    render(<SearchInput value="" onChange={() => {}} onSearch={onSearch} showButton />);
    const button = screen.getByText('common.search');
    expect(button).toBeInTheDocument();
    fireEvent.click(button);
    expect(onSearch).toHaveBeenCalledOnce();
  });

  it('showButton=falseの場合ボタンが表示されない', () => {
    render(<SearchInput value="" onChange={() => {}} onSearch={() => {}} />);
    expect(screen.queryByText('common.search')).toBeNull();
  });

  it('showButton=trueでもonSearch未指定の場合ボタンが表示されない', () => {
    render(<SearchInput value="" onChange={() => {}} showButton />);
    expect(screen.queryByText('common.search')).toBeNull();
  });

  it('検索アイコンが表示される', () => {
    const { container } = render(<SearchInput value="" onChange={() => {}} />);
    expect(container.querySelector('svg')).toBeInTheDocument();
  });
});
