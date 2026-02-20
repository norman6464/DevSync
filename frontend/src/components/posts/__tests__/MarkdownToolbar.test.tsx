import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import MarkdownToolbar from '../MarkdownToolbar';

describe('MarkdownToolbar', () => {
  it('ツールバーがrole="toolbar"で表示される', () => {
    render(<MarkdownToolbar onAction={vi.fn()} />);
    expect(screen.getByRole('toolbar')).toBeInTheDocument();
  });

  it('見出しボタンが表示される', () => {
    render(<MarkdownToolbar onAction={vi.fn()} />);
    expect(screen.getByTitle('見出し')).toBeInTheDocument();
  });

  it('太字ボタンがfont-boldクラスを持つ', () => {
    render(<MarkdownToolbar onAction={vi.fn()} />);
    const boldBtn = screen.getByTitle('太字');
    expect(boldBtn).toHaveClass('font-bold');
  });

  it('斜体ボタンがitalicクラスを持つ', () => {
    render(<MarkdownToolbar onAction={vi.fn()} />);
    const italicBtn = screen.getByTitle('斜体');
    expect(italicBtn).toHaveClass('italic');
  });

  it('取り消し線ボタンがline-throughクラスを持つ', () => {
    render(<MarkdownToolbar onAction={vi.fn()} />);
    const strikeBtn = screen.getByTitle('取り消し線');
    expect(strikeBtn).toHaveClass('line-through');
  });

  it('ボタンクリック時にonActionが正しいアクション名で呼ばれる', () => {
    const onAction = vi.fn();
    render(<MarkdownToolbar onAction={onAction} />);
    fireEvent.click(screen.getByTitle('太字'));
    expect(onAction).toHaveBeenCalledWith('bold');
  });

  it('見出しボタンクリックでheadingアクションが呼ばれる', () => {
    const onAction = vi.fn();
    render(<MarkdownToolbar onAction={onAction} />);
    fireEvent.click(screen.getByTitle('見出し'));
    expect(onAction).toHaveBeenCalledWith('heading');
  });

  it('コードブロックボタンクリックでcodeblockアクションが呼ばれる', () => {
    const onAction = vi.fn();
    render(<MarkdownToolbar onAction={onAction} />);
    fireEvent.click(screen.getByTitle('コードブロック'));
    expect(onAction).toHaveBeenCalledWith('codeblock');
  });

  it('dividerが3つ表示される', () => {
    const { container } = render(<MarkdownToolbar onAction={vi.fn()} />);
    const dividers = container.querySelectorAll('.bg-gray-600.mx-1');
    expect(dividers.length).toBe(3);
  });

  it('全てのボタンにtype="button"が設定されている', () => {
    render(<MarkdownToolbar onAction={vi.fn()} />);
    const buttons = screen.getAllByRole('button');
    buttons.forEach((btn) => {
      expect(btn).toHaveAttribute('type', 'button');
    });
  });
});
