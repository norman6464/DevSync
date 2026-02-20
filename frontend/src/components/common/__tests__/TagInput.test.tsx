import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import TagInput from '../TagInput';

const defaultProps = {
  tags: [] as string[],
  onChange: vi.fn(),
};

describe('TagInput', () => {
  it('入力フィールドと追加ボタンが表示される', () => {
    render(<TagInput {...defaultProps} />);
    expect(screen.getByRole('textbox')).toBeInTheDocument();
    expect(screen.getByText('追加')).toBeInTheDocument();
  });

  it('追加ボタンクリックでonChangeが呼ばれる', () => {
    const onChange = vi.fn();
    render(<TagInput tags={[]} onChange={onChange} />);
    const input = screen.getByRole('textbox');
    fireEvent.change(input, { target: { value: 'React' } });
    fireEvent.click(screen.getByText('追加'));
    expect(onChange).toHaveBeenCalledWith(['React']);
  });

  it('既存タグが表示される', () => {
    render(<TagInput tags={['React', 'TypeScript']} onChange={vi.fn()} />);
    expect(screen.getByText('React')).toBeInTheDocument();
    expect(screen.getByText('TypeScript')).toBeInTheDocument();
  });

  it('prefix付きでタグが表示される', () => {
    render(<TagInput tags={['frontend']} onChange={vi.fn()} prefix="#" />);
    expect(screen.getByText('#frontend')).toBeInTheDocument();
  });

  it('タグ削除ボタンクリックでonChangeが呼ばれる', () => {
    const onChange = vi.fn();
    render(<TagInput tags={['React', 'Go']} onChange={onChange} />);
    const removeButtons = screen.getAllByRole('button').filter(
      btn => btn.querySelector('svg.w-3')
    );
    fireEvent.click(removeButtons[0]);
    expect(onChange).toHaveBeenCalledWith(['Go']);
  });

  it('空入力では追加されない', () => {
    const onChange = vi.fn();
    render(<TagInput tags={[]} onChange={onChange} />);
    fireEvent.click(screen.getByText('追加'));
    expect(onChange).not.toHaveBeenCalled();
  });

  it('labelが表示される', () => {
    render(<TagInput {...defaultProps} label="技術スタック" />);
    expect(screen.getByText('技術スタック')).toBeInTheDocument();
  });

  it('id属性が設定される', () => {
    render(<TagInput {...defaultProps} id="test-tags" />);
    expect(screen.getByRole('textbox')).toHaveAttribute('id', 'test-tags');
  });

  it('placeholderが表示される', () => {
    render(<TagInput {...defaultProps} placeholder="タグを入力" />);
    expect(screen.getByPlaceholderText('タグを入力')).toBeInTheDocument();
  });

  it('maxLength属性が設定される', () => {
    render(<TagInput {...defaultProps} maxLength={100} />);
    expect(screen.getByRole('textbox')).toHaveAttribute('maxLength', '100');
  });
});
