import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Select from '../Select';

const options = [
  { value: 'react', label: 'React' },
  { value: 'vue', label: 'Vue' },
  { value: 'angular', label: 'Angular' },
];

describe('Select', () => {
  it('プレースホルダーが表示される', () => {
    render(<Select options={options} value="" onChange={() => {}} placeholder="選択してください" />);
    expect(screen.getByText('選択してください')).toBeInTheDocument();
  });

  it('選択値が表示される', () => {
    render(<Select options={options} value="react" onChange={() => {}} />);
    expect(screen.getByText('React')).toBeInTheDocument();
  });

  it('クリックでドロップダウンが開く', async () => {
    const user = userEvent.setup();
    render(<Select options={options} value="" onChange={() => {}} placeholder="選択" />);
    await user.click(screen.getByText('選択'));
    expect(screen.getByText('Vue')).toBeInTheDocument();
    expect(screen.getByText('Angular')).toBeInTheDocument();
  });

  it('オプション選択でコールバックが呼ばれる', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<Select options={options} value="" onChange={onChange} placeholder="選択" />);
    await user.click(screen.getByText('選択'));
    await user.click(screen.getByText('Vue'));
    expect(onChange).toHaveBeenCalledWith('vue');
  });

  it('選択後にドロップダウンが閉じる', async () => {
    const user = userEvent.setup();
    render(<Select options={options} value="" onChange={() => {}} placeholder="選択" />);
    await user.click(screen.getByText('選択'));
    await user.click(screen.getByText('Vue'));
    expect(screen.queryByText('Angular')).not.toBeInTheDocument();
  });

  it('シェブロンアイコンが表示される', () => {
    const { container } = render(<Select options={options} value="" onChange={() => {}} />);
    expect(container.querySelector('.lucide-chevron-down')).toBeInTheDocument();
  });

  it('無効状態が適用される', () => {
    const { container } = render(<Select options={options} value="" onChange={() => {}} disabled />);
    expect(container.querySelector('.opacity-50')).toBeInTheDocument();
  });

  it('ラベルが表示される', () => {
    render(<Select options={options} value="" onChange={() => {}} label="フレームワーク" />);
    expect(screen.getByText('フレームワーク')).toBeInTheDocument();
  });

  it('エラーメッセージが表示される', () => {
    render(<Select options={options} value="" onChange={() => {}} error="必須項目です" />);
    expect(screen.getByText('必須項目です')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<Select options={options} value="" onChange={() => {}} className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
