import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import RadioGroup from '../RadioGroup';

const options = [
  { value: 'a', label: 'オプションA' },
  { value: 'b', label: 'オプションB' },
  { value: 'c', label: 'オプションC' },
];

describe('RadioGroup', () => {
  it('全てのオプションが表示される', () => {
    render(<RadioGroup options={options} value="" onChange={() => {}} />);
    expect(screen.getByText('オプションA')).toBeInTheDocument();
    expect(screen.getByText('オプションB')).toBeInTheDocument();
    expect(screen.getByText('オプションC')).toBeInTheDocument();
  });

  it('ラジオボタンが表示される', () => {
    render(<RadioGroup options={options} value="" onChange={() => {}} />);
    const radios = screen.getAllByRole('radio');
    expect(radios.length).toBe(3);
  });

  it('選択値が反映される', () => {
    render(<RadioGroup options={options} value="b" onChange={() => {}} />);
    const radios = screen.getAllByRole('radio');
    expect(radios[1]).toBeChecked();
  });

  it('クリックでコールバックが呼ばれる', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<RadioGroup options={options} value="" onChange={onChange} />);
    await user.click(screen.getByText('オプションB'));
    expect(onChange).toHaveBeenCalledWith('b');
  });

  it('ラベルが表示される', () => {
    render(<RadioGroup options={options} value="" onChange={() => {}} label="選択してください" />);
    expect(screen.getByText('選択してください')).toBeInTheDocument();
  });

  it('説明付きオプションが表示される', () => {
    const opts = [{ value: 'a', label: 'A', description: '説明文' }];
    render(<RadioGroup options={opts} value="" onChange={() => {}} />);
    expect(screen.getByText('説明文')).toBeInTheDocument();
  });

  it('水平レイアウトが適用される', () => {
    const { container } = render(<RadioGroup options={options} value="" onChange={() => {}} direction="horizontal" />);
    expect(container.querySelector('.flex-row')).toBeInTheDocument();
  });

  it('垂直レイアウトがデフォルト', () => {
    const { container } = render(<RadioGroup options={options} value="" onChange={() => {}} />);
    expect(container.querySelector('.flex-col')).toBeInTheDocument();
  });

  it('無効状態が適用される', () => {
    render(<RadioGroup options={options} value="" onChange={() => {}} disabled />);
    const radios = screen.getAllByRole('radio');
    radios.forEach(r => expect(r).toBeDisabled());
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<RadioGroup options={options} value="" onChange={() => {}} className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
