import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import TextArea from '../TextArea';

describe('TextArea', () => {
  it('テキストエリアが表示される', () => {
    render(<TextArea value="" onChange={() => {}} />);
    expect(screen.getByRole('textbox')).toBeInTheDocument();
  });

  it('値が表示される', () => {
    render(<TextArea value="テスト内容" onChange={() => {}} />);
    expect(screen.getByRole('textbox')).toHaveValue('テスト内容');
  });

  it('値の変更でコールバックが呼ばれる', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<TextArea value="" onChange={onChange} />);
    await user.type(screen.getByRole('textbox'), 'a');
    expect(onChange).toHaveBeenCalled();
  });

  it('プレースホルダーが表示される', () => {
    render(<TextArea value="" onChange={() => {}} placeholder="入力してください" />);
    expect(screen.getByPlaceholderText('入力してください')).toBeInTheDocument();
  });

  it('ラベルが表示される', () => {
    render(<TextArea value="" onChange={() => {}} label="説明" />);
    expect(screen.getByText('説明')).toBeInTheDocument();
  });

  it('エラーメッセージが表示される', () => {
    render(<TextArea value="" onChange={() => {}} error="必須項目です" />);
    expect(screen.getByText('必須項目です')).toBeInTheDocument();
  });

  it('文字数カウントが表示される', () => {
    render(<TextArea value="hello" onChange={() => {}} maxLength={100} showCount />);
    expect(screen.getByText('5 / 100')).toBeInTheDocument();
  });

  it('無効状態が適用される', () => {
    render(<TextArea value="" onChange={() => {}} disabled />);
    expect(screen.getByRole('textbox')).toBeDisabled();
  });

  it('読み取り専用が適用される', () => {
    render(<TextArea value="内容" onChange={() => {}} readOnly />);
    expect(screen.getByRole('textbox')).toHaveAttribute('readonly');
  });

  it('行数が設定される', () => {
    render(<TextArea value="" onChange={() => {}} rows={5} />);
    expect(screen.getByRole('textbox')).toHaveAttribute('rows', '5');
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<TextArea value="" onChange={() => {}} className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('エラー時にボーダーが赤くなる', () => {
    render(<TextArea value="" onChange={() => {}} error="エラー" />);
    const textarea = screen.getByRole('textbox');
    expect(textarea.className).toContain('border-red-500');
  });
});
