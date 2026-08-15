import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import CodeBlock from '../CodeBlock';

const mockWriteText = vi.fn().mockResolvedValue(undefined);

describe('CodeBlock', () => {
  beforeEach(() => {
    Object.defineProperty(navigator, 'clipboard', {
      value: { writeText: mockWriteText },
      writable: true,
      configurable: true,
    });
    mockWriteText.mockClear();
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
    render(<CodeBlock code="test" showCopy />);
    expect(screen.getByRole('button')).toBeInTheDocument();
  });

  it('コピーボタンクリックでクリップボードにコピー', async () => {
    render(<CodeBlock code="copied text" showCopy />);
    fireEvent.click(screen.getByRole('button'));
    await waitFor(() => {
      expect(mockWriteText).toHaveBeenCalledWith('copied text');
    });
  });

  it('コピー後にチェックアイコン表示', async () => {
    const { container } = render(<CodeBlock code="test" showCopy />);
    fireEvent.click(screen.getByRole('button'));
    await waitFor(() => {
      const svgs = container.querySelectorAll('svg');
      expect(svgs.length).toBeGreaterThan(0);
    });
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
