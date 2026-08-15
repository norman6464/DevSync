import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import TagInput from '../TagInput';

describe('TagInput', () => {
  it('既存タグが表示される', () => {
    render(<TagInput value={['React', 'TypeScript']} onChange={() => {}} />);
    expect(screen.getByText('React')).toBeInTheDocument();
    expect(screen.getByText('TypeScript')).toBeInTheDocument();
  });

  it('入力フィールドが表示される', () => {
    render(<TagInput value={[]} onChange={() => {}} />);
    expect(screen.getByRole('textbox')).toBeInTheDocument();
  });

  it('Enterキーでタグが追加される', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<TagInput value={['React']} onChange={onChange} />);
    await user.type(screen.getByRole('textbox'), 'Vue{Enter}');
    expect(onChange).toHaveBeenCalledWith(['React', 'Vue']);
  });

  it('空入力でタグが追加されない', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<TagInput value={['React']} onChange={onChange} />);
    await user.type(screen.getByRole('textbox'), '{Enter}');
    expect(onChange).not.toHaveBeenCalled();
  });

  it('重複タグが追加されない', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<TagInput value={['React']} onChange={onChange} />);
    await user.type(screen.getByRole('textbox'), 'React{Enter}');
    expect(onChange).not.toHaveBeenCalled();
  });

  it('削除ボタンでタグが削除される', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<TagInput value={['React', 'Vue']} onChange={onChange} />);
    const removeButtons = screen.getAllByLabelText('削除');
    await user.click(removeButtons[0]);
    expect(onChange).toHaveBeenCalledWith(['Vue']);
  });

  it('最大タグ数に達すると入力が無効になる', () => {
    render(<TagInput value={['React', 'Vue', 'Angular']} onChange={() => {}} maxTags={3} />);
    expect(screen.getByRole('textbox')).toBeDisabled();
  });

  it('プレースホルダーが表示される', () => {
    render(<TagInput value={[]} onChange={() => {}} placeholder="タグを追加" />);
    expect(screen.getByPlaceholderText('タグを追加')).toBeInTheDocument();
  });

  it('ラベルが表示される', () => {
    render(<TagInput value={[]} onChange={() => {}} label="スキル" />);
    expect(screen.getByText('スキル')).toBeInTheDocument();
  });

  it('無効状態で入力が無効になる', () => {
    render(<TagInput value={[]} onChange={() => {}} disabled />);
    expect(screen.getByRole('textbox')).toBeDisabled();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<TagInput value={[]} onChange={() => {}} className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('maxLength が入力欄に適用される', () => {
    render(<TagInput value={[]} onChange={() => {}} maxLength={50} />);

    expect(screen.getByRole('textbox')).toHaveAttribute('maxlength', '50');
  });
});
