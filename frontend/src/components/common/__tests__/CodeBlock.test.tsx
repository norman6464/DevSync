import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import CodeBlock from '../CodeBlock';

describe('CodeBlock', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    Object.assign(navigator, {
      clipboard: { writeText: vi.fn().mockResolvedValue(undefined) },
    });
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('コードが表示される', () => {
    render(<CodeBlock code="const x = 1;" />);
    expect(screen.getByText('const x = 1;')).toBeInTheDocument();
  });

  it('言語名が表示される', () => {
    render(<CodeBlock code="const x = 1;" language="typescript" />);
    expect(screen.getByText('typescript')).toBeInTheDocument();
  });

  it('行番号が表示される', () => {
    render(<CodeBlock code={'line1\nline2\nline3'} showLineNumbers />);
    expect(screen.getByText('1')).toBeInTheDocument();
    expect(screen.getByText('2')).toBeInTheDocument();
    expect(screen.getByText('3')).toBeInTheDocument();
  });

  it('コピーボタンが表示される', () => {
    const { container } = render(<CodeBlock code="test" showCopy />);
    expect(container.querySelector('.lucide-copy')).toBeInTheDocument();
  });

  it('コピーボタンクリックでクリップボードにコピー', async () => {
    const user = userEvent.setup({ advanceTimers: vi.advanceTimersByTime });
    render(<CodeBlock code="copied text" showCopy />);
    await user.click(screen.getByRole('button'));
    expect(navigator.clipboard.writeText).toHaveBeenCalledWith('copied text');
  });

  it('コピー後にチェックアイコン表示', async () => {
    const user = userEvent.setup({ advanceTimers: vi.advanceTimersByTime });
    const { container } = render(<CodeBlock code="test" showCopy />);
    await user.click(screen.getByRole('button'));
    expect(container.querySelector('.lucide-check')).toBeInTheDocument();
  });

  it('pre要素が使用される', () => {
    const { container } = render(<CodeBlock code="test" />);
    expect(container.querySelector('pre')).toBeInTheDocument();
  });

  it('code要素が使用される', () => {
    const { container } = render(<CodeBlock code="test" />);
    expect(container.querySelector('code')).toBeInTheDocument();
  });

  it('タイトルが表示される', () => {
    render(<CodeBlock code="test" title="example.ts" />);
    expect(screen.getByText('example.ts')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<CodeBlock code="test" className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
