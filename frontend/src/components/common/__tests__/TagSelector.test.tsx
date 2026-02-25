import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import TagSelector from '../TagSelector';

const mockOnChange = vi.fn();

describe('TagSelector', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('選択済みタグが表示される', () => {
    render(<TagSelector value="#React,#TypeScript" onChange={mockOnChange} />);

    expect(screen.getByText('React')).toBeInTheDocument();
    expect(screen.getByText('TypeScript')).toBeInTheDocument();
  });

  it('タグ入力フィールドが表示される', () => {
    render(<TagSelector value="" onChange={mockOnChange} />);

    expect(screen.getByPlaceholderText(/タグを入力/)).toBeInTheDocument();
  });

  it('タグをクリックすると削除される', () => {
    render(<TagSelector value="#React,#TypeScript" onChange={mockOnChange} />);

    const reactTag = screen.getByText('React');
    fireEvent.click(reactTag);

    expect(mockOnChange).toHaveBeenCalledWith('#TypeScript');
  });

  it('新しいタグを入力してEnterキーで追加できる', () => {
    render(<TagSelector value="#React" onChange={mockOnChange} />);

    const input = screen.getByPlaceholderText(/タグを入力/);
    fireEvent.change(input, { target: { value: 'TypeScript' } });
    fireEvent.keyDown(input, { key: 'Enter' });

    expect(mockOnChange).toHaveBeenCalledWith('#React,#TypeScript');
  });

  it('候補タグが表示される', () => {
    const suggestions = ['JavaScript', 'TypeScript', 'Python'];
    render(<TagSelector value="" onChange={mockOnChange} suggestions={suggestions} />);

    suggestions.forEach((tag) => {
      expect(screen.getByText(tag)).toBeInTheDocument();
    });
  });

  it('候補タグをクリックすると追加される', () => {
    const suggestions = ['JavaScript', 'TypeScript'];
    render(<TagSelector value="" onChange={mockOnChange} suggestions={suggestions} />);

    const jsTag = screen.getByText('JavaScript');
    fireEvent.click(jsTag);

    expect(mockOnChange).toHaveBeenCalledWith('#JavaScript');
  });

  it('重複するタグは追加されない', () => {
    render(<TagSelector value="#React" onChange={mockOnChange} />);

    const input = screen.getByPlaceholderText(/タグを入力/);
    fireEvent.change(input, { target: { value: 'React' } });
    fireEvent.keyDown(input, { key: 'Enter' });

    // onChangeが呼ばれないことを確認
    expect(mockOnChange).not.toHaveBeenCalled();
  });

  it('タグは自動的に#が付与される', () => {
    render(<TagSelector value="" onChange={mockOnChange} />);

    const input = screen.getByPlaceholderText(/タグを入力/);
    fireEvent.change(input, { target: { value: 'React' } });
    fireEvent.keyDown(input, { key: 'Enter' });

    expect(mockOnChange).toHaveBeenCalledWith('#React');
  });

  it('カンマ区切りでタグが追加される', () => {
    render(<TagSelector value="" onChange={mockOnChange} />);

    const input = screen.getByPlaceholderText(/タグを入力/);
    fireEvent.change(input, { target: { value: 'React,TypeScript' } });
    fireEvent.keyDown(input, { key: 'Enter' });

    expect(mockOnChange).toHaveBeenCalledWith('#React,#TypeScript');
  });

  it('最大タグ数を超えると追加できない', () => {
    render(<TagSelector value="#Tag1,#Tag2,#Tag3" onChange={mockOnChange} maxTags={3} />);

    const input = screen.getByPlaceholderText(/タグを入力/);
    fireEvent.change(input, { target: { value: 'Tag4' } });
    fireEvent.keyDown(input, { key: 'Enter' });

    expect(mockOnChange).not.toHaveBeenCalled();
  });

  it('選択済みタグにホバーエフェクトがある', () => {
    render(<TagSelector value="#React" onChange={mockOnChange} />);

    const reactTag = screen.getByText('React').closest('button');
    expect(reactTag).toHaveClass('hover:bg-red-500/30');
  });

  it('候補タグと選択済みタグが区別される', () => {
    const suggestions = ['React', 'TypeScript'];
    const { container } = render(<TagSelector value="#React" onChange={mockOnChange} suggestions={suggestions} />);

    // 選択済みタグエリアのボタン（bg-blue-500/20 + hover:bg-red-500/30）
    const selectedButtons = container.querySelectorAll('button.bg-blue-500\\/20.hover\\:bg-red-500\\/30');
    expect(selectedButtons.length).toBe(1);

    // 候補タグエリアの未選択タグ（bg-gray-800/50）
    const suggestionButtons = container.querySelectorAll('button.bg-gray-800\\/50');
    expect(suggestionButtons.length).toBe(1); // TypeScriptのみ（Reactは選択済みなので異なるスタイル）
  });
});
