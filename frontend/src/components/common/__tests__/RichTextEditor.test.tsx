import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import RichTextEditor from '../RichTextEditor';

describe('RichTextEditor', () => {
  it('エディター領域が表示される', () => {
    render(<RichTextEditor value="" onChange={() => {}} />);
    expect(screen.getByRole('textbox')).toBeInTheDocument();
  });

  it('ツールバーが表示される', () => {
    render(<RichTextEditor value="" onChange={() => {}} />);
    expect(screen.getByTestId('toolbar')).toBeInTheDocument();
  });

  it('太字ボタンが表示される', () => {
    render(<RichTextEditor value="" onChange={() => {}} />);
    expect(screen.getByLabelText('太字')).toBeInTheDocument();
  });

  it('斜体ボタンが表示される', () => {
    render(<RichTextEditor value="" onChange={() => {}} />);
    expect(screen.getByLabelText('斜体')).toBeInTheDocument();
  });

  it('下線ボタンが表示される', () => {
    render(<RichTextEditor value="" onChange={() => {}} />);
    expect(screen.getByLabelText('下線')).toBeInTheDocument();
  });

  it('プレースホルダーが表示される', () => {
    render(<RichTextEditor value="" onChange={() => {}} placeholder="ここに入力..." />);
    expect(screen.getByText('ここに入力...')).toBeInTheDocument();
  });

  it('値が表示される', () => {
    render(<RichTextEditor value="テスト内容" onChange={() => {}} />);
    expect(screen.getByText('テスト内容')).toBeInTheDocument();
  });

  it('ラベルが表示される', () => {
    render(<RichTextEditor value="" onChange={() => {}} label="本文" />);
    expect(screen.getByText('本文')).toBeInTheDocument();
  });

  it('エラーメッセージが表示される', () => {
    render(<RichTextEditor value="" onChange={() => {}} error="入力してください" />);
    expect(screen.getByText('入力してください')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<RichTextEditor value="" onChange={() => {}} className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
