import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import CheckboxGroup from '../CheckboxGroup';

const options = [
  { value: 'react', label: 'React' },
  { value: 'vue', label: 'Vue' },
  { value: 'angular', label: 'Angular', description: 'Google製フレームワーク' },
];

describe('CheckboxGroup', () => {
  it('すべてのオプションが表示される', () => {
    render(<CheckboxGroup options={options} value={[]} onChange={() => {}} />);
    expect(screen.getByText('React')).toBeInTheDocument();
    expect(screen.getByText('Vue')).toBeInTheDocument();
    expect(screen.getByText('Angular')).toBeInTheDocument();
  });

  it('選択済みの値がチェックされる', () => {
    render(<CheckboxGroup options={options} value={['react']} onChange={() => {}} />);
    const checkboxes = screen.getAllByRole('checkbox');
    expect(checkboxes[0]).toBeChecked();
    expect(checkboxes[1]).not.toBeChecked();
  });

  it('クリックで値が追加される', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<CheckboxGroup options={options} value={['react']} onChange={onChange} />);
    await user.click(screen.getByText('Vue'));
    expect(onChange).toHaveBeenCalledWith(['react', 'vue']);
  });

  it('クリックで値が削除される', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<CheckboxGroup options={options} value={['react', 'vue']} onChange={onChange} />);
    await user.click(screen.getByText('React'));
    expect(onChange).toHaveBeenCalledWith(['vue']);
  });

  it('説明テキストが表示される', () => {
    render(<CheckboxGroup options={options} value={[]} onChange={() => {}} />);
    expect(screen.getByText('Google製フレームワーク')).toBeInTheDocument();
  });

  it('ラベルが表示される', () => {
    render(<CheckboxGroup options={options} value={[]} onChange={() => {}} label="フレームワーク" />);
    expect(screen.getByText('フレームワーク')).toBeInTheDocument();
  });

  it('エラーメッセージが表示される', () => {
    render(<CheckboxGroup options={options} value={[]} onChange={() => {}} error="選択してください" />);
    expect(screen.getByText('選択してください')).toBeInTheDocument();
  });

  it('無効状態で全チェックボックスが無効になる', () => {
    render(<CheckboxGroup options={options} value={[]} onChange={() => {}} disabled />);
    const checkboxes = screen.getAllByRole('checkbox');
    checkboxes.forEach((cb) => expect(cb).toBeDisabled());
  });

  it('全選択ボタンで全て選択される', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<CheckboxGroup options={options} value={[]} onChange={onChange} showSelectAll />);
    await user.click(screen.getByText('すべて選択'));
    expect(onChange).toHaveBeenCalledWith(['react', 'vue', 'angular']);
  });

  it('全解除ボタンで全て解除される', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<CheckboxGroup options={options} value={['react', 'vue', 'angular']} onChange={onChange} showSelectAll />);
    await user.click(screen.getByText('すべて解除'));
    expect(onChange).toHaveBeenCalledWith([]);
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<CheckboxGroup options={options} value={[]} onChange={() => {}} className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
